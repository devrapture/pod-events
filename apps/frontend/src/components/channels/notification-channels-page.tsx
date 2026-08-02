"use client";

import { Bell, Loader2, Plus, RefreshCw } from "lucide-react";
import { useCallback, useState } from "react";

import { AddChannelForm } from "@/components/channels/add-channel-form";
import { ChannelCard } from "@/components/channels/channel-card";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import {
	useDeleteChannel,
	useToggleChannel,
} from "@/hooks/mutations/channel.mutations";
import { useChannels } from "@/hooks/queries/channel.queries";
import { useToast } from "@/hooks/use-toast";
import type { NotificationChannel } from "@/services/types";

function ChannelCardSkeleton() {
	return (
		<Card>
			<div className="animate-pulse p-4">
				<div className="flex gap-3">
					<div className="h-10 w-10 rounded-lg bg-white/10" />
					<div className="flex-1 space-y-2">
						<div className="h-4 w-32 rounded bg-white/10" />
						<div className="h-3 w-48 rounded bg-white/10" />
					</div>
				</div>
			</div>
		</Card>
	);
}

export function NotificationChannelsPage() {
	const [showAddSection, setShowAddSection] = useState(false);
	const [pendingToggle, setPendingToggle] =
		useState<NotificationChannel | null>(null);
	const [pendingDeleteId, setPendingDeleteId] = useState<string | null>(null);
	const { toast } = useToast();

	const { data, isLoading, isFetching, isError, refetch } = useChannels({});

	const { mutate: toggleChannel, isPending: isToggling } = useToggleChannel();
	const { mutate: deleteChannel, isPending: isDeleting } = useDeleteChannel();

	const channels = data?.data ?? [];

	const handleConfirmToggle = useCallback(() => {
		if (!pendingToggle) return;

		const nextActive = !pendingToggle.is_active;
		toggleChannel(
			{ channelID: pendingToggle.id, is_active: nextActive },
			{
				onSuccess: () => {
					toast.success(
						nextActive
							? "Notification channel activated"
							: "Notification channel disabled",
					);
					setPendingToggle(null);
				},
				onError: () => {
					toast.error("Failed to update channel status. Please try again.");
				},
			},
		);
	}, [pendingToggle, toggleChannel, toast]);

	const handleConfirmDelete = useCallback(() => {
		if (!pendingDeleteId) return;

		deleteChannel(
			{ channelID: pendingDeleteId },
			{
				onSuccess: () => {
					toast.success("Notification channel removed");
					setPendingDeleteId(null);
				},
				onError: () => {
					toast.error("Failed to remove channel. Please try again.");
				},
			},
		);
	}, [pendingDeleteId, deleteChannel, toast]);

	return (
		<div className="flex min-h-0 flex-1 flex-col">
			<div className="shrink-0 border-white/6 border-b pb-6">
				<div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
					<div>
						<h1 className="font-semibold text-2xl text-zinc-100 tracking-tight">
							Notification Channels
						</h1>
						<p className="mt-2 text-zinc-400">
							Configure where to receive new episode alerts.
						</p>
						{isFetching && !isLoading && (
							<div className="mt-2 flex items-center gap-2 text-xs text-zinc-500">
								<Loader2 className="h-3 w-3 animate-spin" />
								Refreshing...
							</div>
						)}
					</div>
					<Button
						onClick={() => setShowAddSection((prev) => !prev)}
						variant={showAddSection ? "secondary" : "default"}
					>
						<Plus className="h-4 w-4" />
						Add Channel
					</Button>
				</div>
			</div>

			<div className="mt-6 min-h-0 flex-1 overflow-y-auto">
				{showAddSection && (
					<AddChannelForm onSuccess={() => setShowAddSection(false)} />
				)}

				{isLoading ? (
					<div className="space-y-4">
						{["a", "b", "c"].map((id) => (
							<ChannelCardSkeleton key={`skeleton-${id}`} />
						))}
					</div>
				) : isError ? (
					<div className="flex flex-col items-center justify-center rounded-2xl border border-white/8 bg-white/3 px-6 py-16 text-center">
						<p className="font-medium text-zinc-200">
							Could not load notification channels
						</p>
						<p className="mt-2 text-sm text-zinc-500">
							Check your connection and try again.
						</p>
						<Button
							className="mt-6"
							onClick={() => refetch()}
							variant="secondary"
						>
							<RefreshCw className="h-4 w-4" />
							Try again
						</Button>
					</div>
				) : channels.length === 0 ? (
					<div className="flex flex-col items-center justify-center rounded-2xl border border-white/8 bg-white/3 px-6 py-16 text-center">
						<div className="mb-4 flex h-12 w-12 items-center justify-center rounded-xl border border-emerald-500/20 bg-emerald-500/10">
							<Bell className="h-6 w-6 text-emerald-400" />
						</div>
						<p className="font-medium text-zinc-200">
							No notification channels yet
						</p>
						<p className="mt-2 max-w-md text-sm text-zinc-500">
							Add a channel to start receiving new episode alerts via Slack,
							Discord, WhatsApp, or Telegram.
						</p>
						<Button className="mt-6" onClick={() => setShowAddSection(true)}>
							<Plus className="h-4 w-4" />
							Add Channel
						</Button>
					</div>
				) : (
					<div className="space-y-4">
						{channels.map((channel) => (
							<ChannelCard
								channel={channel}
								key={channel.id}
								onRemove={setPendingDeleteId}
								onToggle={setPendingToggle}
							/>
						))}
					</div>
				)}
			</div>

			<ConfirmDialog
				cancelLabel="Cancel"
				closeOnConfirm={false}
				confirmLabel={pendingToggle?.is_active ? "Disable" : "Activate"}
				description={
					pendingToggle?.is_active
						? "Are you sure you want to disable this notification channel?"
						: "Are you sure you want to activate this notification channel?"
				}
				isPending={isToggling}
				onConfirm={handleConfirmToggle}
				onOpenChange={(open) => {
					if (!open && !isToggling) {
						setPendingToggle(null);
					}
				}}
				open={pendingToggle !== null}
				title={
					pendingToggle?.is_active
						? "Disable notification channel?"
						: "Activate notification channel?"
				}
				variant={pendingToggle?.is_active ? "destructive" : "default"}
			/>

			<ConfirmDialog
				cancelLabel="Cancel"
				closeOnConfirm={false}
				confirmLabel="Remove"
				description="Are you sure you want to remove this notification channel? You will stop receiving alerts through this channel."
				isPending={isDeleting}
				onConfirm={handleConfirmDelete}
				onOpenChange={(open) => {
					if (!open && !isDeleting) {
						setPendingDeleteId(null);
					}
				}}
				open={pendingDeleteId !== null}
				title="Remove notification channel?"
				variant="destructive"
			/>
		</div>
	);
}
