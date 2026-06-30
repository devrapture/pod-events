import { AnimatedSection } from "@/components/animated-section";
import { Badge } from "@/components/ui/badge";

const TEAMS = [
	{ role: "Engineering teams", podcasts: "software podcasts" },
	{ role: "Marketing teams", podcasts: "growth podcasts" },
	{ role: "Product teams", podcasts: "product leaders" },
	{ role: "AI teams", podcasts: "machine learning podcasts" },
] as const;

export function Problem() {
	return (
		<AnimatedSection className="py-24 md:py-32" id="problem">
			<div className="mx-auto max-w-6xl px-6">
				<div className="mx-auto max-w-3xl text-center">
					<Badge className="mb-4" variant="emerald">
						The problem
					</Badge>
					<h2 className="font-semibold text-3xl text-zinc-50 tracking-tight md:text-4xl">
						Teams follow dozens of podcasts. Nobody checks Spotify every day.
					</h2>
				</div>

				<div className="mt-14 grid gap-4 sm:grid-cols-2">
					{TEAMS.map((team) => (
						<div
							className="group rounded-2xl border border-white/[0.06] bg-white/[0.02] p-6 transition-all duration-300 hover:border-white/10 hover:bg-white/[0.04]"
							key={team.role}
						>
							<p className="font-medium text-zinc-200">{team.role}</p>
							<p className="mt-1 text-sm text-zinc-500">
								follow {team.podcasts}
							</p>
						</div>
					))}
				</div>

				<div className="mx-auto mt-12 max-w-2xl rounded-2xl border border-emerald-500/20 bg-emerald-500/[0.04] p-8 text-center">
					<p className="text-lg text-zinc-300 leading-relaxed">
						PodEvents solves this by automatically monitoring your favorite
						podcasts and notifying your entire team the moment a new episode is
						released.
					</p>
				</div>
			</div>
		</AnimatedSection>
	);
}