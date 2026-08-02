export interface APIResponse<T = unknown> {
	success?: boolean;
	message?: string;
	data?: T;
	error?: ErrorInfo;
	meta?: PaginationMeta;
}

export interface ErrorInfo {
	message?: string;
	details?: Record<string, string>;
}

export interface PaginationMeta {
	page?: number;
	pages?: number;
	per_page?: number;
	total?: number;
}

export interface User {
	id: string;
	name: string;
	email: string;
	avatar_url: string;
	spotify_user_id: string;
}

export interface AuthExchangeRequest {
	code: string;
}

export interface AuthExchangeResponse {
	token: string;
	user: User;
}

export type ChannelType =
	| "slack_webhook"
	| "discord_webhook"
	| "whatsapp"
	| "telegram";

export interface CreateChannelRequest {
	channel_type: ChannelType;
	destination: string;
	label?: string;
}

export interface ToggleActiveChannelRequest {
	is_active: boolean;
}

export type TrackingStatus = "tracked" | "untracked";

export interface SavedShowResponse {
	id: string;
	name: string;
	description: string;
	total_episodes: number;
	image_url: string;
	spotify_url: string;
	added_at: string;
	is_tracked?: boolean;
	tracking_status?: TrackingStatus;
}

export interface SavedShowWithTracking extends SavedShowResponse {
	is_tracked: boolean;
	tracking_status: TrackingStatus;
}

export interface SubscribeShowsRequest {
	spotify_show_ids: string[];
}

export interface PodcastShowResponse {
	id: string;
	spotify_show_id: string;
	name: string;
	description: string;
	image_url: string;
	spotify_url: string;
	latest_episode_id?: string;
	latest_episode_published_at?: string;
}

export interface SubscriptionResponse {
	id: string;
	user_id: string;
	podcast_show_id: string;
	podcast_show: PodcastShowResponse;
	created_at: string;
	updated_at: string;
}

export interface NotificationChannel {
	id: string;
	user_id: string;
	channel_type: ChannelType;
	destination: string;
	label?: string;
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export interface TelegramLinkResponse {
	url: string;
}

export interface HealthResponse {
	status: string;
}

export interface DashboardSetupItem {
	key: string;
	label: string;
	completed: boolean;
}

export interface DashboardSetup {
	completed: number;
	total: number;
	percent: number;
	items: DashboardSetupItem[];
}

export interface DashboardStats {
	podcasts_tracked: number;
	active_channels: number;
	new_episodes_this_week: number;
	notification_sent: number;
}

export interface DashboardSummary {
	setup: DashboardSetup;
	stats: DashboardStats;
}

export interface ShowSearchParams {
	q: string;
	limit?: number;
	offset?: number;
}

export interface SavedShowsParams {
	q?: string;
}
