import { ArrowDown, Bell, Headphones, ListMusic, RefreshCw } from "lucide-react";

import {
	AnimatedSection,
	StaggerGrid,
	StaggerItem,
} from "@/components/animated-section";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";

const STEPS = [
	{
		icon: Headphones,
		title: "Connect Spotify",
		description: "Authenticate with OAuth 2.0 and link your Spotify account securely.",
		color: "text-[#1DB954]",
		bg: "bg-[#1DB954]/10 border-[#1DB954]/20",
	},
	{
		icon: ListMusic,
		title: "Choose podcasts to monitor",
		description: "Pick engineering, AI, startup, or any niche shows your team follows.",
		color: "text-violet-400",
		bg: "bg-violet-500/10 border-violet-500/20",
	},
	{
		icon: RefreshCw,
		title: "PodEvents checks Spotify",
		description: "Cron scheduling continuously scans for newly released episodes.",
		color: "text-sky-400",
		bg: "bg-sky-500/10 border-sky-500/20",
	},
	{
		icon: Bell,
		title: "Instant notifications",
		description: "Your team receives alerts on Slack, Discord, Telegram, or WhatsApp.",
		color: "text-emerald-400",
		bg: "bg-emerald-500/10 border-emerald-500/20",
	},
] as const;

export function HowItWorks() {
	return (
		<AnimatedSection className="py-24 md:py-32" id="how-it-works">
			<div className="mx-auto max-w-6xl px-6">
				<div className="mx-auto max-w-2xl text-center">
					<Badge className="mb-4" variant="emerald">
						How it works
					</Badge>
					<h2 className="font-semibold text-3xl text-zinc-50 tracking-tight md:text-4xl">
						From Spotify to your team in four steps
					</h2>
				</div>

				<StaggerGrid className="mx-auto mt-16 flex max-w-lg flex-col items-center gap-3">
					{STEPS.map((step, i) => (
						<StaggerItem className="w-full" key={step.title}>
							<Card className="group transition-all duration-300 hover:border-white/15 hover:bg-white/[0.05]">
								<CardContent className="flex items-start gap-4 p-5">
									<div
										className={`flex h-11 w-11 shrink-0 items-center justify-center rounded-xl border ${step.bg}`}
									>
										<step.icon className={`h-5 w-5 ${step.color}`} />
									</div>
									<div>
										<div className="flex items-center gap-2">
											<span className="font-mono text-emerald-500/60 text-xs">
												0{i + 1}
											</span>
											<h3 className="font-semibold text-zinc-100">{step.title}</h3>
										</div>
										<p className="mt-1 text-sm text-zinc-400">{step.description}</p>
									</div>
								</CardContent>
							</Card>
							{i < STEPS.length - 1 && (
								<div className="flex justify-center py-1">
									<ArrowDown className="h-4 w-4 text-zinc-600" />
								</div>
							)}
						</StaggerItem>
					))}
				</StaggerGrid>
			</div>
		</AnimatedSection>
	);
}