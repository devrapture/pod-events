import Link from "next/link";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";

// Kept as a compatibility boundary for old imports. Personal Spotify-library
// import is intentionally disabled now that users authenticate with Google.
export function ImportSpotifyPage() {
	return (
		<Card>
			<CardContent className="flex flex-col items-center gap-4 p-8 text-center">
				<p className="font-medium text-zinc-200">
					Spotify library import is no longer available.
				</p>
				<p className="max-w-md text-sm text-zinc-500">
					Search the Spotify catalog and subscribe to the podcasts you want
					PodEvents to monitor.
				</p>
				<Button asChild>
					<Link href="/dashboard/search">Search podcasts</Link>
				</Button>
			</CardContent>
		</Card>
	);
}
