import "@/styles/globals.css";

import type { Metadata } from "next";
import { Geist } from "next/font/google";
import Script from "next/script";

import { Providers } from "@/components/providers";
import { SITE_URL } from "@/lib/constants";

const geist = Geist({
	subsets: ["latin"],
	variable: "--font-geist-sans",
});

export const metadata: Metadata = {
	title: {
		default: "PodEvents — Spotify Podcast Notifications for Your Entire Team",
		template: "%s | PodEvents",
	},
	description:
		"PodEvents monitors Spotify podcasts 24/7 and instantly delivers new episode notifications to Slack, Discord, Telegram, and WhatsApp.",
	metadataBase: new URL(SITE_URL),
	openGraph: {
		type: "website",
		locale: "en_US",
		url: SITE_URL,
		siteName: "PodEvents",
		title: "PodEvents — Spotify Podcast Notifications for Your Entire Team",
		description:
			"PodEvents monitors Spotify podcasts 24/7 and instantly delivers new episode notifications to Slack, Discord, Telegram, and WhatsApp.",
		images: [
			{
				url: "/opengraph-image",
				width: 1200,
				height: 630,
				alt: "PodEvents",
			},
		],
	},
	twitter: {
		card: "summary_large_image",
		title: "PodEvents — Spotify Podcast Notifications for Your Entire Team",
		description:
			"PodEvents monitors Spotify podcasts 24/7 and instantly delivers new episode notifications to Slack, Discord, Telegram, and WhatsApp.",
		images: ["/opengraph-image"],
	},
	robots: {
		index: true,
		follow: true,
	},
	keywords: [
		"podcast notifications",
		"Spotify",
		"Slack",
		"Discord",
		"Telegram",
		"WhatsApp",
		"open source",
		"podcast monitoring",
		"engineering podcasts",
	],
};

export default function RootLayout({
	children,
}: Readonly<{ children: React.ReactNode }>) {
	return (
		<html
			className={`${geist.variable} dark`}
			lang="en"
			suppressHydrationWarning
		>
			<body className="font-sans">
				<Providers>{children}</Providers>
				{process.env.NODE_ENV === "production" && (
					<Script
						async
						data-domain="pod-event.vercel.app"
						data-site="ekewhwo941ka"
						src="https://www.sabilytics.com/script.js"
						strategy="afterInteractive"
					/>
				)}
			</body>
		</html>
	);
}
