package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn: mockDB,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return db, mock
}

func testUser() *models.User {
	return &models.User{
		Base: models.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now().Truncate(time.Microsecond),
			UpdatedAt: time.Now().Truncate(time.Microsecond),
		},
		Name:          "Test User",
		Email:         "test@example.com",
		AvatarURL:     "https://example.com/avatar.png",
		SpotifyUserID: "spotify_user_123",
	}
}

// ── Create ──────────────────────────────────────────────────────────────

func TestCreateUser_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewUserRepository(db)
	user := testUser()

	mock.ExpectBegin()
	mock.ExpectExec(
		`INSERT INTO "users" \("id","created_at","updated_at","deleted_at","name","email","avatar_url","spotify_user_id"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8\)`,
	).WithArgs(
		user.ID,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		nil,
		user.Name,
		user.Email,
		user.AvatarURL,
		user.SpotifyUserID,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_DBError(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewUserRepository(db)
	user := testUser()

	mock.ExpectBegin()
	mock.ExpectExec(
		`INSERT INTO "users" \("id","created_at","updated_at","deleted_at","name","email","avatar_url","spotify_user_id"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8\)`,
	).WillReturnError(errors.New("connection refused"))
	mock.ExpectRollback()

	err := repo.Create(context.Background(), user)

	assert.ErrorContains(t, err, "failed to create user")
	assert.ErrorContains(t, err, "connection refused")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByID ────────────────────────────────────────────────────────────

func TestGetByID_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewUserRepository(db)
	user := testUser()

	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"name", "email", "avatar_url", "spotify_user_id",
	}).AddRow(
		user.ID, user.CreatedAt, user.UpdatedAt, nil,
		user.Name, user.Email, user.AvatarURL, user.SpotifyUserID,
	)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(user.ID, 1).WillReturnRows(rows)

	got, err := repo.GetByID(context.Background(), user.ID)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, user.ID, got.ID)
	assert.Equal(t, user.Name, got.Name)
	assert.Equal(t, user.Email, got.Email)
	assert.Equal(t, user.AvatarURL, got.AvatarURL)
	assert.Equal(t, user.SpotifyUserID, got.SpotifyUserID)
	assert.False(t, got.CreatedAt.IsZero())
	assert.False(t, got.UpdatedAt.IsZero())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewUserRepository(db)
	userID := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(userID, 1).WillReturnError(gorm.ErrRecordNotFound)

	got, err := repo.GetByID(context.Background(), userID)

	assert.Nil(t, got)
	assert.ErrorIs(t, err, apperrors.ErrUserNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_DBError(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewUserRepository(db)
	userID := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(userID, 1).WillReturnError(errors.New("database is down"))

	got, err := repo.GetByID(context.Background(), userID)

	assert.Nil(t, got)
	assert.ErrorContains(t, err, "failed to get user")
	assert.ErrorContains(t, err, "database is down")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Update ─────────────────────────────────────────────────────────────

func TestUpdateUser_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewUserRepository(db)
	user := testUser()

	mock.ExpectBegin()
	mock.ExpectExec(
		`UPDATE "users" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"name"=\$4,"email"=\$5,"avatar_url"=\$6,"spotify_user_id"=\$7 WHERE "users"\."deleted_at" IS NULL AND "id" = \$8`,
	).WithArgs(
		user.CreatedAt,
		sqlmock.AnyArg(),
		nil,
		user.Name,
		user.Email,
		user.AvatarURL,
		user.SpotifyUserID,
		user.ID,
	).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser_DBError(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewUserRepository(db)
	user := testUser()

	mock.ExpectBegin()
	mock.ExpectExec(
		`UPDATE "users" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"name"=\$4,"email"=\$5,"avatar_url"=\$6,"spotify_user_id"=\$7 WHERE "users"\."deleted_at" IS NULL AND "id" = \$8`,
	).WillReturnError(errors.New("connection lost"))
	mock.ExpectRollback()

	err := repo.Update(context.Background(), user)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "connection lost")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Helpers ────────────────────────────────────────────────────────────

func TestNewUserRepository_NilDB(t *testing.T) {
	repo := NewUserRepository(nil)
	assert.NotNil(t, repo)
}
