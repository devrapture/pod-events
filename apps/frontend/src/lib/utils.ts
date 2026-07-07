import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

import type { ChannelType } from "@/services/types";

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function formatNumber(value: number): string {
	return value.toLocaleString("en-US");
}

export function getTimeGreeting(name: string): string {
	const hour = new Date().getHours();
	const period =
		hour < 12 ? "Good morning" : hour < 17 ? "Good afternoon" : "Good evening";
	const firstName = name.trim().split(/\s+/)[0] ?? name;
	return `${period}, ${firstName}`;
}

const CHANNEL_TYPE_LABELS: Record<ChannelType, string> = {
	slack_webhook: "Slack",
	discord_webhook: "Discord",
	whatsapp: "WhatsApp",
	telegram: "Telegram",
};

export function getChannelTypeLabel(type: ChannelType): string {
	return CHANNEL_TYPE_LABELS[type];
}

function maskWebhookUrl(url: string): string {
	try {
		const parsed = new URL(url);
		const segments = parsed.pathname.split("/").filter(Boolean);
		const lastSegment = segments.at(-1) ?? "";
		const prefix = `${parsed.origin}${parsed.pathname.split("/").slice(0, 2).join("/")}`;
		if (!lastSegment) return `${prefix}/...`;
		return `${prefix}/.../${lastSegment.slice(-6)}`;
	} catch {
		if (url.length <= 12) return url;
		return `${url.slice(0, 8)}...${url.slice(-4)}`;
	}
}

export function maskDestination(
	type: ChannelType,
	destination: string,
): string {
	switch (type) {
		case "slack_webhook":
		case "discord_webhook":
			return maskWebhookUrl(destination);
		case "telegram":
			return "Connected";
		default:
			return destination;
	}
}
