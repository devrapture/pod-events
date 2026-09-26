"use client";

import { Loader2, RefreshCw, Search, X } from "lucide-react";
import pluralize from "pluralize";
import { useCallback, useMemo, useState } from "react";

import { SearchResultCard } from "@/components/search/search-result-card";
import { Button } from "@/components/ui/button";
import { useSubscribeToShow } from "@/hooks/mutations/show.mutations";
import { useSearchShows } from "@/hooks/queries/show.queries";
import { useSubscriptions } from "@/hooks/queries/subscription.queries";
import { useDebouncedValue } from "@/hooks/use-debounced-value";
import { useToast } from "@/hooks/use-toast";
import { getAPIErrorMessage } from "@/lib/api-error";
import { enrichShowsWithTracking } from "@/lib/show-utils";
import { cn, formatNumber } from "@/lib/utils";

function SearchResultSkeleton() {
	return (
		<div className="animate-pulse rounded-2xl border border-white/8 bg-white/3 p-4">
			<div className="flex gap-4">
				<div className="h-20 w-20 rounded-lg bg-white/10" />
				<div className="flex-1 space-y-2">
					<div className="h-4 w-3/4 rounded bg-white/10" />
					<div className="h-3 w-1/4 rounded bg-white/10" />
					<div className="h-3 w-full rounded bg-white/10" />
					<div className="h-3 w-2/3 rounded bg-white/10" />
				</div>
				<div className="hidden h-10 w-32 rounded-lg bg-white/10 sm:block" />
			</div>
		</div>
	);
}

function getResultCountText(count: number, query: string): string {
	return `${formatNumber(count)} ${pluralize("result", count)} for "${query}"`;
}

export function SearchPodcastsPage() {
	const [searchQuery, setSearchQuery] = useState("");
	const [trackingId, setTrackingId] = useState<string | null>(null);
	const debouncedQuery = useDebouncedValue(searchQuery, 300);
	const trimmedQuery = debouncedQuery.trim();
	const { toast } = useToast();

	const { data, isLoading, isFetching, isError, refetch } = useSearchShows({
		variables: { q: trimmedQuery, limit: 10 },
		enabled: trimmedQuery.length > 0,
	});

	const { data: subsData } = useSubscriptions({});

	const results = useMemo(
		() => enrichShowsWithTracking(data?.data ?? [], subsData?.data ?? []),
		[data?.data, subsData?.data],
	);

	const { mutate: subscribe, isPending: isSubscribing } = useSubscribeToShow();

	const handleTrack = useCallback(
		(spotifyShowId: string) => {
			const show = results.find((s) => s.id === spotifyShowId);
			setTrackingId(spotifyShowId);

			subscribe(
				{ spotifyShowId },
				{
					onSuccess: () => {
						toast.success(
							show
								? `Now tracking ${show.name}`
								: "Podcast tracked successfully",
						);
						setTrackingId(null);
					},
					onError: (error) => {
						toast.error(
							getAPIErrorMessage(
								error,
								"Failed to track podcast. Please try again.",
							),
						);
						setTrackingId(null);
					},
				},
			);
		},
		[results, subscribe, toast],
	);

	const hasQuery = trimmedQuery.length > 0;
	const isSearching =
		hasQuery && (isLoading || (isFetching && results.length === 0));

	return (
		<div className="flex min-h-0 flex-1 flex-col">
			<div className="shrink-0 border-white/6 border-b pb-6">
				<h1 className="font-semibold text-2xl text-zinc-100 tracking-tight">
					Search Podcasts
				</h1>
				<p className="mt-2 text-zinc-400">
					Discover podcasts from the Spotify catalog and start tracking them.
				</p>

				<div className="relative mt-6 max-w-2xl">
					<Search
						aria-hidden
						className="pointer-events-none absolute top-1/2 left-4 h-5 w-5 -translate-y-1/2 text-zinc-500"
					/>
					<input
						aria-label="Search podcasts"
						className={cn(
							"h-12 w-full rounded-xl border border-white/10 bg-white/5 pr-12 pl-12 text-sm text-zinc-100",
							"placeholder:text-zinc-500 focus:border-emerald-500/40 focus:outline-none focus:ring-2 focus:ring-emerald-500/20",
						)}
						onChange={(e) => setSearchQuery(e.target.value)}
						placeholder='Try "Tech", "Productivity", or any podcast name'
						type="search"
						value={searchQuery}
					/>
					{searchQuery && (
						<button
							aria-label="Clear search"
							className="absolute top-1/2 right-4 -translate-y-1/2 text-zinc-500 hover:text-zinc-300"
							onClick={() => setSearchQuery("")}
							type="button"
						>
							<X className="h-4 w-4" />
						</button>
					)}
				</div>

				{hasQuery && (
					<div className="mt-3 text-sm text-zinc-500">
						{isSearching ? (
							<span className="flex items-center gap-2">
								<Loader2 className="h-3.5 w-3.5 animate-spin" />
								Searching...
							</span>
						) : isError ? null : (
							getResultCountText(results.length, trimmedQuery)
						)}
					</div>
				)}
			</div>

			<div className="mt-6 min-h-0 flex-1 overflow-y-auto">
				{!hasQuery ? (
					<div className="flex flex-col items-center justify-center rounded-2xl border border-white/8 bg-white/3 px-6 py-16 text-center">
						<div className="mb-4 flex h-12 w-12 items-center justify-center rounded-xl border border-emerald-500/20 bg-emerald-500/10">
							<Search className="h-6 w-6 text-emerald-400" />
						</div>
						<p className="font-medium text-zinc-200">Search for a podcast</p>
						<p className="mt-2 max-w-md text-sm text-zinc-500">
							Try &ldquo;tech&rdquo;, &ldquo;Productivity&rdquo;,
							&ldquo;Entrepreneurship&rdquo;, or any podcast name.
						</p>
					</div>
				) : isLoading ? (
					<div className="flex flex-col gap-3">
						{["a", "b", "c", "d", "e", "f"].map((id) => (
							<SearchResultSkeleton key={`skeleton-${id}`} />
						))}
					</div>
				) : isError ? (
					<div className="flex flex-col items-center justify-center rounded-2xl border border-white/8 bg-white/3 px-6 py-16 text-center">
						<p className="font-medium text-zinc-200">
							Could not search podcasts
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
				) : results.length === 0 ? (
					<div className="flex flex-col items-center justify-center rounded-2xl border border-white/8 bg-white/3 px-6 py-16 text-center">
						<p className="font-medium text-zinc-200">No podcasts found</p>
						<p className="mt-2 text-sm text-zinc-500">
							Try a different search term.
						</p>
					</div>
				) : (
					<div className="flex flex-col gap-3">
						{results.map((show) => (
							<SearchResultCard
								isTracking={isSubscribing && trackingId === show.id}
								key={show.id}
								onTrack={handleTrack}
								show={show}
							/>
						))}
					</div>
				)}
			</div>
		</div>
	);
}
