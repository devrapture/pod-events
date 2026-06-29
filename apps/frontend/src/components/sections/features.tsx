import {
	Database,
	Lock,
	RefreshCw,
	Server,
	Shield,
	Users,
	Zap,
} from "lucide-react";

import {
	AnimatedSection,
	StaggerGrid,
	StaggerItem,
} from "@/components/animated-section";
import { Badge } from "@/components/ui/badge";
import {
	Card,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";

const FEATURES = [
	{
		icon: RefreshCw,
		title: "Automatic Podcast Monitoring",
		description:
			"Continuously checks Spotify for newly released episodes on a reliable cron schedule.",
	},
	{
		icon: Zap,
		title: "Multi-channel Notifications",
		description:
			"Deliver alerts directly to Slack, Discord, Telegram, or WhatsApp.",
	},
	{
		icon: Users,
		title: "Team Ready",
		description:
			"Keep engineering, marketing, product, and leadership teams informed without manual effort.",
	},
	{
		icon: Shield,
		title: "Duplicate Prevention",
		description:
			"Built-in idempotency guarantees notifications are never sent twice.",
	},
	{
		icon: Lock,
		title: "Secure by Default",
		description: "Spotify tokens are encrypted using AES-256-GCM at rest.",
	},
	{
		icon: Server,
		title: "Production Architecture",
		description:
			"PostgreSQL, cron scheduling, connection pooling, OAuth, middleware, and structured services.",
	},
] as const;

export function Features() {
	return (
		<AnimatedSection className="py-24 md:py-32" id="features">
			<div className="mx-auto max-w-6xl px-6">
				<div className="mx-auto max-w-2xl text-center">
					<Badge className="mb-4" variant="emerald">
						Features
					</Badge>
					<h2 className="font-semibold text-3xl text-zinc-50 tracking-tight md:text-4xl">
						Production-ready from day one
					</h2>
					<p className="mt-4 text-zinc-400">
						Everything you need to monitor podcasts reliably and notify your team
						instantly.
					</p>
				</div>

				<StaggerGrid className="mt-14 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
					{FEATURES.map((feature) => (
						<StaggerItem key={feature.title}>
							<Card className="h-full transition-all duration-300 hover:border-white/15 hover:bg-white/[0.05] hover:shadow-lg hover:shadow-emerald-500/5">
								<CardHeader>
									<div className="mb-3 flex h-10 w-10 items-center justify-center rounded-lg border border-emerald-500/20 bg-emerald-500/10">
										<feature.icon className="h-5 w-5 text-emerald-400" />
									</div>
									<CardTitle>{feature.title}</CardTitle>
									<CardDescription>{feature.description}</CardDescription>
								</CardHeader>
							</Card>
						</StaggerItem>
					))}
				</StaggerGrid>

				<div className="mt-8 flex justify-center">
					<div className="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/[0.03] px-4 py-2 text-sm text-zinc-400">
						<Database className="h-4 w-4 text-emerald-400" />
						Backed by PostgreSQL with connection pooling
					</div>
				</div>
			</div>
		</AnimatedSection>
	);
}