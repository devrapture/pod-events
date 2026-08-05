"use client";

import { Loader2 } from "lucide-react";
import { useCallback, useState } from "react";

import { Button } from "@/components/ui/button";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import {
	useCreateChannel,
	useGenerateTelegramLink,
} from "@/hooks/mutations/channel.mutations";
import { useToast } from "@/hooks/use-toast";
import { getAPIErrorMessage } from "@/lib/api-error";
import { cn } from "@/lib/utils";
import type { ChannelType } from "@/services/types";

type UiChannelType = "slack" | "discord" | "whatsapp" | "telegram";

const CHANNEL_OPTIONS: { key: UiChannelType; label: string }[] = [
	{ key: "slack", label: "Slack" },
	{ key: "discord", label: "Discord" },
	{ key: "whatsapp", label: "WhatsApp" },
	{ key: "telegram", label: "Telegram" },
];

const UI_TO_API_TYPE: Record<UiChannelType, ChannelType | null> = {
	slack: "slack_webhook",
	discord: "discord_webhook",
	whatsapp: "whatsapp",
	telegram: null,
};

const INPUT_CLASS =
	"h-10 w-full rounded-lg border border-white/10 bg-white/5 px-3 text-sm text-zinc-100 placeholder:text-zinc-500 focus:border-emerald-500/40 focus:outline-none focus:ring-2 focus:ring-emerald-500/20";

interface AddChannelFormProps {
	onSuccess: () => void;
}

