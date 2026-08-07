import * as Sentry from "@sentry/nextjs";
import { createMutation } from "react-query-kit";
import { queryClient } from "@/lib/query-client";
import { apis } from "@/services/api-services";
import type {
	APIResponse,
	CreateChannelRequest,
	NotificationChannel,
	TelegramLinkResponse,
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
	onSuccess: (_data, variables) => {
		Sentry.metrics.count("podevents.channel.created", 1, {
			attributes: { channel_type: variables.channel_type },
		});
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
		Sentry.metrics.count("podevents.channel.deleted");
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
		const previousIsActive = previousData?.data?.find(
			(channel) => channel.id === variables.channelID,
		)?.is_active;

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

		return { previousIsActive };
	},
	onError: (_err, variables, context) => {
		if (context?.previousIsActive !== undefined) {
			const previousIsActive = context.previousIsActive;
			queryClient.setQueryData<APIResponse<NotificationChannel[]> | undefined>(
				channelKeys.list(),
				(currentData) => {
					if (!currentData?.data) return currentData;

					return {
						...currentData,
						data: currentData.data.map((channel) =>
							channel.id === variables.channelID
								? { ...channel, is_active: previousIsActive }
								: channel,
						),
					};
				},
			);
		}
	},
	onSettled: () => {
		queryClient.invalidateQueries({ queryKey: channelKeys.all });
	},
	onSuccess: (_data, variables) => {
		Sentry.metrics.count("podevents.channel.status_changed", 1, {
			attributes: { active: variables.is_active },
		});
	},
});

export const useGenerateTelegramLink = createMutation<
	APIResponse<TelegramLinkResponse>,
	void,
	Error
>({
	mutationFn: async () => {
		const response = await apis.telegram.generateLink();
		return response.data;
	},
	onSuccess: () => {
		Sentry.metrics.count("podevents.telegram.link_generated");
	},
});
