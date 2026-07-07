"use client";

import { useAuth } from "@/components/auth-provider";
import { Badge } from "@/components/ui/badge";
import { getTimeGreeting } from "@/lib/utils";

export function OverviewHeader() {
	const { user } = useAuth();
	const greeting = user ? getTimeGreeting(user.name) : "Welcome";

	return (
		<div className="mb-8 flex items-start justify-between gap-4">
			<h1 className="font-semibold text-2xl text-zinc-100 tracking-tight">
				{greeting}
			</h1>
		</div>
	);
}