"use client";

import { ExternalLink, Loader2 } from "lucide-react";
import pluralize from "pluralize";

import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { formatNumber } from "@/lib/utils";
import type { SavedShowWithTracking } from "@/services/types";

interface SearchResultCardProps {
	show: SavedShowWithTracking;
	isTracking: boolean;
	onTrack: (spotifyShowId: string) => void;
}

export function SearchResultCard({
	show,
	isTracking,
	onTrack,
}: SearchResultCardProps) {
	const isTracked = show.is_tracked;

	return (
		<Card>
			<div className="flex flex-col gap-4 p-4 sm:flex-row sm:items-center">
				<div className="flex min-w-0 flex-1 gap-4">
					{show.image_url ? (
						// biome-ignore lint/performance/noImgElement: dynamic Spotify image URL
						<img
							alt={`${show.name} cover art`}
							className="h-20 w-20 shrink-0 rounded-lg object-cover"
							src={show.image_url}
						/>
					) : (
						<div className="flex h-20 w-20 shrink-0 items-center justify-center rounded-lg bg-white/5 font-semibold text-emerald-400">
							{show.name.charAt(0).toUpperCase()}
						</div>
					)}

					<div className="min-w-0 flex-1">
						<h3 className="line-clamp-1 font-medium text-sm text-zinc-100">
							{show.name}
						</h3>
						<p className="mt-1 text-xs text-zinc-500">
							{formatNumber(show.total_episodes)}{" "}
							{pluralize("episode", show.total_episodes)}
						</p>
						{show.description && (
							<p className="mt-1 line-clamp-2 text-xs text-zinc-400">
								{show.description}
							</p>
						)}
						{show.spotify_url && (
							<Button asChild className="mt-2" size="sm" variant="ghost">
								<a
									aria-label={`Open ${show.name} on Spotify`}
									href={show.spotify_url}
									rel="noopener noreferrer"
									target="_blank"
								>
									<ExternalLink className="h-3.5 w-3.5" />
									Spotify
								</a>
							</Button>
						)}
					</div>
				</div>

				<Button
					className="w-full shrink-0 sm:w-auto"
					disabled={isTracked || isTracking}
					onClick={() => onTrack(show.id)}
				>
					{isTracking ? (
						<>
							<Loader2 className="h-4 w-4 animate-spin" />
							Tracking...
						</>
					) : isTracked ? (
						"Tracking"
					) : (
						"Track Podcast"
					)}
				</Button>
			</div>
		</Card>
	);
}
