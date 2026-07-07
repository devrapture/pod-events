import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

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
