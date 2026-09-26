import { createQuery } from "react-query-kit";

import { apis } from "@/services/api-services";
import type {
	APIResponse,
	SavedShowResponse,
	ShowSearchParams,
} from "@/services/types";

export const useSearchShows = createQuery({
	queryKey: ["shows", "search"] as const,
	fetcher: async (
		variables: ShowSearchParams,
	): Promise<APIResponse<SavedShowResponse[]>> => {
		const response = await apis.shows.search(variables);
		return response.data;
	},
});
