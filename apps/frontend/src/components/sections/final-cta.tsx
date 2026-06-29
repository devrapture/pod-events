import { ArrowRight, BookOpen } from "lucide-react";

import { AnimatedSection } from "@/components/animated-section";
import { GitHubIcon } from "@/components/icons/github";
import { Button } from "@/components/ui/button";
import { DOCS_URL, GITHUB_URL } from "@/lib/constants";

export function FinalCTA() {
	return (
		<AnimatedSection className="py-24 md:py-32">
			<div className="mx-auto max-w-6xl px-6">
				<div className="relative overflow-hidden rounded-3xl border border-white/10 bg-gradient-to-b from-white/[0.05] to-transparent px-8 py-16 text-center md:px-16">
					<div className="pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_at_center,_var(--tw-gradient-stops))] from-emerald-500/10 via-transparent to-transparent" />

					<div className="relative">
						<h2 className="font-semibold text-3xl text-zinc-50 tracking-tight md:text-4xl">
							Your Team Should Never Miss Another Podcast Episode.
						</h2>
						<p className="mx-auto mt-4 max-w-xl text-lg text-zinc-400">
							Monitor Spotify podcasts once. Keep your entire team informed
							forever.
						</p>

						<div className="mt-8 flex flex-wrap justify-center gap-3">
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
									Read Documentation
								</a>
							</Button>
						</div>
					</div>
				</div>
			</div>
		</AnimatedSection>
	);
}