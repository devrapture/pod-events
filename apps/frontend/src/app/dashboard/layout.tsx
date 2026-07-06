import type { Metadata } from "next";

import { DashboardShell } from "@/components/dashboard-shell";

export const metadata: Metadata = {
	title: "Dashboard",
};

export default function DashboardLayout({
	children,
}: Readonly<{ children: React.ReactNode }>) {
	return <DashboardShell>{children}</DashboardShell>;
}
