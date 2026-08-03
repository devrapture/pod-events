package repositories

import (
	"context"
	"testing"

	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestCreateChannelDuplicate(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewChannelRepository(db, "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=")
	channel := &models.NotificationChannel{
		UserID:      uuid.New(),
		ChannelType: models.ChannelTypeTelegram,
		Destination: "123456789",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "notification_channels"`).
		WillReturnError(&pgconn.PgError{Code: "23505", ConstraintName: "idx_user_channel_destination"})
	mock.ExpectRollback()

	err := repo.Create(context.Background(), channel)

	assert.ErrorIs(t, err, apperrors.ErrNotificationChannelAlreadyExists)
	assert.NoError(t, mock.ExpectationsWereMet())
}
