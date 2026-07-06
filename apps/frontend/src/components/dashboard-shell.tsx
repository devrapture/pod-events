"use client";

import {
	Download,
	LayoutDashboard,
	LogOut,
	Menu,
	Podcast,
	Radio,
	Search,
	X,
} from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import type { ReactNode } from "react";
import { useState } from "react";

import { useAuth } from "@/components/auth-provider";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import { cn } from "@/lib/utils";

const NAV_ITEMS = [
	{ label: "Overview", href: "/dashboard", icon: LayoutDashboard, exact: true },
	{
		label: "My Subscriptions",
		href: "/dashboard/subscriptions",
		icon: Podcast,
		exact: true,
	},
	{
		label: "Search Podcasts",
		href: "/dashboard/search",
		icon: Search,
		exact: true,
	},
	{
		label: "Import from Spotify",
		href: "/dashboard/import",
		icon: Download,
		exact: false,
	},
] as const;

function isNavActive(pathname: string, href: string, exact: boolean) {
	if (exact) return pathname === href;
	return pathname === href || pathname.startsWith(`${href}/`);
}

export function DashboardShell({ children }: { children: ReactNode }) {
	const pathname = usePathname();
	const { user, isLoading, logout } = useAuth();
	const [mobileOpen, setMobileOpen] = useState(false);
	const [confirmSignOut, setConfirmSignOut] = useState(false);

	return (
		<div className="min-h-screen bg-zinc-950 md:h-screen md:overflow-hidden">
			<aside className="fixed inset-y-0 left-0 z-40 hidden w-64 flex-col border-white/6 border-r bg-zinc-950 md:flex">
				<div className="flex h-16 items-center gap-2.5 border-white/6 border-b px-6">
					<div className="flex h-8 w-8 items-center justify-center rounded-lg border border-emerald-500/20 bg-emerald-500/10">
						<Radio className="h-4 w-4 text-emerald-400" />
					</div>
					<span className="font-semibold text-zinc-100 tracking-tight">
						PodEvents
					</span>
				</div>

				<nav className="flex flex-1 flex-col gap-1 p-4">
					{NAV_ITEMS.map((item) => {
						const isActive = isNavActive(pathname, item.href, item.exact);
						const Icon = item.icon;

						return (
							<Link
								className={cn(
									"flex items-center gap-3 rounded-lg px-3 py-2 font-medium text-sm transition-colors",
									isActive
										? "bg-emerald-500/10 text-emerald-400"
										: "text-zinc-400 hover:bg-white/5 hover:text-zinc-100",
								)}
								href={item.href}
								key={item.href}
							>
								<Icon className="h-4 w-4" />
								{item.label}
							</Link>
						);
					})}
				</nav>

				<div className="border-white/6 border-t p-4">
					{isLoading ? (
						<div className="h-10 animate-pulse rounded-lg bg-white/5" />
					) : user ? (
						<div className="flex items-center gap-3">
							{user.avatar_url ? (
								// biome-ignore lint/performance/noImgElement: dynamic avatar URL from Spotify
								<img
									alt={user.name}
									className="h-9 w-9 rounded-full"
									src={user.avatar_url}
								/>
							) : (
								<div className="flex h-9 w-9 items-center justify-center rounded-full bg-emerald-500/20 font-medium text-emerald-400 text-sm">
									{user.name.charAt(0).toUpperCase()}
								</div>
							)}
							<div className="min-w-0 flex-1">
								<p className="truncate font-medium text-sm text-zinc-100">
									{user.name}
								</p>
								<p className="truncate text-xs text-zinc-500">{user.email}</p>
							</div>
							<button
								aria-label="Sign out"
								className="flex items-center gap-2 rounded-lg px-3 py-2 text-sm text-zinc-400 transition-colors hover:bg-white/5 hover:text-zinc-100"
								onClick={() => setConfirmSignOut(true)}
								type="button"
							>
								<LogOut className="h-4 w-4" />
								Sign Out
							</button>
						</div>
					) : null}
				</div>
			</aside>

			<div className="flex min-h-screen min-w-0 flex-col md:ml-64 md:h-screen md:overflow-hidden">
				<header className="flex h-16 items-center justify-between border-white/6 border-b px-4 md:hidden">
					<div className="flex items-center gap-2.5">
						<div className="flex h-8 w-8 items-center justify-center rounded-lg border border-emerald-500/20 bg-emerald-500/10">
							<Radio className="h-4 w-4 text-emerald-400" />
						</div>
						<span className="font-semibold text-zinc-100 tracking-tight">
							PodEvents
						</span>
					</div>
					<button
						aria-label="Toggle menu"
						className="rounded-lg p-2 text-zinc-400 hover:bg-white/5 hover:text-zinc-100"
						onClick={() => setMobileOpen((open) => !open)}
						type="button"
					>
						{mobileOpen ? (
							<X className="h-5 w-5" />
						) : (
							<Menu className="h-5 w-5" />
						)}
					</button>
				</header>

				{mobileOpen && (
					<div className="border-white/6 border-b bg-zinc-950 px-4 py-4 md:hidden">
						<nav className="flex flex-col gap-1">
							{NAV_ITEMS.map((item) => {
								const isActive = isNavActive(pathname, item.href, item.exact);
								const Icon = item.icon;

								return (
									<Link
										className={cn(
											"flex items-center gap-3 rounded-lg px-3 py-2 font-medium text-sm transition-colors",
											isActive
												? "bg-emerald-500/10 text-emerald-400"
												: "text-zinc-400 hover:bg-white/5 hover:text-zinc-100",
										)}
										href={item.href}
										key={item.href}
										onClick={() => setMobileOpen(false)}
									>
										<Icon className="h-4 w-4" />
										{item.label}
									</Link>
								);
							})}
						</nav>
						{user && (
							<div className="mt-4 flex items-center justify-between border-white/6 border-t pt-4">
								<div className="flex items-center gap-3">
									{user.avatar_url ? (
										// biome-ignore lint/performance/noImgElement: dynamic avatar URL from Spotify
										<img
											alt={user.name}
											className="h-8 w-8 rounded-full"
											src={user.avatar_url}
										/>
									) : (
										<div className="flex h-8 w-8 items-center justify-center rounded-full bg-emerald-500/20 font-medium text-emerald-400 text-xs">
											{user.name.charAt(0).toUpperCase()}
										</div>
									)}
									<span className="text-sm text-zinc-300">{user.name}</span>
								</div>
								<button
									className="text-sm text-zinc-400 hover:text-zinc-100"
									onClick={() => {
										setConfirmSignOut(true);
										setMobileOpen(false);
									}}
									type="button"
								>
									Sign Out
								</button>
							</div>
						)}
					</div>
				)}

				<main className="flex min-h-0 flex-1 flex-col overflow-hidden p-6 md:p-8">
					{children}
				</main>
			</div>

			<ConfirmDialog
				confirmLabel="Sign Out"
				description="Are you sure you want to sign out?"
				onConfirm={logout}
				onOpenChange={setConfirmSignOut}
				open={confirmSignOut}
				title="Sign out"
				variant="destructive"
			/>
		</div>
	);
}
