package models

import (
	"time"

	"github.com/google/uuid"
)

type TelegramConnection struct {
	Base
	UserID    uuid.UUID `json:"user_id" gorm:"not null;uniqueIndex:idx_telegram_connections_user_id,where:deleted_at IS NULL"`
	TokenHash string    `json:"token_hash" gorm:"not null;uniqueIndex:idx_telegram_connections_token_hash,where:deleted_at IS NULL"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	Consumed  bool      `json:"consumed" gorm:"not null;default:false"`

	User User `json:"-" gorm:"foreignKey:UserID"`
}
