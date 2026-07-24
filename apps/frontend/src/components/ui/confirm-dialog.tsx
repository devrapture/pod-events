"use client";

import * as Dialog from "@radix-ui/react-dialog";
import { AnimatePresence, motion } from "framer-motion";
import { Loader2 } from "lucide-react";

import { cn } from "@/lib/utils";

interface ConfirmDialogProps {
	open: boolean;
	onOpenChange: (open: boolean) => void;
	title: string;
	description: string;
	onConfirm: () => void;
	confirmLabel?: string;
	cancelLabel?: string;
	variant?: "default" | "destructive";
	isPending?: boolean;
	closeOnConfirm?: boolean;
}

export function ConfirmDialog({
	open,
	onOpenChange,
	title,
	description,
	onConfirm,
	confirmLabel = "Confirm",
	cancelLabel = "Cancel",
	variant = "default",
	isPending = false,
	closeOnConfirm = true,
}: ConfirmDialogProps) {
	const handleOpenChange = (nextOpen: boolean) => {
		if (isPending && !nextOpen) return;
		onOpenChange(nextOpen);
	};

	return (
		<Dialog.Root onOpenChange={handleOpenChange} open={open}>
			<AnimatePresence>
				{open && (
					<Dialog.Portal forceMount>
						<Dialog.Overlay asChild>
							<motion.div
								animate={{ opacity: 1 }}
								className="fixed inset-0 z-50 bg-black/60 backdrop-blur-sm md:left-64"
								exit={{ opacity: 0 }}
								initial={{ opacity: 0 }}
								transition={{ duration: 0.15 }}
							/>
						</Dialog.Overlay>
						<Dialog.Content asChild>
							<motion.div
								animate={{ opacity: 1 }}
								className="pointer-events-none fixed inset-0 z-50 flex items-center justify-center p-4 md:left-64"
								exit={{ opacity: 0 }}
								initial={{ opacity: 0 }}
								transition={{ duration: 0.15 }}
							>
								<motion.div
									animate={{ opacity: 1, scale: 1 }}
									className="pointer-events-auto w-full max-w-sm rounded-xl border border-white/[0.06] bg-zinc-950 p-6 shadow-2xl"
									exit={{ opacity: 0, scale: 0.95 }}
									initial={{ opacity: 0, scale: 0.95 }}
									transition={{ duration: 0.15 }}
								>
									<Dialog.Title className="font-semibold text-lg text-zinc-100">
										{title}
									</Dialog.Title>
									<Dialog.Description className="mt-2 text-sm text-zinc-400">
										{description}
									</Dialog.Description>

									<div className="mt-6 flex justify-end gap-3">
										<Dialog.Close asChild>
											<button
												className="inline-flex h-10 items-center justify-center rounded-lg border border-white/10 bg-white/5 px-5 font-medium text-sm text-zinc-100 transition-all duration-200 hover:border-white/20 hover:bg-white/10 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/50 disabled:cursor-not-allowed disabled:opacity-50"
												disabled={isPending}
												type="button"
											>
												{cancelLabel}
											</button>
										</Dialog.Close>
										<button
											className={cn(
												"inline-flex h-10 items-center justify-center gap-2 rounded-lg px-5 font-medium text-sm transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/50 disabled:cursor-not-allowed disabled:opacity-50",
												variant === "destructive"
													? "bg-red-500 text-white shadow-lg shadow-red-500/20 hover:bg-red-400"
													: "bg-emerald-500 text-zinc-950 shadow-emerald-500/20 shadow-lg hover:bg-emerald-400",
											)}
											disabled={isPending}
											onClick={() => {
												onConfirm();
												if (closeOnConfirm) {
													onOpenChange(false);
												}
											}}
											type="button"
										>
											{isPending && (
												<Loader2 className="h-4 w-4 animate-spin" />
											)}
											{confirmLabel}
										</button>
									</div>
								</motion.div>
							</motion.div>
						</Dialog.Content>
					</Dialog.Portal>
				)}
			</AnimatePresence>
		</Dialog.Root>
	);
}
