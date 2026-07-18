import type { Metadata } from "next";
import type { ReactNode } from "react";

export const metadata: Metadata = {
	title: "Notification Channels",
};

export default function ChannelsLayout({ children }: { children: ReactNode }) {
	return children;
}