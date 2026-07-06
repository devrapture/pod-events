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

export interface SavedShowResponse {
	id: string;
	name: string;
	description: string;
	total_episodes: number;
	image_url: string;
	spotify_url: string;
	added_at: string;
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

export interface ShowSearchParams {
	q: string;
	limit?: number;
	offset?: number;
}

export interface SavedShowsParams {
	q?: string;
}
