package models

import (
	"time"

	"github.com/google/uuid"
)

type TelegramConnection struct {
	Base
	UserID    uuid.UUID `json:"user_id" gorm:"uniqueIndex;not null"`
	TokenHash string    `json:"token_hash" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expired_at" gorm:"not null"`
	Consumed  bool      `json:"consumed" gorm:"not null;default:false"`

	User User `json:"-" gorm:"foreignKey:UserID"`
}
