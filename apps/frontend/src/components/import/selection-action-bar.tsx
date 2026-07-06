"use client";

import { Loader2 } from "lucide-react";

import { Button } from "@/components/ui/button";

interface SelectionActionBarProps {
	selectedCount: number;
	isPending: boolean;
	onTrack: () => void;
}

export function SelectionActionBar({
	selectedCount,
	isPending,
	onTrack,
}: SelectionActionBarProps) {
	if (selectedCount === 0) return null;

	return (
		<div className="fixed right-0 bottom-0 left-0 z-40 border-white/10 border-t bg-zinc-950/90 px-4 py-4 backdrop-blur-xl md:left-64">
			<div className="mx-auto flex max-w-6xl items-center justify-between gap-4">
				<p className="text-sm text-zinc-300">
					<span className="font-medium text-zinc-100">{selectedCount}</span>{" "}
					podcast{selectedCount === 1 ? "" : "s"} selected
				</p>
				<Button disabled={isPending || selectedCount === 0} onClick={onTrack}>
					{isPending ? (
						<>
							<Loader2 className="h-4 w-4 animate-spin" />
							Tracking...
						</>
					) : (
						"Track selected podcasts"
					)}
				</Button>
			</div>
		</div>
	);
}
