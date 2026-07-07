import { createQuery } from "react-query-kit";

import { apis } from "@/services/api-services";
import type { APIResponse, DashboardSummary } from "@/services/types";

import { dashboardKeys } from "../keys/dashboard.keys";

export const useDashboardSummary = createQuery({
	queryKey: dashboardKeys.summary(),
	fetcher: async (): Promise<APIResponse<DashboardSummary>> => {
		const response = await apis.dashboard.summary();
		return response.data;
	},
});