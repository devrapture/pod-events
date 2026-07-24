"use client";

import { ExternalLink } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import type { SubscriptionResponse } from "@/services/types";

interface SubscriptionCardProps {
	subscription: SubscriptionResponse;
	onUnsubscribe: (id: string) => void;
}

export function SubscriptionCard({
	subscription,
	onUnsubscribe,
}: SubscriptionCardProps) {
	const show = subscription.podcast_show;

	return (
		<Card>
			<div className="flex gap-4 p-4">
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
					<h3 className="line-clamp-1 font-medium text-sm text-zinc-100">
						{show.name}
					</h3>
					{show.description && (
						<p className="mt-1 line-clamp-2 text-xs text-zinc-400">
							{show.description}
						</p>
					)}
					<div className="mt-3 flex flex-wrap items-center gap-2">
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
						<Button
							onClick={() => onUnsubscribe(subscription.id)}
							size="sm"
							variant="secondary"
						>
							Unsubscribe
						</Button>
					</div>
				</div>
			</div>
		</Card>
	);
}
