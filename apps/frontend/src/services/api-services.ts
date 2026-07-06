import { server, serverWithInterceptors } from "@/lib/axios-util";

import type {
	APIResponse,
	AuthExchangeRequest,
	AuthExchangeResponse,
	CreateChannelRequest,
	HealthResponse,
	NotificationChannel,
	SavedShowResponse,
	SubscriptionResponse,
	TelegramLinkResponse,
	ToggleActiveChannelRequest,
} from "./types";

export const apis = {
	auth: {
		exchangeCode: (data: AuthExchangeRequest) =>
			server.post<APIResponse<AuthExchangeResponse>>("/auth/exchange", data),

		me: () => serverWithInterceptors.get<APIResponse>("/auth/me"),

		spotifyLogin: () => server.get("/auth/spotify/login"),

		spotifyCallback: (params: { code: string; state: string }) =>
			server.get("/auth/spotify/callback", { params }),
	},

	shows: {
		saved: (q?: string) =>
			serverWithInterceptors.get<APIResponse<SavedShowResponse[]>>(
				"/shows/saved",
				{
					params: q ? { q } : undefined,
				},
			),

		search: (params: { q: string; limit?: number; offset?: number }) =>
			serverWithInterceptors.get<APIResponse<SavedShowResponse[]>>(
				"/shows/search",
				{ params },
			),

		subscribe: (spotifyShowId: string) =>
			serverWithInterceptors.post<APIResponse<SubscriptionResponse>>(
				`/shows/${spotifyShowId}/subscribe`,
			),
	},

	subscriptions: {
		list: () =>
			serverWithInterceptors.get<APIResponse<SubscriptionResponse[]>>(
				"/subscriptions",
			),

		delete: (id: string) =>
			serverWithInterceptors.delete<APIResponse>(`/subscriptions/${id}`),
	},

	channels: {
		list: () =>
			serverWithInterceptors.get<APIResponse<NotificationChannel[]>>(
				"/channels",
			),

		create: (data: CreateChannelRequest) =>
			serverWithInterceptors.post<APIResponse<NotificationChannel>>(
				"/channels",
				data,
			),

		toggle: (channelID: string, data: ToggleActiveChannelRequest) =>
			serverWithInterceptors.post<APIResponse>(
				`/channels/${channelID}/toggle`,
				data,
			),

		delete: (channelID: string) =>
			serverWithInterceptors.delete<APIResponse>(`/channels/${channelID}`),
	},

	telegram: {
		generateLink: () =>
			serverWithInterceptors.post<APIResponse<TelegramLinkResponse>>(
				"/telegram/generate-link",
			),
	},

	health: {
		check: () => server.get<HealthResponse>("/health"),
	},
};
