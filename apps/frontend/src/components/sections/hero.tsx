"use client";

import { motion } from "framer-motion";
import { ArrowRight, Bell, BookOpen, Sparkles } from "lucide-react";

import { GitHubIcon } from "@/components/icons/github";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { DOCS_URL, GITHUB_URL, TECH_BADGES } from "@/lib/constants";

function NotificationPreview() {
	return (
		<div className="relative">
			<div className="absolute -inset-4 rounded-3xl bg-emerald-500/10 blur-3xl" />
			<div className="relative space-y-3">
				<motion.div
					animate={{ y: [0, -6, 0] }}
					className="ml-auto max-w-xs rounded-xl border border-white/10 bg-zinc-900/90 p-4 shadow-2xl backdrop-blur-xl"
					transition={{ duration: 4, repeat: Number.POSITIVE_INFINITY, ease: "easeInOut" }}
				>
					<div className="mb-2 flex items-center gap-2">
						<div className="h-2 w-2 rounded-full bg-[#E01E5A]" />
						<span className="font-medium text-xs text-zinc-400">#engineering</span>
					</div>
					<p className="font-semibold text-sm text-zinc-100">🎙 New episode: Go Time</p>
					<p className="mt-1 text-xs text-zinc-400">
						&quot;Building resilient APIs at scale&quot; just dropped on Spotify.
					</p>
				</motion.div>

				<motion.div
					animate={{ y: [0, 6, 0] }}
					className="max-w-xs rounded-xl border border-white/10 bg-zinc-900/90 p-4 shadow-2xl backdrop-blur-xl"
					transition={{ duration: 5, repeat: Number.POSITIVE_INFINITY, ease: "easeInOut", delay: 0.5 }}
				>
					<div className="mb-2 flex items-center gap-2">
						<div className="h-2 w-2 rounded-full bg-[#5865F2]" />
						<span className="font-medium text-xs text-zinc-400">ai-research</span>
					</div>
					<p className="font-semibold text-sm text-zinc-100">Latent Space — new episode</p>
					<p className="mt-1 text-xs text-zinc-400">
						Agents, evals, and production ML. Listen now →
					</p>
				</motion.div>

				<motion.div
					animate={{ y: [0, -4, 0] }}
					className="ml-8 max-w-[220px] rounded-xl border border-white/10 bg-zinc-900/90 p-3 shadow-2xl backdrop-blur-xl"
					transition={{ duration: 3.5, repeat: Number.POSITIVE_INFINITY, ease: "easeInOut", delay: 1 }}
				>
					<div className="flex items-center gap-2">
						<Bell className="h-3.5 w-3.5 text-emerald-400" />
						<p className="text-xs text-zinc-300">Telegram · Syntax FM</p>
					</div>
				</motion.div>
			</div>
		</div>
	);
}

export function Hero() {
	return (
		<section className="relative overflow-hidden pt-32 pb-20 md:pt-40 md:pb-28">
			<div className="hero-glow pointer-events-none absolute inset-0" />
			<div className="grid-overlay pointer-events-none absolute inset-0 opacity-40" />

			<div className="relative mx-auto grid max-w-6xl items-center gap-16 px-6 lg:grid-cols-2">
				<div>
					<motion.div
						animate={{ opacity: 1, y: 0 }}
						initial={{ opacity: 0, y: 20 }}
						transition={{ duration: 0.6 }}
					>
						<Badge className="mb-6 gap-1.5" variant="emerald">
							<Sparkles className="h-3 w-3" />
							Open-source podcast monitoring
						</Badge>

						<h1 className="font-semibold text-4xl text-zinc-50 tracking-tight sm:text-5xl lg:text-[3.25rem] lg:leading-[1.1]">
							Spotify Podcast Notifications for Your Entire Team.
						</h1>

						<p className="mt-6 max-w-xl text-lg text-zinc-400 leading-relaxed">
							PodEvents monitors Spotify podcasts 24/7 and instantly delivers new
							episode notifications to Slack, Discord, Telegram, and WhatsApp—so
							your team never misses important industry conversations.
						</p>

						<div className="mt-8 flex flex-wrap gap-3">
							<Button asChild size="lg">
								<a href={GITHUB_URL} rel="noopener noreferrer" target="_blank">
									<GitHubIcon className="h-4 w-4" />
									View on GitHub
									<ArrowRight className="h-4 w-4" />
								</a>
							</Button>
							<Button asChild size="lg" variant="secondary">
								<a href={DOCS_URL} rel="noopener noreferrer" target="_blank">
									<BookOpen className="h-4 w-4" />
									Documentation
								</a>
							</Button>
						</div>

						<div className="mt-8 flex flex-wrap gap-2">
							{TECH_BADGES.map((badge) => (
								<Badge key={badge} variant="outline">
									{badge}
								</Badge>
							))}
						</div>
					</motion.div>
				</div>

				<motion.div
					animate={{ opacity: 1, x: 0 }}
					className="hidden lg:block"
					initial={{ opacity: 0, x: 30 }}
					transition={{ duration: 0.7, delay: 0.2 }}
				>
					<NotificationPreview />
				</motion.div>
			</div>
		</section>
	);
}