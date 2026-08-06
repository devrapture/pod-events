package dto

type DashboardSummaryDTO struct {
	Setup DashboardSetupResponse `json:"setup"`
	Stats DashboardStatsResponse `json:"stats"`
}

type DashboardSetupResponse struct {
	Completed      int                      `json:"completed"`
	Total          int                      `json:"total"`
	Percent        int                      `json:"percent"`
	Items          []DashboardItems         `json:"items"`
	RecentEpisodes []RecentEpisodesResponse `json:"recent_episodes"`
}

type DashboardStatsResponse struct {
	PodcastsTracked     int64 `json:"podcasts_tracked"`
	ActiveChannels      int64 `json:"active_channels"`
	NewEpisodesThisWeek int64 `json:"new_episodes_this_week"`
	NotificationSent    int64 `json:"notification_sent"`
}

type RecentEpisodesResponse struct {
	EpisodeTitle string `json:"episode_title"`
	PodcastName  string `json:"podcast_name"`
	ReleaseDate  string `json:"release_date"`
	Duration     string `json:"duration"`
	URL          string `json:"url"`
}

type DashboardItems struct {
	Key       string `json:"key"`
	Label     string `json:"label"`
	Completed bool   `json:"completed"`
}
