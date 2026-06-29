import { ImageResponse } from "next/og";

export const runtime = "edge";
export const alt = "PodEvents — Spotify Podcast Notifications";
export const size = { width: 1200, height: 630 };
export const contentType = "image/png";

export default function Image() {
	return new ImageResponse(
		(
			<div
				style={{
					width: "100%",
					height: "100%",
					display: "flex",
					flexDirection: "column",
					justifyContent: "center",
					padding: "80px",
					background: "linear-gradient(135deg, #09090b 0%, #18181b 50%, #0c1f17 100%)",
					fontFamily: "system-ui, sans-serif",
				}}
			>
				<div
					style={{
						display: "flex",
						alignItems: "center",
						gap: "16px",
						marginBottom: "32px",
					}}
				>
					<div
						style={{
							width: "48px",
							height: "48px",
							borderRadius: "12px",
							background: "rgba(16, 185, 129, 0.15)",
							border: "1px solid rgba(16, 185, 129, 0.3)",
							display: "flex",
							alignItems: "center",
							justifyContent: "center",
							fontSize: "24px",
						}}
					>
						📻
					</div>
					<span style={{ fontSize: "32px", fontWeight: 600, color: "#fafafa" }}>
						PodEvents
					</span>
				</div>
				<h1
					style={{
						fontSize: "56px",
						fontWeight: 700,
						color: "#fafafa",
						lineHeight: 1.15,
						margin: 0,
						maxWidth: "900px",
					}}
				>
					Spotify Podcast Notifications for Your Entire Team
				</h1>
				<p
					style={{
						fontSize: "24px",
						color: "#a1a1aa",
						marginTop: "24px",
						maxWidth: "800px",
						lineHeight: 1.5,
					}}
				>
					Monitor podcasts 24/7. Notify Slack, Discord, Telegram & WhatsApp instantly.
				</p>
			</div>
		),
		{ ...size },
	);
}