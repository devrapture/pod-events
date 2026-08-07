import * as Sentry from "@sentry/nextjs";
import { createMutation } from "react-query-kit";
import { queryClient } from "@/lib/query-client";
import { apis } from "@/services/api-services";
import type { APIResponse } from "@/services/types";

import { dashboardKeys } from "../keys/dashboard.keys";
import { showKeys } from "../keys/show.keys";
import { subscriptionKeys } from "../keys/subscription.keys";

export const useSubscribeToShow = createMutation({
	mutationFn: async (variables: {
		spotifyShowId: string;
	}): Promise<APIResponse> => {
		const response = await apis.shows.subscribe([variables.spotifyShowId]);
		return response.data;
	},
	onSuccess: () => {
		Sentry.metrics.count("podevents.subscription.created", 1, {
			attributes: { source: "search" },
		});
		queryClient.invalidateQueries({ queryKey: showKeys.all });
		queryClient.invalidateQueries({ queryKey: subscriptionKeys.all });
		queryClient.invalidateQueries({
			queryKey: dashboardKeys.all,
			refetchType: "all",
		});
	},
});

export const useBulkTrackShows = createMutation({
	mutationFn: async (variables: {
		spotifyShowIds: string[];
	}): Promise<APIResponse> => {
		const response = await apis.shows.subscribe(variables.spotifyShowIds);
		return response.data;
	},
	onSuccess: (_data, variables) => {
		Sentry.metrics.count(
			"podevents.subscription.created",
			variables.spotifyShowIds.length,
			{ attributes: { source: "spotify_import" } },
		);
		queryClient.invalidateQueries({ queryKey: showKeys.all });
		queryClient.invalidateQueries({ queryKey: subscriptionKeys.all });
		queryClient.invalidateQueries({
			queryKey: dashboardKeys.all,
			refetchType: "all",
		});
	},
});
