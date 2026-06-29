import { MessageCircle, MessageSquare, Send, Slack } from "lucide-react";

import {
	AnimatedSection,
	StaggerGrid,
	StaggerItem,
} from "@/components/animated-section";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const PLATFORMS = [
	{
		icon: Slack,
		name: "Slack",
		description: "Rich Block Kit notifications.",
		color: "#E01E5A",
		preview: (
			<div className="mt-4 space-y-2 rounded-lg border border-white/5 bg-zinc-950/80 p-3 font-mono text-xs">
				<p className="text-[#E01E5A]">#dev-alerts</p>
				<p className="text-zinc-300">🎙 Changelog — new episode</p>
				<p className="text-zinc-500">Episode 612: Platform engineering</p>
			</div>
		),
	},
	{
		icon: MessageSquare,
		name: "Discord",
		description: "Beautiful embedded alerts.",
		color: "#5865F2",
		preview: (
			<div className="mt-4 rounded-lg border border-[#5865F2]/30 bg-[#5865F2]/5 p-3 text-xs">
				<p className="font-semibold text-[#5865F2]">PodEvents</p>
				<p className="mt-1 text-zinc-300">New episode: Latent Space</p>
				<p className="text-zinc-500">Agents in production — listen now</p>
			</div>
		),
	},
	{
		icon: Send,
		name: "Telegram",
		description: "Instant mobile notifications.",
		color: "#26A5E4",
		preview: (
			<div className="mt-4 rounded-lg border border-[#26A5E4]/30 bg-[#26A5E4]/5 p-3 text-xs">
				<p className="text-zinc-300">PodEvents Bot</p>
				<p className="mt-1 text-zinc-400">Syntax FM just released a new episode.</p>
			</div>
		),
	},
	{
		icon: MessageCircle,
		name: "WhatsApp",
		description: "Business-ready message delivery.",
		color: "#25D366",
		preview: (
			<div className="mt-4 rounded-lg border border-[#25D366]/30 bg-[#25D366]/5 p-3 text-xs">
				<p className="text-zinc-300">PodEvents</p>
				<p className="mt-1 text-zinc-400">🎧 New: Invest Like the Best</p>
			</div>
		),
	},
] as const;

export function NotificationPlatforms() {
	return (
		<AnimatedSection className="py-24 md:py-32" id="platforms">
			<div className="mx-auto max-w-6xl px-6">
				<div className="mx-auto max-w-2xl text-center">
					<Badge className="mb-4" variant="emerald">
						Notification platforms
					</Badge>
					<h2 className="font-semibold text-3xl text-zinc-50 tracking-tight md:text-4xl">
						Meet your team where they already are
					</h2>
				</div>

				<StaggerGrid className="mt-14 grid gap-4 sm:grid-cols-2">
					{PLATFORMS.map((platform) => (
						<StaggerItem key={platform.name}>
							<Card className="h-full transition-all duration-300 hover:border-white/15 hover:bg-white/[0.05]">
								<CardHeader>
									<div className="flex items-center gap-3">
										<div
											className="flex h-10 w-10 items-center justify-center rounded-lg border"
											style={{
												borderColor: `${platform.color}33`,
												backgroundColor: `${platform.color}14`,
											}}
										>
											<platform.icon
												className="h-5 w-5"
												style={{ color: platform.color }}
											/>
										</div>
										<div>
											<CardTitle>{platform.name}</CardTitle>
											<p className="text-sm text-zinc-400">{platform.description}</p>
										</div>
									</div>
								</CardHeader>
								<CardContent>{platform.preview}</CardContent>
							</Card>
						</StaggerItem>
					))}
				</StaggerGrid>
			</div>
		</AnimatedSection>
	);
}