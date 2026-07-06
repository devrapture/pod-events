import { createMutation } from "react-query-kit";
import { queryClient } from "@/lib/query-client";
import { apis } from "@/services/api-services";
import type { APIResponse } from "@/services/types";

import { subscriptionKeys } from "../keys/subscription.keys";

export const useUnsubscribe = createMutation({
	mutationFn: async (variables: { id: string }): Promise<APIResponse> => {
		const response = await apis.subscriptions.delete(variables.id);
		return response.data;
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: subscriptionKeys.all });
	},
});
