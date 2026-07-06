"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useEffect, useRef, useState } from "react";

import { useAuth } from "@/components/auth-provider";
import { getNgrokHeaders } from "@/lib/api-headers";
import type { User } from "@/lib/auth";

type AuthExchangeData = {
	token: string;
	user: User;
};

function CallbackContent() {
	const router = useRouter();
	const searchParams = useSearchParams();
	const { setToken } = useAuth();
	const [error, setError] = useState<string | null>(null);
	const handledCodeRef = useRef<string | null>(null);

	useEffect(() => {
		const code = searchParams.get("code");
		const errorParam = searchParams.get("error");

		if (errorParam) {
			setError(decodeURIComponent(errorParam));
			return;
		}

		if (!code) {
			setError("No authentication code received");
			return;
		}

		if (handledCodeRef.current === code) {
			return;
		}
		handledCodeRef.current = code;

		fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/exchange`, {
			method: "POST",
			headers: { ...getNgrokHeaders(), "Content-Type": "application/json" },
			body: JSON.stringify({ code }),
		})
			.then(async (res) => {
				const body = await res.json();
				if (!res.ok)
					throw new Error(body.error?.message ?? "Authentication failed");
				return body.data as AuthExchangeData;
			})
			.then(({ token, user }) => setToken(token, user))
			.then((signedIn) => {
				if (!signedIn) {
					throw new Error("Unable to verify your account");
				}
				router.replace("/dashboard");
			})
			.catch((err) => setError(err.message));
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
