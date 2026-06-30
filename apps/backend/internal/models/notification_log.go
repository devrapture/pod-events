package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationStatus string

const (
	NotificationStatusSent   NotificationStatus = "sent"
	NotificationStatusFailed NotificationStatus = "failed"
)

type NotificationLog struct {
	Base
	UserID       uuid.UUID          `json:"user_id" gorm:"type:uuid;not null;index"`
	EpisodeID    uuid.UUID          `json:"episode_id" gorm:"type:uuid;not null;index"`
	ChannelType  ChannelType        `json:"channel_type" gorm:"not null"`
	Status       NotificationStatus `json:"status" gorm:"not null"`
	ErrorMessage string             `json:"error_message,omitempty" gorm:"type:text"`
	SentAt       *time.Time         `json:"sent_at,omitempty"`

	User    User    `json:"-" gorm:"foreignKey:UserID"`
	Episode Episode `json:"-" gorm:"foreignKey:EpisodeID"`
}
