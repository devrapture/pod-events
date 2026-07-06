import { createQuery } from "react-query-kit";

import { apis } from "@/services/api-services";
import type { APIResponse, SubscriptionResponse } from "@/services/types";

import { subscriptionKeys } from "../keys/subscription.keys";

export const useSubscriptions = createQuery({
	queryKey: subscriptionKeys.list(),
	fetcher: async (): Promise<APIResponse<SubscriptionResponse[]>> => {
		const response = await apis.subscriptions.list();
		return response.data;
	},
});
