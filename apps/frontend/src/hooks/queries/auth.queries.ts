import { createQuery } from "react-query-kit";

import { apis } from "@/services/api-services";
import type { APIResponse } from "@/services/types";

import { authKeys } from "../keys/auth.keys";

export const useCurrentUser = createQuery({
	queryKey: authKeys.me(),
	fetcher: async (): Promise<APIResponse> => {
		const response = await apis.auth.me();
		return response.data;
	},
});
