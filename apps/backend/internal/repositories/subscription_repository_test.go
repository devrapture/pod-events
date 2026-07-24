package repositories

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSubscriptionsBatch_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewSubscriptionRepository(db)
	userID := uuid.New()
	subscriptions := []models.Subscription{
		{UserID: userID, PodcastShowID: uuid.New()},
		{UserID: userID, PodcastShowID: uuid.New()},
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "subscriptions".*VALUES \(.+\),\(.+\)`).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit()

	err := repo.CreateBatch(context.Background(), subscriptions)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, subscriptions[0].ID)
	assert.NotEqual(t, uuid.Nil, subscriptions[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateSubscriptionsBatch_Duplicate(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewSubscriptionRepository(db)
	subscriptions := []models.Subscription{
		{UserID: uuid.New(), PodcastShowID: uuid.New()},
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "subscriptions"`).
		WillReturnError(&pgconn.PgError{Code: "23505"})
	mock.ExpectRollback()

	err := repo.CreateBatch(context.Background(), subscriptions)

	assert.ErrorIs(t, err, apperrors.ErrSubscriptionAlreadyExists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateSubscriptionsBatch_Empty(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewSubscriptionRepository(db)

	err := repo.CreateBatch(context.Background(), nil)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
