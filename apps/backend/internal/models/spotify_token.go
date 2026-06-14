package models

import (
	"time"

	"github.com/google/uuid"
)

type SpotifyToken struct {
	Base
	UserID       uuid.UUID `json:"user_id" gorm:"uniqueIndex;not null"`
	AccessToken  string    `json:"-" gorm:"not null"`
	RefreshToken string    `json:"-" gorm:"not null"`
	ExpiresAt    time.Time `json:"expired_at" gorm:"not null"`
	Scope        string    `json:"scope"`

	User User `json:"-" gorm:"foreignKey:UserID"`
}

func (t *SpotifyToken) IsExpired() bool {
	return time.Now().UTC().After(t.ExpiresAt.Add(-60 * time.Second))
}
