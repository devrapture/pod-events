"use client";

import { ExternalLink } from "lucide-react";
import pluralize from "pluralize";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { cn, formatNumber } from "@/lib/utils";
import type { SavedShowWithTracking } from "@/services/types";

interface PodcastCardProps {
	show: SavedShowWithTracking;
	isSelected: boolean;
	onToggle: (id: string) => void;
}

export function PodcastCard({ show, isSelected, onToggle }: PodcastCardProps) {
	const isTracked = show.is_tracked;

	return (
		<Card
			className={cn(
				"transition-opacity",
				isTracked && "opacity-60",
				!isTracked && isSelected && "border-emerald-500/30 bg-emerald-500/5",
			)}
		>
			<div className="flex gap-4 p-4">
				<div className="flex shrink-0 items-start pt-1">
					<input
						aria-label={`Select ${show.name}`}
						checked={isSelected}
						className="h-4 w-4 rounded border-white/20 bg-white/5 accent-emerald-500 disabled:cursor-not-allowed disabled:opacity-40"
						disabled={isTracked}
						id={`show-${show.id}`}
						onChange={() => onToggle(show.id)}
						type="checkbox"
					/>
				</div>

				<div className="min-w-0 flex-1">
					<div className="flex gap-3">
						{show.image_url ? (
							// biome-ignore lint/performance/noImgElement: dynamic Spotify image URL
							<img
								alt={`${show.name} cover art`}
								className="h-16 w-16 shrink-0 rounded-lg object-cover"
								src={show.image_url}
							/>
						) : (
							<div className="flex h-16 w-16 shrink-0 items-center justify-center rounded-lg bg-white/5 font-semibold text-emerald-400">
								{show.name.charAt(0).toUpperCase()}
							</div>
						)}

						<div className="min-w-0 flex-1">
							<label
								className={cn(
									"block font-medium text-sm text-zinc-100",
									!isTracked && "cursor-pointer",
								)}
								htmlFor={`show-${show.id}`}
							>
								<span className="line-clamp-1">{show.name}</span>
							</label>
							{show.description && (
								<p className="mt-1 line-clamp-2 text-xs text-zinc-400">
									{show.description}
								</p>
							)}
							<p className="mt-1 text-xs text-zinc-500">
								{formatNumber(show.total_episodes)}{" "}
								{pluralize("episode", show.total_episodes)}
							</p>
							<div className="mt-2 flex flex-wrap items-center gap-2">
								<Badge variant={isTracked ? "emerald" : "outline"}>
									{isTracked ? "Tracked" : "Not tracked"}
								</Badge>
								{show.spotify_url && (
									<Button asChild size="sm" variant="ghost">
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
					</div>
				</div>
			</div>
		</Card>
	);
}
