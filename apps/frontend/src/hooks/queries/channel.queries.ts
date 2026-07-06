import { createQuery } from "react-query-kit";

import { apis } from "@/services/api-services";
import type { APIResponse, NotificationChannel } from "@/services/types";

import { channelKeys } from "../keys/channel.keys";

export const useChannels = createQuery({
	queryKey: channelKeys.list(),
	fetcher: async (): Promise<APIResponse<NotificationChannel[]>> => {
		const response = await apis.channels.list();
		return response.data;
	},
});
