"use client";

import { useAuth } from "@/components/auth-provider";

export default function DashboardPage() {
	const { user, isLoading } = useAuth();

	if (isLoading) {
		return (
			<div className="flex min-h-[50vh] items-center justify-center">
				<p className="text-zinc-400">Loading...</p>
			</div>
		);
	}

	return (
		<div className="mx-auto max-w-4xl">
			<div className="mb-8">
				<h1 className="font-semibold text-2xl text-zinc-100">
					Welcome{user ? `, ${user.name}` : ""}
				</h1>
				<p className="mt-1 text-zinc-400">
					Manage your podcast notifications from here.
				</p>
			</div>
		</div>
	);
}
