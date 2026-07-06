import { createQuery } from "react-query-kit";

import { apis } from "@/services/api-services";
import type { HealthResponse } from "@/services/types";

import { healthKeys } from "../keys/health.keys";

export const useHealthCheck = createQuery({
	queryKey: healthKeys.check(),
	fetcher: async (): Promise<HealthResponse> => {
		const response = await apis.health.check();
		return response.data;
	},
	staleTime: 30 * 1000,
	retry: 1,
});
