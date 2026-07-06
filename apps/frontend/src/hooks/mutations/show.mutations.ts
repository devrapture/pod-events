import { createMutation } from "react-query-kit";
import { queryClient } from "@/lib/query-client";
import { apis } from "@/services/api-services";
import type { APIResponse, SubscriptionResponse } from "@/services/types";

import { showKeys } from "../keys/show.keys";
import { subscriptionKeys } from "../keys/subscription.keys";

export const useSubscribeToShow = createMutation({
	mutationFn: async (variables: {
		spotifyShowId: string;
	}): Promise<APIResponse<SubscriptionResponse>> => {
		const response = await apis.shows.subscribe(variables.spotifyShowId);
		return response.data;
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: showKeys.all });
		queryClient.invalidateQueries({ queryKey: subscriptionKeys.all });
	},
});
