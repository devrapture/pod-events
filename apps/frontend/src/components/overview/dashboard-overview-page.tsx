"use client";

import { RefreshCw } from "lucide-react";

import { OverviewHeader } from "@/components/overview/overview-header";
import { RecentEpisodes } from "@/components/overview/recent-episodes";
import { SetupChecklist } from "@/components/overview/setup-checklist";
import { StatsCards } from "@/components/overview/stats-cards";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { useDashboardSummary } from "@/hooks/queries/dashboard.queries";

const SETUP_CHECKLIST_SKELETON_KEYS = [
	"connect-spotify",
	"import-spotify",
	"subscribe-podcast",
	"add-channel",
	"receive-notification",
] as const;

const STATS_CARDS_SKELETON_KEYS = [
	"podcasts-tracked",
	"active-channels",
	"new-episodes",
	"notifications-sent",
] as const;

function SetupChecklistSkeleton() {
	return (
		<Card>
			<CardContent className="space-y-6 p-6">
				<div className="space-y-2">
					<div className="h-5 w-36 animate-pulse rounded bg-white/10" />
					<div className="h-4 w-28 animate-pulse rounded bg-white/10" />
				</div>
				<div className="h-2 animate-pulse rounded-full bg-white/10" />
				<div className="space-y-3">
					{SETUP_CHECKLIST_SKELETON_KEYS.map((key) => (
						<div className="flex items-center gap-3" key={key}>
							<div className="h-5 w-5 animate-pulse rounded-full bg-white/10" />
							<div className="h-4 flex-1 animate-pulse rounded bg-white/10" />
						</div>
					))}
				</div>
			</CardContent>
		</Card>
	);
}

function StatsCardsSkeleton() {
	return (
		<div className="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
			{STATS_CARDS_SKELETON_KEYS.map((key) => (
				<Card key={key}>
					<CardContent className="space-y-2 p-6">
						<div className="h-4 w-24 animate-pulse rounded bg-white/10" />
						<div className="h-8 w-16 animate-pulse rounded bg-white/10" />
					</CardContent>
				</Card>
			))}
		</div>
	);
}

export function DashboardOverviewPage() {
	const { data, isLoading, isError, refetch, isFetching } = useDashboardSummary(
		{},
	);

	const summary = data?.data;

	return (
		<div className="mx-auto w-[95%]">
			<OverviewHeader />

			{isLoading && (
				<>
					<SetupChecklistSkeleton />
					<StatsCardsSkeleton />
				</>
			)}

			{isError && !isLoading && (
				<Card>
					<CardContent className="flex flex-col items-center gap-4 p-8 text-center">
						<p className="text-zinc-400">
							Failed to load dashboard summary. Please try again.
						</p>
						<Button onClick={() => refetch()} variant="outline">
							<RefreshCw className={isFetching ? "animate-spin" : undefined} />
							Retry
						</Button>
					</CardContent>
				</Card>
			)}

			{summary && (
				<>
					{summary.setup.percent < 100 && (
						<SetupChecklist setup={summary.setup} />
					)}
					<StatsCards stats={summary.stats} />
					<RecentEpisodes episodes={summary.recent_episodes ?? []} />
				</>
			)}
		</div>
	);
}
