"use client";

import { Check } from "lucide-react";
import Link from "next/link";

import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import { LOGIN_URL } from "@/lib/auth";
import { cn } from "@/lib/utils";
import type { DashboardSetup } from "@/services/types";

const CHECKLIST_LINKS: Record<string, string | null> = {
	connect_spotify: LOGIN_URL,
	import_spotify: "/dashboard/import",
	subscribe_podcast: "/dashboard/search",
	add_channel: "/dashboard/search",
	receive_notification: null,
};

const CHECKLIST_CTA_LABELS: Record<string, string> = {
	connect_spotify: "Connect Spotify",
	import_spotify: "Import from Spotify",
	subscribe_podcast: "Subscribe to a podcast",
	add_channel: "Add a notification channel",
};

interface SetupChecklistProps {
	setup: DashboardSetup;
}

export function SetupChecklist({ setup }: SetupChecklistProps) {
	const nextItem = setup.items.find((item) => !item.completed);
	const ctaHref = nextItem ? CHECKLIST_LINKS[nextItem.key] : null;
	const ctaLabel = nextItem ? CHECKLIST_CTA_LABELS[nextItem.key] : null;

	return (
		<Card>
			<CardHeader>
				<CardTitle>Setup checklist</CardTitle>
				<CardDescription>
					{setup.completed} of {setup.total} completed
				</CardDescription>
			</CardHeader>
			<CardContent className="space-y-6">
				<div className="flex items-center gap-3">
					<div className="h-2 flex-1 overflow-hidden rounded-full bg-white/10">
						<div
							className="h-full rounded-full bg-emerald-500 transition-all duration-500"
							style={{ width: `${setup.percent}%` }}
						/>
					</div>
					<span className="shrink-0 font-medium text-emerald-400 text-sm">
						{setup.percent}%
					</span>
				</div>

				<ul className="space-y-3">
					{setup.items.map((item) => (
						<li key={item.key} className="flex items-center gap-3">
							<span
								className={cn(
									"flex h-6 w-6 shrink-0 items-center justify-center rounded-full",
									item.completed
										? "bg-emerald-500"
										: "border border-zinc-600 bg-transparent",
								)}
							>
								{item.completed && (
									<Check className="h-3.5 w-3.5 text-white" strokeWidth={3} />
								)}
							</span>
							<span
								className={cn(
									"text-sm",
									item.completed ? "text-zinc-300" : "text-zinc-400",
								)}
							>
								{item.label}
							</span>
						</li>
					))}
				</ul>

				{ctaHref && ctaLabel &&
					(nextItem?.key === "connect_spotify" ? (
						<a
							href={ctaHref}
							className="inline-flex items-center gap-1 font-medium text-emerald-400 text-sm transition-colors hover:text-emerald-300"
						>
							{ctaLabel} →
						</a>
					) : (
						<Link
							href={ctaHref}
							className="inline-flex items-center gap-1 font-medium text-emerald-400 text-sm transition-colors hover:text-emerald-300"
						>
							{ctaLabel} →
						</Link>
					))}
			</CardContent>
		</Card>
	);
}