import type { AxiosError } from "axios";

import { server, serverWithInterceptors } from "@/lib/axios-util";

import type {
	APIResponse,
	AuthExchangeRequest,
	AuthExchangeResponse,
	BulkSubscribeItemResult,
	BulkSubscribeResponse,
	CreateChannelRequest,
	HealthResponse,
	NotificationChannel,
	SavedShowResponse,
	SubscriptionResponse,
	TelegramLinkResponse,
	ToggleActiveChannelRequest,
} from "./types";

async function subscribeShow(
	spotifyShowId: string,
): Promise<BulkSubscribeItemResult> {
	try {
		await serverWithInterceptors.post<APIResponse<SubscriptionResponse>>(
			`/shows/${spotifyShowId}/subscribe`,
		);
		return { spotify_show_id: spotifyShowId, success: true };
	} catch (error) {
		const axiosError = error as AxiosError<APIResponse>;
		if (axiosError.response?.status === 409) {
			return { spotify_show_id: spotifyShowId, success: true };
		}
		return {
			spotify_show_id: spotifyShowId,
			success: false,
			error:
				axiosError.response?.data?.error?.message ??
				"Failed to subscribe to show",
		};
	}
}

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

		// TODO: swap to POST /shows/subscribe/bulk when backend ships
		bulkSubscribe: async (
			spotifyShowIds: string[],
		): Promise<{ data: APIResponse<BulkSubscribeResponse> }> => {
			const results = await Promise.all(
				spotifyShowIds.map((id) => subscribeShow(id)),
			);
			const succeeded = results.filter((r) => r.success).length;
			const failed = results.length - succeeded;
			return {
				data: {
					success: true,
					data: { succeeded, failed, results },
				},
			};
		},
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
