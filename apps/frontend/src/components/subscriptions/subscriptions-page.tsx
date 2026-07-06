"use client";

import { Loader2, Podcast, RefreshCw } from "lucide-react";
import Link from "next/link";
import { useCallback, useState } from "react";
import pluralize from "pluralize";

import { SubscriptionCard } from "@/components/subscriptions/subscription-card";
import { Button } from "@/components/ui/button";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import { useUnsubscribe } from "@/hooks/mutations/subscription.mutations";
import { useSubscriptions } from "@/hooks/queries/subscription.queries";
import { useToast } from "@/hooks/use-toast";
import { formatNumber } from "@/lib/utils";

function SubscriptionCardSkeleton() {
	return (
		<div className="animate-pulse rounded-2xl border border-white/8 bg-white/3 p-4">
			<div className="flex gap-4">
				<div className="h-16 w-16 rounded-lg bg-white/10" />
				<div className="flex-1 space-y-2">
					<div className="h-4 w-3/4 rounded bg-white/10" />
					<div className="h-3 w-full rounded bg-white/10" />
					<div className="h-3 w-2/3 rounded bg-white/10" />
					<div className="h-8 w-28 rounded-lg bg-white/10" />
				</div>
			</div>
		</div>
	);
}

function getTrackingSubtitle(count: number): string {
	if (count === 0) {
		return "You are not tracking any podcasts yet.";
	}
	return `You are tracking ${formatNumber(count)} ${pluralize("podcast", count)}.`;
}

export function SubscriptionsPage() {
	const [pendingId, setPendingId] = useState<string | null>(null);
	const { toast } = useToast();

	const {
		data,
		isLoading,
		isFetching,
		isError,
		refetch,
	} = useSubscriptions({});

	const { mutate: unsubscribe, isPending: isUnsubscribing } = useUnsubscribe();

	const subscriptions = data?.data ?? [];
	const count = subscriptions.length;

	const handleConfirmUnsubscribe = useCallback(() => {
		if (!pendingId) return;

		unsubscribe(
			{ id: pendingId },
			{
				onSuccess: () => {
					toast.success("Unsubscribed successfully");
					setPendingId(null);
				},
				onError: () => {
					toast.error("Failed to unsubscribe. Please try again.");
				},
			},
		);
	}, [pendingId, unsubscribe, toast]);

	return (
		<div className="flex min-h-0 flex-1 flex-col">
			<div className="shrink-0 border-white/6 border-b pb-6">
				<h1 className="font-semibold text-2xl text-zinc-100 tracking-tight">
					My Subscriptions
				</h1>
				<p className="mt-2 text-zinc-400">{getTrackingSubtitle(count)}</p>
				{isFetching && !isLoading && (
					<div className="mt-2 flex items-center gap-2 text-xs text-zinc-500">
						<Loader2 className="h-3 w-3 animate-spin" />
						Refreshing...
					</div>
				)}
			</div>

			<div className="mt-6 min-h-0 flex-1 overflow-y-auto">
				{isLoading ? (
					<div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
						{["a", "b", "c", "d", "e", "f"].map((id) => (
							<SubscriptionCardSkeleton key={`skeleton-${id}`} />
						))}
					</div>
				) : isError ? (
					<div className="flex flex-col items-center justify-center rounded-2xl border border-white/8 bg-white/3 px-6 py-16 text-center">
						<p className="font-medium text-zinc-200">
							Could not load your subscriptions
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
				) : subscriptions.length === 0 ? (
					<div className="flex flex-col items-center justify-center rounded-2xl border border-white/8 bg-white/3 px-6 py-16 text-center">
						<div className="mb-4 flex h-12 w-12 items-center justify-center rounded-xl border border-emerald-500/20 bg-emerald-500/10">
							<Podcast className="h-6 w-6 text-emerald-400" />
						</div>
						<p className="font-medium text-zinc-200">No subscriptions yet</p>
						<p className="mt-2 max-w-md text-sm text-zinc-500">
							You&apos;re not tracking any podcasts yet.{" "}
							<Link
								className="text-emerald-400 hover:text-emerald-300"
								href="/dashboard/import"
							>
								Import podcasts from Spotify
							</Link>{" "}
							or{" "}
							<Link
								className="text-emerald-400 hover:text-emerald-300"
								href="/dashboard/search"
							>
								search for podcasts
							</Link>{" "}
							to start tracking them.
						</p>
					</div>
				) : (
					<div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
						{subscriptions.map((subscription) => (
							<SubscriptionCard
								key={subscription.id}
								onUnsubscribe={setPendingId}
								subscription={subscription}
							/>
						))}
					</div>
				)}
			</div>

			<ConfirmDialog
				cancelLabel="Cancel"
				closeOnConfirm={false}
				confirmLabel="Unsubscribe"
				description="Are you sure you want to unsubscribe from this podcast? You will stop receiving notifications for new episodes."
				isPending={isUnsubscribing}
				onConfirm={handleConfirmUnsubscribe}
				onOpenChange={(open) => {
					if (!open && !isUnsubscribing) {
						setPendingId(null);
					}
				}}
				open={pendingId !== null}
				title="Unsubscribe from this podcast?"
				variant="destructive"
			/>
		</div>
	);
}