export function AddChannelForm({ onSuccess }: AddChannelFormProps) {
	const [selectedType, setSelectedType] = useState<UiChannelType>("slack");
	const [destination, setDestination] = useState("");
	const [label, setLabel] = useState("");
	const [validationError, setValidationError] = useState<string | null>(null);
	const { toast } = useToast();

	const { mutate: createChannel, isPending: isCreating } = useCreateChannel();
	const { mutate: generateTelegramLink, isPending: isConnectingTelegram } =
		useGenerateTelegramLink();

	const resetForm = useCallback(() => {
		setDestination("");
		setLabel("");
		setValidationError(null);
	}, []);

	const handleSave = useCallback(() => {
		const channelType = UI_TO_API_TYPE[selectedType];
		if (!channelType) return;

		const trimmed = destination.trim();
		if (!trimmed) {
			setValidationError("This field is required.");
			return;
		}

		if (channelType === "slack_webhook" || channelType === "discord_webhook") {
			try {
				const parsed = new URL(trimmed);
				if (parsed.protocol !== "https:") {
					setValidationError("Webhook URL must use https.");
					return;
				}
			} catch {
				setValidationError("Please enter a valid webhook URL.");
				return;
			}
		}

		if (channelType === "whatsapp") {
			if (!/^\[1-9]\d{1,14}$/.test(trimmed)) {
				setValidationError("Use E.164 format, e.g. 2348012345678.");
				return;
			}
		}
		setValidationError(null);
		createChannel(
			{
				channel_type: channelType,
				destination: trimmed,
				label: label.trim() || undefined,
			},
			{
				onSuccess: () => {
					toast.success("Notification channel saved");
					resetForm();
					onSuccess();
				},
				onError: (error) => {
					toast.error(
						getAPIErrorMessage(
							error,
							"Failed to save channel. Please try again.",
						),
					);
				},
			},
		);
	}, [
		selectedType,
		destination,
		label,
		createChannel,
		toast,
		resetForm,
		onSuccess,
	]);

	const handleConnectTelegram = useCallback(() => {
		generateTelegramLink(void 0, {
			onSuccess: (response) => {
				const url = response.data?.url;
				if (!url) {
					toast.error("Failed to generate Telegram link. Please try again.");
					return;
				}
				window.open(url, "_blank", "noopener,noreferrer");
				const win = window.open(url, "_blank", "noopener,noreferrer");
				if (!win) {
					toast.error(
						"Popup blocked. Please allow popups to connect Telegram.",
					);
				}
			},
			onError: (error) => {
				toast.error(
					getAPIErrorMessage(
						error,
						"Failed to connect Telegram. Please try again.",
					),
				);
			},
		});
	}, [generateTelegramLink, toast]);

	const handleTypeChange = (type: UiChannelType) => {
		setSelectedType(type);
		setDestination("");
		setLabel("");
		setValidationError(null);
	};

	const isPending = isCreating || isConnectingTelegram;

	return (
		<Card className="mb-6">
			<CardHeader>
				<CardTitle>Add Notification Channel</CardTitle>
				<CardDescription>
					Choose a channel type and configure where to receive alerts.
				</CardDescription>
			</CardHeader>
			<CardContent className="space-y-6">
				<div className="flex flex-wrap gap-2">
					{CHANNEL_OPTIONS.map((option) => (
						<Button
							key={option.key}
							onClick={() => handleTypeChange(option.key)}
							size="sm"
							type="button"
							variant={selectedType === option.key ? "default" : "secondary"}
						>
							{option.label}
						</Button>
					))}
				</div>

				{selectedType === "slack" && (
					<div className="space-y-3">
						<div>
							<label
								className="mb-1.5 block text-sm text-zinc-300"
								htmlFor="slack-label"
							>
								Channel label <span className="text-zinc-500">(optional)</span>
							</label>
							<input
								className={INPUT_CLASS}
								id="slack-label"
								onChange={(e) => setLabel(e.target.value)}
								placeholder="e.g. Work alerts — helps you identify this channel"
								type="text"
								value={label}
							/>
						</div>
						<div>
							<label
								className="mb-1.5 block text-sm text-zinc-300"
								htmlFor="slack-webhook"
							>
								Slack webhook URL
							</label>
							<input
								className={cn(
									INPUT_CLASS,
									validationError && "border-red-500/50",
								)}
								id="slack-webhook"
								onChange={(e) => {
									setDestination(e.target.value);
									setValidationError(null);
								}}
								placeholder="https://hooks.slack.com/services/..."
								type="url"
								value={destination}
							/>
							{validationError && (
								<p className="mt-1.5 text-red-400 text-xs">{validationError}</p>
							)}
						</div>
						<p className="text-xs text-zinc-500 leading-relaxed">
							Create an Incoming Webhook in your Slack workspace, copy the
							webhook URL, and paste it here.
						</p>
						<Button disabled={isPending} onClick={handleSave} type="button">
							{isCreating && <Loader2 className="h-4 w-4 animate-spin" />}
							Save Channel
						</Button>
					</div>
				)}

				{selectedType === "discord" && (
					<div className="space-y-3">
						<div>
							<label
								className="mb-1.5 block text-sm text-zinc-300"
								htmlFor="discord-label"
							>
								Channel label <span className="text-zinc-500">(optional)</span>
							</label>
							<input
								className={INPUT_CLASS}
								id="discord-label"
								onChange={(e) => setLabel(e.target.value)}
								placeholder="e.g. Podcast alerts — helps you identify this channel"
								type="text"
								value={label}
							/>
						</div>
						<div>
							<label
								className="mb-1.5 block text-sm text-zinc-300"
								htmlFor="discord-webhook"
							>
								Discord webhook URL
							</label>
							<input
								className={cn(
									INPUT_CLASS,
									validationError && "border-red-500/50",
								)}
								id="discord-webhook"
								onChange={(e) => {
									setDestination(e.target.value);
									setValidationError(null);
								}}
								placeholder="https://discord.com/api/webhooks/..."
								type="url"
								value={destination}
							/>
							{validationError && (
								<p className="mt-1.5 text-red-400 text-xs">{validationError}</p>
							)}
						</div>
						<p className="text-xs text-zinc-500 leading-relaxed">
							Get this from Discord Server Settings → Integrations → Webhooks.
						</p>
						<Button disabled={isPending} onClick={handleSave} type="button">
							{isCreating && <Loader2 className="h-4 w-4 animate-spin" />}
							Save Channel
						</Button>
					</div>
				)}

				{selectedType === "whatsapp" && (
					<div className="space-y-3">
						<div>
							<label
								className="mb-1.5 block text-sm text-zinc-300"
								htmlFor="whatsapp-contact"
							>
								WhatsApp phone number
							</label>
							<input
								className={cn(
									INPUT_CLASS,
									validationError && "border-red-500/50",
								)}
								id="whatsapp-contact"
								onChange={(e) => {
									setDestination(e.target.value);
									setValidationError(null);
								}}
								placeholder="+2348012345678"
								type="tel"
								value={destination}
							/>
							{validationError && (
								<p className="mt-1.5 text-red-400 text-xs">{validationError}</p>
							)}
						</div>
						<p className="text-xs text-zinc-500 leading-relaxed">
							Use E.164 format with country code, e.g. +2348012345678.
						</p>
						<Button disabled={isPending} onClick={handleSave} type="button">
							{isCreating && <Loader2 className="h-4 w-4 animate-spin" />}
							Save Channel
						</Button>
					</div>
				)}

				{selectedType === "telegram" && (
					<div className="space-y-3">
						<p className="text-sm text-zinc-400 leading-relaxed">
							Connect your Telegram account to receive new episode alerts.
							Telegram will open in a new tab so you can start our bot.
						</p>
						<Button
							disabled={isPending}
							onClick={handleConnectTelegram}
							type="button"
						>
							{isConnectingTelegram && (
								<Loader2 className="h-4 w-4 animate-spin" />
							)}
							Connect Telegram
						</Button>
					</div>
				)}
			</CardContent>
		</Card>
	);
}
