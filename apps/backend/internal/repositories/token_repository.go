package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/devrapture/pod-events/internal/models"
	appcrypto "github.com/devrapture/pod-events/pkg/crypto"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TokenRepository interface {
	Upsert(ctx context.Context, token *models.SpotifyToken) error
	UpdateAccessToken(ctx context.Context, userID uuid.UUID, accessToken string, expires_at time.Time) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.SpotifyToken, error)
}

type tokenRepository struct {
	db            *gorm.DB
	encryptionKey string
}

func NewTokenRepository(db *gorm.DB, encryptionKey string) TokenRepository {
	return &tokenRepository{
		db:            db,
		encryptionKey: encryptionKey,
	}
}

func (r *tokenRepository) Upsert(ctx context.Context, token *models.SpotifyToken) error {
	accessTokenEncrypted, err := appcrypto.EncryptText(token.AccessToken, r.encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt access token: %w", err)
	}
	refreshTokenEncrypted, err := appcrypto.EncryptText(token.RefreshToken, r.encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt refresh token: %w", err)
	}

	toSave := *token
	toSave.AccessToken = accessTokenEncrypted
	toSave.RefreshToken = refreshTokenEncrypted

	err = r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{
			Name: "user_id",
		}},
		DoUpdates: clause.AssignmentColumns([]string{
			"access_token",
			"refresh_token",
			"expires_at",
			"scope",
			"updated_at",
		}),
	}).Create(&toSave).Error
	if err != nil {
		return fmt.Errorf("upsert spotify token: %w", err)
	}
	// Copy the generated ID back to the caller's struct
	token.ID = toSave.ID
	return nil
}

// UpdateAccessToken updates only the access token and expiry (after a refresh).
func (r *tokenRepository) UpdateAccessToken(ctx context.Context, userID uuid.UUID, accessToken string, expires_at time.Time) error {
	accessTokenEncrypted, err := appcrypto.EncryptText(accessToken, r.encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt access token: %w", err)
	}

	result := r.db.WithContext(ctx).Model(&models.SpotifyToken{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"access_token": accessTokenEncrypted,
		"expires_at":   expires_at,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update access token: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrorSpotifyTokenNotFound
	}
	return nil
}

func (r *tokenRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.SpotifyToken, error) {
	var spotifyToken models.SpotifyToken
	result := r.db.WithContext(ctx).First(&spotifyToken, "user_id = ?", userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrorSpotifyTokenNotFound
		}
		return nil, fmt.Errorf("failed to get spotify token: %w", result.Error)
	}
	accessToken, err := appcrypto.DecryptText(spotifyToken.AccessToken, r.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt access token: %w", err)
	}

	refreshToken, err := appcrypto.DecryptText(spotifyToken.RefreshToken, r.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt refresh token: %w", err)
	}

	spotifyToken.AccessToken = accessToken
	spotifyToken.RefreshToken = refreshToken

	return &spotifyToken, nil
}

// Delete removes the Spotify token for a user (when they disconnect Spotify).
func (r *tokenRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.SpotifyToken{}).Error
}
