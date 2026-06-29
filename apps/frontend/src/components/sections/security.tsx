import { Check } from "lucide-react";

import { AnimatedSection } from "@/components/animated-section";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";

const SECURITY_FEATURES = [
	"OAuth 2.0 authentication",
	"AES-256-GCM encrypted credentials",
	"Duplicate notification protection",
	"Secure cron endpoints",
	"Automatic token refresh",
	"Connection pooling",
	"Panic recovery",
	"Structured logging",
] as const;

export function Security() {
	return (
		<AnimatedSection className="py-24 md:py-32" id="security">
			<div className="mx-auto max-w-6xl px-6">
				<div className="grid items-center gap-12 lg:grid-cols-2">
					<div>
						<Badge className="mb-4" variant="emerald">
							Security
						</Badge>
						<h2 className="font-semibold text-3xl text-zinc-50 tracking-tight md:text-4xl">
							Enterprise-grade reliability
						</h2>
						<p className="mt-4 text-zinc-400 leading-relaxed">
							PodEvents is built with the security patterns you expect from
							production infrastructure—encryption, idempotency, and secure
							authentication out of the box.
						</p>
					</div>

					<Card className="p-6">
						<ul className="space-y-4">
							{SECURITY_FEATURES.map((feature) => (
								<li className="flex items-start gap-3" key={feature}>
									<div className="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-emerald-500/15">
										<Check className="h-3 w-3 text-emerald-400" />
									</div>
									<span className="text-sm text-zinc-300">{feature}</span>
								</li>
							))}
						</ul>
					</Card>
				</div>
			</div>
		</AnimatedSection>
	);
}