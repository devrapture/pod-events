import { createMutation } from "react-query-kit";
import { queryClient } from "@/lib/query-client";
import { apis } from "@/services/api-services";
import type {
	APIResponse,
	CreateChannelRequest,
	NotificationChannel,
	ToggleActiveChannelRequest,
} from "@/services/types";

import { channelKeys } from "../keys/channel.keys";

export const useCreateChannel = createMutation({
	mutationFn: async (
		variables: CreateChannelRequest,
	): Promise<APIResponse<NotificationChannel>> => {
		const response = await apis.channels.create(variables);
		return response.data;
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: channelKeys.all });
	},
});

export const useDeleteChannel = createMutation({
	mutationFn: async (variables: {
		channelID: string;
	}): Promise<APIResponse> => {
		const response = await apis.channels.delete(variables.channelID);
		return response.data;
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: channelKeys.all });
	},
});

export const useToggleChannel = createMutation({
	mutationFn: async (
		variables: { channelID: string } & ToggleActiveChannelRequest,
	): Promise<APIResponse> => {
		const response = await apis.channels.toggle(variables.channelID, {
			is_active: variables.is_active,
		});
		return response.data;
	},
	onMutate: async (variables) => {
		await queryClient.cancelQueries({ queryKey: channelKeys.all });

		const previousData = queryClient.getQueryData<
			APIResponse<NotificationChannel[]>
		>(channelKeys.list());

		if (previousData?.data) {
			const optimised = previousData.data.map((channel) =>
				channel.id === variables.channelID
					? { ...channel, is_active: variables.is_active }
					: channel,
			);
			queryClient.setQueryData(channelKeys.list(), {
				...previousData,
				data: optimised,
			});
		}

		return { previousData };
	},
	onError: (_err, _vars, context) => {
		if (context?.previousData) {
			queryClient.setQueryData(channelKeys.list(), context.previousData);
		}
	},
	onSettled: () => {
		queryClient.invalidateQueries({ queryKey: channelKeys.all });
	},
});

export const useGenerateTelegramLink = createMutation({
	mutationFn: async (): Promise<APIResponse> => {
		const response = await apis.telegram.generateLink();
		return response.data;
	},
});
