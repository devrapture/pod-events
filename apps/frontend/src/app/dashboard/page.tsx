"use client";

import { useAuth } from "@/components/auth-provider";

export default function DashboardPage() {
	const { user, isLoading, logout } = useAuth();

	if (isLoading) {
		return (
			<div className="flex min-h-screen items-center justify-center">
				<p className="text-zinc-400">Loading...</p>
			</div>
		);
	}

	return (
		<div className="mx-auto max-w-4xl px-6 pt-32">
			<div className="mb-8 flex items-center justify-between">
				<div>
					<h1 className="font-semibold text-2xl text-zinc-100">
						Welcome{user ? `, ${user.name}` : ""}
					</h1>
					<p className="mt-1 text-zinc-400">Dashboard</p>
				</div>
				<button
					className="inline-flex h-10 items-center justify-center rounded-lg border border-white/10 bg-white/5 px-5 font-medium text-sm text-zinc-100 backdrop-blur-sm hover:border-white/20 hover:bg-white/10"
					onClick={logout}
					type="button"
				>
					Sign Out
				</button>
			</div>
		</div>
	);
}
