import { createQuery } from "react-query-kit";

import { apis } from "@/services/api-services";
import type { APIResponse, TelegramLinkResponse } from "@/services/types";

import { telegramKeys } from "../keys/telegram.keys";

export const useTelegramLink = createQuery({
	queryKey: telegramKeys.link(),
	fetcher: async (): Promise<APIResponse<TelegramLinkResponse>> => {
		const response = await apis.telegram.generateLink();
		return response.data;
	},
	enabled: false,
});
