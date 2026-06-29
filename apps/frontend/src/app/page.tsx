import { Architecture } from "@/components/sections/architecture";
import { Features } from "@/components/sections/features";
import { FinalCTA } from "@/components/sections/final-cta";
import { Hero } from "@/components/sections/hero";
import { HowItWorks } from "@/components/sections/how-it-works";
import { NotificationPlatforms } from "@/components/sections/notification-platforms";
import { OpenSource } from "@/components/sections/open-source";
import { Problem } from "@/components/sections/problem";
import { Security } from "@/components/sections/security";
import { UseCases } from "@/components/sections/use-cases";

export default function Home() {
	return (
		<>
			<Hero />
			<Problem />
			<HowItWorks />
			<Features />
			<UseCases />
			<NotificationPlatforms />
			<Architecture />
			<Security />
			<OpenSource />
			<FinalCTA />
		</>
	);
}