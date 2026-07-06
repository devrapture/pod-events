"use client";

import * as Dialog from "@radix-ui/react-dialog";
import { AnimatePresence, motion } from "framer-motion";

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
}: ConfirmDialogProps) {
	return (
		<Dialog.Root onOpenChange={onOpenChange} open={open}>
			<AnimatePresence>
				{open && (
					<Dialog.Portal forceMount>
						<Dialog.Overlay asChild>
							<motion.div
								animate={{ opacity: 1 }}
								className="fixed inset-0 z-50 bg-black/60 backdrop-blur-sm"
								exit={{ opacity: 0 }}
								initial={{ opacity: 0 }}
								transition={{ duration: 0.15 }}
							/>
						</Dialog.Overlay>
						<Dialog.Content asChild>
							<motion.div
								animate={{ opacity: 1, scale: 1, y: "-50%", x: "-50%" }}
								className="fixed top-1/2 left-1/2 z-50 w-full max-w-sm -translate-x-1/2 -translate-y-1/2 rounded-xl border border-white/[0.06] bg-zinc-950 p-6 shadow-2xl"
								exit={{ opacity: 0, scale: 0.95, y: "-50%", x: "-50%" }}
								initial={{ opacity: 0, scale: 0.95, y: "-50%", x: "-50%" }}
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
											className="inline-flex h-10 items-center justify-center rounded-lg border border-white/10 bg-white/5 px-5 font-medium text-sm text-zinc-100 transition-all duration-200 hover:border-white/20 hover:bg-white/10 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/50"
											type="button"
										>
											{cancelLabel}
										</button>
									</Dialog.Close>
									<button
										className={cn(
											"inline-flex h-10 items-center justify-center rounded-lg px-5 font-medium text-sm transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/50",
											variant === "destructive"
												? "bg-red-500 text-white shadow-lg shadow-red-500/20 hover:bg-red-400"
												: "bg-emerald-500 text-zinc-950 shadow-emerald-500/20 shadow-lg hover:bg-emerald-400",
										)}
										onClick={() => {
											onConfirm();
											onOpenChange(false);
										}}
										type="button"
									>
										{confirmLabel}
									</button>
								</div>
							</motion.div>
						</Dialog.Content>
					</Dialog.Portal>
				)}
			</AnimatePresence>
		</Dialog.Root>
	);
}
