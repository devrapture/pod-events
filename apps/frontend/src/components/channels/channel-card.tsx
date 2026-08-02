"use client";

import type { LucideIcon } from "lucide-react";
import { Hash, Phone, Send } from "lucide-react";

import { DiscordIcon } from "@/components/icons/discord";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { getChannelTypeLabel, maskDestination } from "@/lib/utils";
import type { ChannelType, NotificationChannel } from "@/services/types";

const CHANNEL_ICONS: Record<ChannelType, LucideIcon | typeof DiscordIcon> = {
	slack_webhook: Hash,
	discord_webhook: DiscordIcon,
	whatsapp: Phone,
	telegram: Send,
};

interface ChannelCardProps {
	channel: NotificationChannel;
	onToggle: (channel: NotificationChannel) => void;
	onRemove: (id: string) => void;
}

function ChannelTypeIcon({
	type,
	className,
}: {
	type: ChannelType;
	className?: string;
}) {
	const Icon = CHANNEL_ICONS[type];

	if (type === "discord_webhook") {
		return <DiscordIcon className={className} />;
	}

	const LucideComponent = Icon as LucideIcon;
	return <LucideComponent className={className} />;
}

export function ChannelCard({ channel, onToggle, onRemove }: ChannelCardProps) {
	const channelName = getChannelTypeLabel(channel.channel_type);
	const isDiscord = channel.channel_type === "discord_webhook";

	return (
		<Card>
			<div className="flex flex-col gap-4 p-4 sm:flex-row sm:items-center sm:justify-between">
				<div className="flex min-w-0 items-start gap-3">
					<div
						className={
							isDiscord
								? "flex h-10 w-10 shrink-0 items-center justify-center rounded-lg border border-[#5865F2]/30 bg-[#5865F2]/10"
								: "flex h-10 w-10 shrink-0 items-center justify-center rounded-lg border border-white/10 bg-white/5"
						}
					>
						<ChannelTypeIcon
							className={
								isDiscord
									? "h-5 w-5 text-[#5865F2]"
									: "h-5 w-5 text-emerald-400"
							}
							type={channel.channel_type}
						/>
					</div>
					<div className="min-w-0">
						<div className="flex flex-wrap items-center gap-2">
							<h3 className="font-medium text-sm text-zinc-100">
								{channelName}
							</h3>
							<Badge variant={channel.is_active ? "emerald" : "outline"}>
								{channel.is_active ? "Active" : "Inactive"}
							</Badge>
						</div>
						{channel.label && (
							<p className="mt-0.5 text-sm text-zinc-300">{channel.label}</p>
						)}
						<p className="mt-1 truncate font-mono text-xs text-zinc-400">
							{maskDestination(channel.channel_type, channel.destination)}
						</p>
					</div>
				</div>

				<div className="flex shrink-0 flex-wrap items-center gap-2">
					<Button
						onClick={() => onToggle(channel)}
						size="sm"
						variant={channel.is_active ? "destructive" : "secondary"}
					>
						{channel.is_active ? "Disable" : "Enable"}
					</Button>
					<Button
						onClick={() => onRemove(channel.id)}
						size="sm"
						variant="ghost"
					>
						Remove
					</Button>
				</div>
			</div>
		</Card>
	);
}
