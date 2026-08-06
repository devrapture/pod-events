package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dashboardTokenRepositoryStub struct{}

func (dashboardTokenRepositoryStub) Upsert(context.Context, *models.SpotifyToken) error {
	return nil
}

func (dashboardTokenRepositoryStub) UpdateAccessToken(context.Context, uuid.UUID, string, time.Time) error {
	return nil
}

func (dashboardTokenRepositoryStub) GetByUserID(context.Context, uuid.UUID) (*models.SpotifyToken, error) {
	return &models.SpotifyToken{}, nil
}

func TestGetDashboardSummaryReturnsRecentSentEpisodes(t *testing.T) {
	db, mock := setupMockDB(t)
	userID := uuid.New()
	releasedAt := time.Date(2026, time.August, 5, 12, 0, 0, 0, time.UTC)
	repo := NewDashboardSummaryRepository(db, dashboardTokenRepositoryStub{})

	mock.ExpectQuery(`SELECT count\(\*\) FROM "subscriptions".*`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`SELECT count\(\*\) FROM "notification_channels".*`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`SELECT count\(\*\) FROM "episodes".*`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`SELECT count\(\*\) FROM "notification_logs".*`).
		WithArgs(userID, models.NotificationStatusSent).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	mock.ExpectQuery(`(?s)SELECT .*episodes\.name AS episode_title.*FROM "notification_logs".*GROUP BY episodes\.id, podcast_shows\.id ORDER BY MAX\(notification_logs\.sent_at\) DESC LIMIT \$3`).
		WithArgs(userID, models.NotificationStatusSent, 5).
		WillReturnRows(sqlmock.NewRows([]string{
			"episode_title",
			"podcast_name",
			"release_date",
			"duration",
			"url",
		}).AddRow(
			"A recent episode",
			"A podcast",
			releasedAt,
			3_600_000,
			"https://open.spotify.com/episode/test",
		))

	summary, err := repo.GetDashboardSummary(context.Background(), userID)

	require.NoError(t, err)
	require.Len(t, summary.RecentEpisodes, 1)
	assert.Equal(t, "A recent episode", summary.RecentEpisodes[0].EpisodeTitle)
	assert.Equal(t, "A podcast", summary.RecentEpisodes[0].PodcastName)
	assert.Equal(t, releasedAt, summary.RecentEpisodes[0].ReleaseDate)
	assert.Equal(t, 3_600_000, summary.RecentEpisodes[0].Duration)
	assert.Equal(t, "https://open.spotify.com/episode/test", summary.RecentEpisodes[0].URL)
	assert.NoError(t, mock.ExpectationsWereMet())
}
