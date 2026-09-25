"use client";

import { CheckCircle2, Loader2, LockKeyhole, Music2 } from "lucide-react";
import { useSearchParams } from "next/navigation";
import { Suspense, useState } from "react";

import { useAuth } from "@/components/auth-provider";
import { Button } from "@/components/ui/button";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import { getAPIErrorMessage } from "@/lib/api-error";
import { apis } from "@/services/api-services";

function SpotifyOwnerContent() {
	const searchParams = useSearchParams();
	const { user, isLoading } = useAuth();
	const [isConnecting, setIsConnecting] = useState(false);
	const [requestError, setRequestError] = useState<string | null>(null);
	const connected = searchParams.get("connected") === "1";
	const callbackError = searchParams.get("error");

	const connectSpotify = async () => {
		setIsConnecting(true);
		setRequestError(null);
		try {
			const response = await apis.auth.spotifyOwnerLogin();
			const authorizationURL = response.data.data?.authorization_url;
			if (!authorizationURL) {
				throw new Error("Spotify authorization URL was not returned");
			}
			window.location.assign(authorizationURL);
		} catch (error) {
			setRequestError(getAPIErrorMessage(error, "Unable to connect Spotify"));
			setIsConnecting(false);
		}
	};

	return (
		<main className="flex min-h-screen items-center justify-center bg-zinc-950 px-6">
			<Card className="w-full max-w-lg">
				<CardHeader>
					<div className="mb-3 flex h-11 w-11 items-center justify-center rounded-xl border border-emerald-500/20 bg-emerald-500/10">
						<LockKeyhole className="h-5 w-5 text-emerald-400" />
					</div>
					<CardTitle>Spotify service account</CardTitle>
					<CardDescription>
						Owner-only setup for the shared Spotify credential used by
						PodEvents. This page is not linked from the application.
					</CardDescription>
				</CardHeader>
				<CardContent className="space-y-5">
					{connected && (
						<div className="flex items-center gap-2 rounded-lg border border-emerald-500/20 bg-emerald-500/10 p-3 text-emerald-300 text-sm">
							<CheckCircle2 className="h-4 w-4" />
							Spotify credential connected successfully.
						</div>
					)}
					{(callbackError || requestError) && (
						<div className="rounded-lg border border-red-500/20 bg-red-500/10 p-3 text-red-300 text-sm">
							{callbackError ?? requestError}
						</div>
					)}
					<p className="text-sm text-zinc-400">
						Signed in as {isLoading ? "…" : (user?.email ?? "unknown user")}
					</p>
					<Button
						disabled={isLoading || !user || isConnecting}
						onClick={connectSpotify}
						type="button"
					>
						{isConnecting ? (
							<Loader2 className="h-4 w-4 animate-spin" />
						) : (
							<Music2 className="h-4 w-4" />
						)}
						{isConnecting ? "Connecting…" : "Connect Spotify"}
					</Button>
				</CardContent>
			</Card>
		</main>
	);
}

export default function SpotifyOwnerPage() {
	return (
		<Suspense>
			<SpotifyOwnerContent />
		</Suspense>
	);
}
