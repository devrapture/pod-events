import { Radio } from "lucide-react";
import Link from "next/link";

import {
	CONTRIBUTING_URL,
	DOCS_URL,
	GITHUB_URL,
	ISSUES_URL,
	LICENSE_URL,
} from "@/lib/constants";

const FOOTER_LINKS = [
	{ label: "GitHub", href: GITHUB_URL },
	{ label: "Documentation", href: DOCS_URL },
	{ label: "Issues", href: ISSUES_URL },
	{ label: "Contributing", href: CONTRIBUTING_URL },
	{ label: "License", href: LICENSE_URL },
] as const;

export function Footer() {
	return (
		<footer className="border-white/[0.06] border-t py-12">
			<div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-8 px-6 md:flex-row">
				<Link className="flex items-center gap-2.5" href="/">
					<div className="flex h-7 w-7 items-center justify-center rounded-lg border border-emerald-500/20 bg-emerald-500/10">
						<Radio className="h-3.5 w-3.5 text-emerald-400" />
					</div>
					<span className="font-medium text-sm text-zinc-300">PodEvents</span>
				</Link>

				<nav className="flex flex-wrap justify-center gap-6">
					{FOOTER_LINKS.map((link) => (
						<a
							className="text-sm text-zinc-500 transition-colors hover:text-zinc-300"
							href={link.href}
							key={link.href}
							rel="noopener noreferrer"
							target="_blank"
						>
							{link.label}
						</a>
					))}
				</nav>

				<p className="text-sm text-zinc-600">
					© {new Date().getFullYear()} PodEvents. Open source.
				</p>
			</div>
		</footer>
	);
}