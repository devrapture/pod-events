"use client";

import { Loader2, Search, X } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

interface ImportToolbarProps {
	searchQuery: string;
	onSearchChange: (value: string) => void;
	isSearching: boolean;
	selectedCount: number;
	untrackedCount: number;
	onSelectAllUntracked: () => void;
	onClearSelection: () => void;
}

export function ImportToolbar({
	searchQuery,
	onSearchChange,
	isSearching,
	selectedCount,
	untrackedCount,
	onSelectAllUntracked,
	onClearSelection,
}: ImportToolbarProps) {
	return (
		<div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
			<div className="relative max-w-md flex-1">
				<Search
					aria-hidden
					className="pointer-events-none absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-zinc-500"
				/>
				<input
					aria-label="Search saved podcasts"
					className={cn(
						"h-10 w-full rounded-lg border border-white/10 bg-white/5 pr-10 pl-9 text-sm text-zinc-100",
						"placeholder:text-zinc-500 focus:border-emerald-500/40 focus:outline-none focus:ring-2 focus:ring-emerald-500/20",
					)}
					onChange={(e) => onSearchChange(e.target.value)}
					placeholder="Search your saved podcasts..."
					type="search"
					value={searchQuery}
				/>
				{searchQuery && (
					<button
						aria-label="Clear search"
						className="absolute top-1/2 right-3 -translate-y-1/2 text-zinc-500 hover:text-zinc-300"
						onClick={() => onSearchChange("")}
						type="button"
					>
						<X className="h-4 w-4" />
					</button>
				)}
				{isSearching && (
					<div className="mt-2 flex items-center gap-2 text-xs text-zinc-500">
						<Loader2 className="h-3 w-3 animate-spin" />
						Searching...
					</div>
				)}
			</div>

			<div className="flex flex-wrap items-center gap-2">
				{selectedCount > 0 && (
					<Badge aria-live="polite" variant="emerald">
						{selectedCount} selected
					</Badge>
				)}
				<Button
					disabled={untrackedCount === 0}
					onClick={onSelectAllUntracked}
					size="sm"
					variant="secondary"
				>
					Select all untracked
				</Button>
				<Button
					disabled={selectedCount === 0}
					onClick={onClearSelection}
					size="sm"
					variant="ghost"
				>
					Clear selection
				</Button>
			</div>
		</div>
	);
}
