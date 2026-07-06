import { createQuery } from "react-query-kit";

import { apis } from "@/services/api-services";
import type {
	APIResponse,
	SavedShowResponse,
	SavedShowsParams,
	ShowSearchParams,
} from "@/services/types";

import { showKeys } from "../keys/show.keys";

export const useSavedShows = createQuery({
	queryKey: showKeys.saved(),
	fetcher: async (
		variables: SavedShowsParams | undefined,
	): Promise<APIResponse<SavedShowResponse[]>> => {
		const response = await apis.shows.saved(variables?.q);
		return response.data;
	},
});

export const useSearchShows = createQuery({
	queryKey: ["shows", "search"] as const,
	fetcher: async (
		variables: ShowSearchParams,
	): Promise<APIResponse<SavedShowResponse[]>> => {
		const response = await apis.shows.search(variables);
		return response.data;
	},
});
