import { Code2, GitPullRequest, Globe, Package, Plug } from "lucide-react";

import { AnimatedSection } from "@/components/animated-section";
import { GitHubIcon } from "@/components/icons/github";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { GITHUB_URL } from "@/lib/constants";

const BENEFITS = [
	{ icon: Globe, label: "Self-host anywhere" },
	{ icon: GitPullRequest, label: "Contribute features" },
	{ icon: Plug, label: "Build integrations" },
	{ icon: Package, label: "Add notification providers" },
	{ icon: Code2, label: "Deploy on your stack" },
] as const;

export function OpenSource() {
	return (
		<AnimatedSection className="py-24 md:py-32" id="open-source">
			<div className="mx-auto max-w-6xl px-6">
				<Card className="relative overflow-hidden border-emerald-500/20 bg-gradient-to-br from-emerald-500/[0.06] via-transparent to-violet-500/[0.04] p-8 md:p-12">
					<div className="pointer-events-none absolute -top-24 -right-24 h-64 w-64 rounded-full bg-emerald-500/10 blur-3xl" />

					<div className="relative grid items-center gap-10 lg:grid-cols-2">
						<div>
							<Badge className="mb-4" variant="emerald">
								Open source
							</Badge>
							<h2 className="font-semibold text-3xl text-zinc-50 tracking-tight md:text-4xl">
								PodEvents is completely open source
							</h2>
							<p className="mt-4 text-zinc-400 leading-relaxed">
								Developers can self-host, contribute, build integrations, add new
								notification providers, and deploy anywhere. No vendor lock-in.
							</p>
							<Button asChild className="mt-8" size="lg">
								<a href={GITHUB_URL} rel="noopener noreferrer" target="_blank">
									<GitHubIcon className="h-5 w-5" />
									View on GitHub
								</a>
							</Button>
						</div>

						<div className="grid gap-3 sm:grid-cols-2">
							{BENEFITS.map((item) => (
								<div
									className="flex items-center gap-3 rounded-xl border border-white/[0.06] bg-white/[0.03] px-4 py-3 transition-colors hover:border-white/10 hover:bg-white/[0.05]"
									key={item.label}
								>
									<item.icon className="h-4 w-4 shrink-0 text-emerald-400" />
									<span className="text-sm text-zinc-300">{item.label}</span>
								</div>
							))}
						</div>
					</div>
				</Card>
			</div>
		</AnimatedSection>
	);
}