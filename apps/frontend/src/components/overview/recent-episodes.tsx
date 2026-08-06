import { ExternalLink } from "lucide-react";

import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import type { RecentEpisode } from "@/services/types";

interface RecentEpisodesProps {
	episodes: RecentEpisode[];
}

function formatDuration(durationMs: number): string {
	const totalMinutes = Math.round(durationMs / 60_000);
	const hours = Math.floor(totalMinutes / 60);
	const minutes = totalMinutes % 60;

	return hours > 0 ? `${hours}h ${minutes}m` : `${minutes}m`;
}

function formatReleaseDate(releaseDate: string): string {
	const date = new Date(releaseDate);
	if (Number.isNaN(date.getTime())) return releaseDate;

	return new Intl.DateTimeFormat("en-US", {
		dateStyle: "medium",
		timeZone: "UTC",
	}).format(date);
}

export function RecentEpisodes({ episodes }: RecentEpisodesProps) {
	return (
		<Card className="mt-8">
			<CardHeader>
				<CardTitle>Recent episodes</CardTitle>
				<CardDescription>
					The latest episodes successfully delivered to your channels.
				</CardDescription>
			</CardHeader>
			<CardContent>
				{episodes.length === 0 ? (
					<p className="text-sm text-zinc-500">
						Delivered episodes will appear here.
					</p>
				) : (
					<ul className="divide-y divide-white/8">
						{episodes.map((episode) => (
							<li className="py-4 first:pt-0 last:pb-0" key={episode.url}>
								<a
									className="group block"
									href={episode.url}
									rel="noreferrer"
									target="_blank"
								>
									<span className="flex items-center gap-2 font-medium text-sm text-zinc-200 transition-colors group-hover:text-emerald-400">
										{episode.episode_title}
										<ExternalLink className="h-3.5 w-3.5 shrink-0" />
									</span>
									<span className="mt-1 block text-xs text-zinc-500">
										{episode.podcast_name} ·{" "}
										{formatReleaseDate(episode.release_date)} ·{" "}
										{formatDuration(episode.duration)}
									</span>
								</a>
							</li>
						))}
					</ul>
				)}
			</CardContent>
		</Card>
	);
}
