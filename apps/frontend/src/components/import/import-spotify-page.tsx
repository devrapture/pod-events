"use client";

import { Podcast, RefreshCw } from "lucide-react";
import { useCallback, useMemo, useState } from "react";

import { ImportToolbar } from "@/components/import/import-toolbar";
import { PodcastCard } from "@/components/import/podcast-card";
import { SelectionActionBar } from "@/components/import/selection-action-bar";
import { Button } from "@/components/ui/button";
import { useBulkTrackShows } from "@/hooks/mutations/show.mutations";
import { useSavedShows } from "@/hooks/queries/show.queries";
import { useSubscriptions } from "@/hooks/queries/subscription.queries";
import { useDebouncedValue } from "@/hooks/use-debounced-value";
import { useToast } from "@/hooks/use-toast";
import { getAPIErrorMessage } from "@/lib/api-error";
import { enrichShowsWithTracking } from "@/lib/show-utils";
import { cn } from "@/lib/utils";

const MAX_SHOWS_PER_SUBSCRIPTION_REQUEST = 50;

function PodcastCardSkeleton() {
	return (
		<div className="animate-pulse rounded-2xl border border-white/8 bg-white/3 p-4">
			<div className="flex gap-4">
				<div className="h-4 w-4 rounded bg-white/10" />
				<div className="h-16 w-16 rounded-lg bg-white/10" />
				<div className="flex-1 space-y-2">
					<div className="h-4 w-3/4 rounded bg-white/10" />
					<div className="h-3 w-full rounded bg-white/10" />
					<div className="h-3 w-2/3 rounded bg-white/10" />
					<div className="h-3 w-1/3 rounded bg-white/10" />
					<div className="h-5 w-20 rounded-full bg-white/10" />
				</div>
			</div>
		</div>
	);
}

export function ImportSpotifyPage() {
	const [searchQuery, setSearchQuery] = useState("");
	const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
	const debouncedQuery = useDebouncedValue(searchQuery, 300);
	const { toast } = useToast();

	const {
		data: savedData,
		isLoading,
		isFetching,
		isError,
		refetch,
	} = useSavedShows({
		variables: debouncedQuery ? { q: debouncedQuery } : undefined,
	});

	const { data: subsData } = useSubscriptions({});

	const shows = savedData?.data ?? [];
	const enrichedShows = useMemo(
		() => enrichShowsWithTracking(shows, subsData?.data ?? []),
		[shows, subsData?.data],
	);

	const untrackedShows = useMemo(
		() => enrichedShows.filter((show) => !show.is_tracked),
		[enrichedShows],
	);

	const { mutate: trackShows, isPending } = useBulkTrackShows();

	const toggleSelection = useCallback(
		(id: string) => {
			const show = enrichedShows.find((s) => s.id === id);
			if (!show || show.is_tracked) return;

			if (
				!selectedIds.has(id) &&
				selectedIds.size >= MAX_SHOWS_PER_SUBSCRIPTION_REQUEST
			) {
				toast.error(
					`You can track up to ${MAX_SHOWS_PER_SUBSCRIPTION_REQUEST} podcasts at once`,
				);
				return;
			}

			setSelectedIds((prev) => {
				const next = new Set(prev);
				if (next.has(id)) next.delete(id);
				else next.add(id);
				return next;
			});
		},
		[enrichedShows, selectedIds, toast],
	);

	const selectAllUntracked = useCallback(() => {
		const ids = untrackedShows
			.slice(0, MAX_SHOWS_PER_SUBSCRIPTION_REQUEST)
			.map((show) => show.id);

		setSelectedIds(new Set(ids));

		if (untrackedShows.length > MAX_SHOWS_PER_SUBSCRIPTION_REQUEST) {
			toast.error(
				`You can track up to ${MAX_SHOWS_PER_SUBSCRIPTION_REQUEST} podcasts at once`,
			);
		}
	}, [toast, untrackedShows]);

	const clearSelection = useCallback(() => {
		setSelectedIds(new Set());
	}, []);

	const handleTrack = useCallback(() => {
		const ids = Array.from(selectedIds);
		if (ids.length === 0) return;

		trackShows(
			{ spotifyShowIds: ids },
			{
				onSuccess: () => {
					toast.success(
						`Tracked ${ids.length} podcast${ids.length === 1 ? "" : "s"}`,
					);
					setSelectedIds(new Set());
				},
				onError: (error) => {
					toast.error(
						getAPIErrorMessage(
							error,
							"Failed to track podcasts. Please try again.",
						),
					);
				},
			},
		);
	}, [selectedIds, trackShows, toast]);

	const hasSelection = selectedIds.size > 0;

	return (
		<div className="flex min-h-0 flex-1 flex-col">
			<div className="shrink-0 border-white/6 border-b pb-6">
				<div className="mb-6">
					<h1 className="font-semibold text-2xl text-zinc-100 tracking-tight">
						Bring Your Spotify Podcasts Home
					</h1>
					<p className="mt-2 max-w-2xl text-zinc-400">
						These are the podcasts you already follow on Spotify. Select any to
						start tracking new episodes from PodEvents.
					</p>
				</div>

				<ImportToolbar
					isSearching={isFetching && !isLoading}
					onClearSelection={clearSelection}
					onSearchChange={setSearchQuery}
					onSelectAllUntracked={selectAllUntracked}
					searchQuery={searchQuery}
					selectedCount={selectedIds.size}
					untrackedCount={untrackedShows.length}
				/>
			</div>

			<div
				className={cn(
					"mt-6 min-h-0 flex-1 overflow-y-auto",
					hasSelection && "pb-24",
				)}
			>
				{isLoading ? (
					<div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
						{["a", "b", "c", "d", "e", "f"].map((id) => (
							<PodcastCardSkeleton key={`skeleton-${id}`} />
						))}
					</div>
				) : isError ? (
					<div className="flex flex-col items-center justify-center rounded-2xl border border-white/8 bg-white/3 px-6 py-16 text-center">
						<p className="font-medium text-zinc-200">
							Could not load your saved podcasts
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
				) : enrichedShows.length === 0 ? (
					<div className="flex flex-col items-center justify-center rounded-2xl border border-white/8 bg-white/3 px-6 py-16 text-center">
						<div className="mb-4 flex h-12 w-12 items-center justify-center rounded-xl border border-emerald-500/20 bg-emerald-500/10">
							<Podcast className="h-6 w-6 text-emerald-400" />
						</div>
						{debouncedQuery ? (
							<>
								<p className="font-medium text-zinc-200">No results found</p>
								<p className="mt-2 text-sm text-zinc-500">
									No saved podcasts match &ldquo;{debouncedQuery}&rdquo;. Try a
									different search term.
								</p>
							</>
						) : (
							<>
								<p className="font-medium text-zinc-200">
									No saved podcasts found
								</p>
								<p className="mt-2 text-sm text-zinc-500">
									Follow podcasts on Spotify first, then come back to import
									them here.
								</p>
							</>
						)}
					</div>
				) : (
					<div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
						{enrichedShows.map((show) => (
							<PodcastCard
								isSelected={selectedIds.has(show.id)}
								key={show.id}
								onToggle={toggleSelection}
								show={show}
							/>
						))}
					</div>
				)}
			</div>

			<SelectionActionBar
				isPending={isPending}
				onTrack={handleTrack}
				selectedCount={selectedIds.size}
			/>
		</div>
	);
}
