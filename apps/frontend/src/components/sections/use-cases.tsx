import { Bot, Briefcase, Code2, Megaphone, Rocket, Users } from "lucide-react";

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

const USE_CASES = [
	{
		icon: Code2,
		title: "Engineering Teams",
		description: "Track Go Time, Software Engineering Daily, Syntax, Changelog, and more.",
		accent: "from-emerald-500/20 to-transparent",
	},
	{
		icon: Megaphone,
		title: "Marketing Teams",
		description:
			"Know immediately when marketing leaders publish new episodes.",
		accent: "from-violet-500/20 to-transparent",
	},
	{
		icon: Bot,
		title: "AI Teams",
		description: "Stay updated with AI podcasts without checking Spotify.",
		accent: "from-sky-500/20 to-transparent",
	},
	{
		icon: Rocket,
		title: "Startup Founders",
		description: "Never miss discussions from founders and investors.",
		accent: "from-amber-500/20 to-transparent",
	},
	{
		icon: Users,
		title: "Developer Communities",
		description:
			"Automatically share podcast releases inside Slack or Discord communities.",
		accent: "from-rose-500/20 to-transparent",
	},
	{
		icon: Briefcase,
		title: "Product Teams",
		description:
			"Follow product leaders and stay ahead of industry conversations.",
		accent: "from-indigo-500/20 to-transparent",
	},
] as const;

export function UseCases() {
	return (
		<AnimatedSection className="py-24 md:py-32" id="use-cases">
			<div className="mx-auto max-w-6xl px-6">
				<div className="mx-auto max-w-2xl text-center">
					<Badge className="mb-4" variant="emerald">
						Who uses PodEvents
					</Badge>
					<h2 className="font-semibold text-3xl text-zinc-50 tracking-tight md:text-4xl">
						Built for teams that learn from podcasts
					</h2>
				</div>

				<StaggerGrid className="mt-14 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
					{USE_CASES.map((useCase) => (
						<StaggerItem key={useCase.title}>
							<Card className="group relative h-full overflow-hidden transition-all duration-300 hover:border-white/15 hover:bg-white/[0.05]">
								<div
									className={`pointer-events-none absolute inset-0 bg-gradient-to-br ${useCase.accent} opacity-0 transition-opacity duration-300 group-hover:opacity-100`}
								/>
								<CardHeader className="relative">
									<div className="mb-3 flex h-10 w-10 items-center justify-center rounded-lg border border-white/10 bg-white/5">
										<useCase.icon className="h-5 w-5 text-zinc-200" />
									</div>
									<CardTitle>{useCase.title}</CardTitle>
									<CardDescription>{useCase.description}</CardDescription>
								</CardHeader>
							</Card>
						</StaggerItem>
					))}
				</StaggerGrid>
			</div>
		</AnimatedSection>
	);
}