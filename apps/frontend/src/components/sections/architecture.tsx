import { ArrowDown, Database, KeyRound, Layers, ShieldCheck } from "lucide-react";

import { AnimatedSection } from "@/components/animated-section";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";

const LABELS = [
	"OAuth",
	"PostgreSQL",
	"Cron Scheduling",
	"Notification Queue",
	"AES Encryption",
	"Idempotency",
] as const;

const CHANNELS = ["Slack", "Discord", "Telegram", "WhatsApp"] as const;

function FlowNode({
	label,
	sub,
	accent = false,
}: {
	label: string;
	sub?: string;
	accent?: boolean;
}) {
	return (
		<Card
			className={`px-6 py-4 text-center transition-all duration-300 hover:border-white/15 ${
				accent
					? "border-emerald-500/30 bg-emerald-500/[0.06] shadow-lg shadow-emerald-500/10"
					: ""
			}`}
		>
			<p className="font-semibold text-zinc-100">{label}</p>
			{sub && <p className="mt-1 text-xs text-zinc-500">{sub}</p>}
		</Card>
	);
}

function FlowArrow() {
	return (
		<div className="flex justify-center py-2">
			<ArrowDown className="h-4 w-4 text-zinc-600" />
		</div>
	);
}

export function Architecture() {
	return (
		<AnimatedSection className="py-24 md:py-32" id="architecture">
			<div className="mx-auto max-w-6xl px-6">
				<div className="mx-auto max-w-2xl text-center">
					<Badge className="mb-4" variant="emerald">
						Architecture
					</Badge>
					<h2 className="font-semibold text-3xl text-zinc-50 tracking-tight md:text-4xl">
						From Spotify to your channels
					</h2>
					<p className="mt-4 text-zinc-400">
						A reliable pipeline designed for production workloads.
					</p>
				</div>

				<div className="mx-auto mt-14 max-w-md">
					<FlowNode label="Spotify" sub="Podcast catalog & episodes API" />
					<FlowArrow />
					<FlowNode accent label="PodEvents" sub="Go API · OAuth · Middleware" />
					<FlowArrow />
					<FlowNode label="Episode Detection" sub="Cron scheduler · Idempotency" />
					<FlowArrow />
					<FlowNode label="Notification Engine" sub="Multi-channel dispatch" />
					<FlowArrow />
					<div className="grid grid-cols-2 gap-2">
						{CHANNELS.map((ch) => (
							<Card
								className="px-4 py-3 text-center text-sm text-zinc-300 transition-all hover:border-white/15 hover:bg-white/[0.04]"
								key={ch}
							>
								{ch}
							</Card>
						))}
					</div>
				</div>

				<div className="mt-12 flex flex-wrap justify-center gap-2">
					{LABELS.map((label) => (
						<span
							className="inline-flex items-center gap-1.5 rounded-full border border-white/10 bg-white/[0.03] px-3 py-1.5 text-xs text-zinc-400"
							key={label}
						>
							{label === "OAuth" && <KeyRound className="h-3 w-3 text-emerald-400" />}
							{label === "PostgreSQL" && <Database className="h-3 w-3 text-emerald-400" />}
							{label === "AES Encryption" && <ShieldCheck className="h-3 w-3 text-emerald-400" />}
							{label === "Notification Queue" && <Layers className="h-3 w-3 text-emerald-400" />}
							{label}
						</span>
					))}
				</div>

				<div className="mx-auto mt-10 max-w-2xl overflow-hidden rounded-xl border border-white/10 bg-zinc-950/80">
					<div className="flex items-center gap-2 border-white/5 border-b px-4 py-3">
						<div className="h-2.5 w-2.5 rounded-full bg-red-500/80" />
						<div className="h-2.5 w-2.5 rounded-full bg-yellow-500/80" />
						<div className="h-2.5 w-2.5 rounded-full bg-green-500/80" />
						<span className="ml-2 font-mono text-xs text-zinc-500">cron/episode-check</span>
					</div>
					<pre className="overflow-x-auto p-4 font-mono text-xs text-zinc-400 leading-relaxed">
{`POST /api/v1/cron/check-episodes
Authorization: Bearer <cron-secret>

→ Fetch saved shows from PostgreSQL
→ Query Spotify API for new episodes
→ Deduplicate via idempotency keys
→ Dispatch to notification channels`}
					</pre>
				</div>
			</div>
		</AnimatedSection>
	);
}