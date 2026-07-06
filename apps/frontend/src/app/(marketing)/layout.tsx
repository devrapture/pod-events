import { LandingHeader } from "@/components/landing-header";
import { Footer } from "@/components/sections/footer";

export default function MarketingLayout({
	children,
}: Readonly<{ children: React.ReactNode }>) {
	return (
		<>
			<LandingHeader />
			<main>{children}</main>
			<Footer />
		</>
	);
}
