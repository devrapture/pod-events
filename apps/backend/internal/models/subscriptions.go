package models

import "github.com/google/uuid"

type Subscription struct {
	Base
	UserID        uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex:idx_user_show"`
	PodcastShowID uuid.UUID `json:"podcast_show_id" gorm:"type:uuid;not null;uniqueIndex:idx_user_show"`

	User        User        `json:"-" gorm:"foreignKey:UserID"`
	PodcastShow PodcastShow `json:"-" gorm:"foreignKey:PodcastShowID"`
}
