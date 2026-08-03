import { Radio } from "lucide-react";
import Link from "next/link";

import {
	CONTRIBUTING_URL,
	DOCS_URL,
	GITHUB_URL,
	ISSUES_URL,
	LICENSE_URL,
} from "@/lib/constants";
import LogoWithText from "../icons/logo-with-text";

const FOOTER_LINKS = [
	{ label: "GitHub", href: GITHUB_URL, external: true },
	{ label: "Documentation", href: DOCS_URL, external: true },
	{ label: "Issues", href: ISSUES_URL, external: true },
	{ label: "Contributing", href: CONTRIBUTING_URL, external: true },
	{ label: "License", href: LICENSE_URL, external: true },
	{ label: "Privacy", href: "/privacy", external: false },
] as const;

export function Footer() {
	return (
		<footer className="border-white/6 border-t py-12">
			<div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-8 px-6 md:flex-row">
				<Link className="flex items-center gap-2.5" href="/">
					<LogoWithText className="h-8 w-auto" />
				</Link>

				<nav className="flex flex-wrap justify-center gap-6">
					{FOOTER_LINKS.map((link) =>
						link.external ? (
							<a
								className="text-sm text-zinc-500 transition-colors hover:text-zinc-300"
								href={link.href}
								key={link.href}
								rel="noopener noreferrer"
								target="_blank"
							>
								{link.label}
							</a>
						) : (
							<Link
								className="text-sm text-zinc-500 transition-colors hover:text-zinc-300"
								href={link.href}
								key={link.href}
							>
								{link.label}
							</Link>
						),
					)}
				</nav>

				<p className="text-sm text-zinc-600">
					© {new Date().getFullYear()} PodEvents. Open source.
				</p>
			</div>
		</footer>
	);
}
