package models

type User struct {
	Base
	Name          string `json:"name" gorm:"type:text;not null"`
	Email         string `json:"email" gorm:"uniqueIndex;not null"`
	AvatarURL     string `json:"avatar_url,omitempty"`
	SpotifyUserID string `json:"spotify_user_id,omitempty" gorm:"uniqueIndex;not null"`

	SpotifyToken         *SpotifyToken         `json:"-" gorm:"foreignKey:UserID"`
	Subscriptions        []Subscription        `gorm:"foreignKey:UserID" json:"-"`
	NotificationChannels []NotificationChannel `json:"-" gorm:"foreignKey:UserID"`
}
