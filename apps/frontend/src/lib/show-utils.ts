import type {
	SavedShowResponse,
	SavedShowWithTracking,
	SubscriptionResponse,
	TrackingStatus,
} from "@/services/types";

export function getShowTrackingStatus(
	show: SavedShowResponse,
	trackedSpotifyIds: Set<string>,
): { is_tracked: boolean; tracking_status: TrackingStatus } {
	if (show.is_tracked !== undefined || show.tracking_status !== undefined) {
		const isTracked = show.is_tracked ?? show.tracking_status === "tracked";
		return {
			is_tracked: isTracked,
			tracking_status:
				show.tracking_status ?? (isTracked ? "tracked" : "untracked"),
		};
	}

	const isTracked = trackedSpotifyIds.has(show.id);
	return {
		is_tracked: isTracked,
		tracking_status: isTracked ? "tracked" : "untracked",
	};
}

export function enrichShowsWithTracking(
	shows: SavedShowResponse[],
	subscriptions: SubscriptionResponse[],
): SavedShowWithTracking[] {
	const trackedSpotifyIds = new Set(
		subscriptions.map((sub) => sub.podcast_show.spotify_show_id),
	);

	return shows.map((show) => {
		const tracking = getShowTrackingStatus(show, trackedSpotifyIds);
		return { ...show, ...tracking };
	});
}
