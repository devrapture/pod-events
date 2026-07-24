// internal/dto/subscription_dto.go

package dto

import (
	"time"

	"github.com/devrapture/pod-events/internal/models"
	"github.com/google/uuid"
)

type SubscribeShowsRequest struct {
	SpotifyShowIDs []string `json:"spotify_show_ids" binding:"required,min=1,max=50,dive,required"`
}

// PodcastShowResponse represents a podcast show entity.
type PodcastShowResponse struct {
	ID                       uuid.UUID  `json:"id"`
	SpotifyShowID            string     `json:"spotify_show_id"`
	Name                     string     `json:"name"`
	Description              string     `json:"description"`
	ImageURL                 string     `json:"image_url"`
	SpotifyURL               string     `json:"spotify_url"`
	LatestEpisodeID          *string    `json:"latest_episode_id,omitempty"`
	LatestEpisodePublishedAt *time.Time `json:"latest_episode_published_at,omitempty"`
}

// SubscriptionResponse represents a user's subscription to a podcast show.
type SubscriptionResponse struct {
	ID            uuid.UUID           `json:"id"`
	UserID        uuid.UUID           `json:"user_id"`
	PodcastShowID uuid.UUID           `json:"podcast_show_id"`
	PodcastShow   PodcastShowResponse `json:"podcast_show"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}

func ToSubscriptionResponse(subscription models.Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:            subscription.ID,
		UserID:        subscription.UserID,
		PodcastShowID: subscription.PodcastShowID,
		CreatedAt:     subscription.CreatedAt,
		UpdatedAt:     subscription.UpdatedAt,
		PodcastShow: PodcastShowResponse{
			ID:                       subscription.PodcastShow.ID,
			SpotifyShowID:            subscription.PodcastShow.SpotifyShowID,
			Name:                     subscription.PodcastShow.Name,
			Description:              subscription.PodcastShow.Description,
			ImageURL:                 subscription.PodcastShow.ImageURL,
			SpotifyURL:               subscription.PodcastShow.SpotifyURL,
			LatestEpisodeID:          subscription.PodcastShow.LatestEpisodeID,
			LatestEpisodePublishedAt: subscription.PodcastShow.LatestEpisodePublishedAt,
		},
	}
}

func ToSubscriptionResponses(subscriptions []models.Subscription) []SubscriptionResponse {
	responses := make([]SubscriptionResponse, 0, len(subscriptions))

	for _, subscription := range subscriptions {
		responses = append(responses, ToSubscriptionResponse(subscription))
	}

	return responses
}
