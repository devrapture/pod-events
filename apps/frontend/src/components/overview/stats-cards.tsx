"use client";

import { Card, CardContent } from "@/components/ui/card";
import { formatNumber } from "@/lib/utils";
import type { DashboardStats } from "@/services/types";

const STAT_ITEMS = [
	{ key: "podcasts_tracked", label: "Podcasts tracked" },
	{ key: "active_channels", label: "Active channels" },
	{ key: "new_episodes_this_week", label: "New episodes this week" },
	{ key: "notification_sent", label: "Notifications sent" },
] as const satisfies ReadonlyArray<{
	key: keyof DashboardStats;
	label: string;
}>;

interface StatsCardsProps {
	stats: DashboardStats;
}

export function StatsCards({ stats }: StatsCardsProps) {
	return (
		<div className="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
			{STAT_ITEMS.map((item) => (
				<Card key={item.key}>
					<CardContent className="p-6">
						<p className="text-sm text-zinc-400">{item.label}</p>
						<p className="mt-2 font-semibold text-2xl text-zinc-100">
							{formatNumber(stats[item.key])}
						</p>
					</CardContent>
				</Card>
			))}
		</div>
	);
}
