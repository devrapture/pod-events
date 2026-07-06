import type { Metadata } from "next";
import type { ReactNode } from "react";

export const metadata: Metadata = {
	title: "My Subscriptions",
};

export default function SubscriptionsLayout({
	children,
}: {
	children: ReactNode;
}) {
	return children;
}