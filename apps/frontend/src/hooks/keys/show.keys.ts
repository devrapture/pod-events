import type { SavedShowsParams, ShowSearchParams } from "@/services/types";

export const showKeys = {
	all: ["shows"] as const,
	saved: (params?: SavedShowsParams) =>
		params
			? (["shows", "saved", params] as const)
			: (["shows", "saved"] as const),
	search: (params: ShowSearchParams) => ["shows", "search", params] as const,
};
