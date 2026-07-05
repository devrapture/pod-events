"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useEffect, useState } from "react";

import { useAuth } from "@/components/auth-provider";

function CallbackContent() {
	const router = useRouter();
	const searchParams = useSearchParams();
	const { setToken } = useAuth();
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		const token = searchParams.get("token");
		const errorParam = searchParams.get("error");

		if (errorParam) {
			setError(decodeURIComponent(errorParam));
			return;
		}

		if (!token) {
			setError("No token received from authentication");
			return;
		}

		setToken(token).then(() => {
			router.replace("/dashboard");
		});
	}, [searchParams, setToken, router]);

	if (error) {
		return (
			<div className="flex min-h-screen items-center justify-center">
				<div className="text-center">
					<h1 className="mb-4 font-semibold text-2xl text-zinc-100">
						Authentication Failed
					</h1>
					<p className="mb-6 text-zinc-400">{error}</p>
					<a
						className="inline-flex h-10 items-center justify-center rounded-lg bg-emerald-500 px-5 font-medium text-sm text-zinc-950 shadow-emerald-500/20 shadow-lg hover:bg-emerald-400"
						href="/"
					>
						Go Home
					</a>
				</div>
			</div>
		);
	}

	return (
		<div className="flex min-h-screen items-center justify-center">
			<p className="text-zinc-400">Signing you in...</p>
		</div>
	);
}

export default function AuthCallbackPage() {
	return (
		<Suspense
			fallback={
				<div className="flex min-h-screen items-center justify-center">
					<p className="text-zinc-400">Signing you in...</p>
				</div>
			}
		>
			<CallbackContent />
		</Suspense>
	);
}
