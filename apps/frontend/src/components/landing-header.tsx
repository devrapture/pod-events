"use client";

import { Menu, Radio, X } from "lucide-react";
import Link from "next/link";
import { useEffect, useState } from "react";

import { GitHubIcon } from "@/components/icons/github";
import { Button } from "@/components/ui/button";
import { DOCS_URL, GITHUB_URL, NAV_LINKS } from "@/lib/constants";
import { cn } from "@/lib/utils";

export function LandingHeader() {
	const [scrolled, setScrolled] = useState(false);
	const [mobileOpen, setMobileOpen] = useState(false);

	useEffect(() => {
		const onScroll = () => setScrolled(window.scrollY > 12);
		onScroll();
		window.addEventListener("scroll", onScroll, { passive: true });
		return () => window.removeEventListener("scroll", onScroll);
	}, []);

	return (
		<header
			className={cn(
				"fixed inset-x-0 top-0 z-50 transition-all duration-300",
				scrolled
					? "border-white/[0.06] border-b bg-zinc-950/80 backdrop-blur-xl"
					: "bg-transparent",
			)}
		>
			<div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-6">
				<Link className="flex items-center gap-2.5" href="/">
					<div className="flex h-8 w-8 items-center justify-center rounded-lg border border-emerald-500/20 bg-emerald-500/10">
						<Radio className="h-4 w-4 text-emerald-400" />
					</div>
					<span className="font-semibold text-zinc-100 tracking-tight">
						PodEvents
					</span>
				</Link>

				<nav className="hidden items-center gap-8 md:flex">
					{NAV_LINKS.map((link) => (
						<a
							className="text-sm text-zinc-400 transition-colors hover:text-zinc-100"
							href={link.href}
							key={link.href}
						>
							{link.label}
						</a>
					))}
				</nav>

				<div className="hidden items-center gap-3 md:flex">
					<Button asChild size="sm" variant="ghost">
						<a href={DOCS_URL} rel="noopener noreferrer" target="_blank">
							Docs
						</a>
					</Button>
					<Button asChild size="sm">
						<a href={GITHUB_URL} rel="noopener noreferrer" target="_blank">
							<GitHubIcon className="h-4 w-4" />
							GitHub
						</a>
					</Button>
				</div>

				<button
					aria-label="Toggle menu"
					className="rounded-lg p-2 text-zinc-400 hover:bg-white/5 hover:text-zinc-100 md:hidden"
					onClick={() => setMobileOpen((v) => !v)}
					type="button"
				>
					{mobileOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
				</button>
			</div>

			{mobileOpen && (
				<div className="border-white/[0.06] border-t bg-zinc-950/95 px-6 py-4 backdrop-blur-xl md:hidden">
					<nav className="flex flex-col gap-4">
						{NAV_LINKS.map((link) => (
							<a
								className="text-sm text-zinc-300"
								href={link.href}
								key={link.href}
								onClick={() => setMobileOpen(false)}
							>
								{link.label}
							</a>
						))}
						<div className="flex gap-3 pt-2">
							<Button asChild className="flex-1" size="sm" variant="secondary">
								<a href={DOCS_URL} rel="noopener noreferrer" target="_blank">
									Docs
								</a>
							</Button>
							<Button asChild className="flex-1" size="sm">
								<a href={GITHUB_URL} rel="noopener noreferrer" target="_blank">
									GitHub
								</a>
							</Button>
						</div>
					</nav>
				</div>
			)}
		</header>
	);
}