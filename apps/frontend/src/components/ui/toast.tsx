"use client";

import { CheckCircle2, X, XCircle } from "lucide-react";
import {
	createContext,
	type ReactNode,
	useCallback,
	useContext,
	useMemo,
	useState,
} from "react";

import { cn } from "@/lib/utils";

export type ToastVariant = "success" | "error";

interface ToastItem {
	id: string;
	message: string;
	variant: ToastVariant;
}

interface ToastContextValue {
	toast: {
		success: (message: string) => void;
		error: (message: string) => void;
	};
}

const ToastContext = createContext<ToastContextValue | null>(null);

const TOAST_DURATION_MS = 4000;

export function ToastProvider({ children }: { children: ReactNode }) {
	const [toasts, setToasts] = useState<ToastItem[]>([]);

	const dismiss = useCallback((id: string) => {
		setToasts((current) => current.filter((toast) => toast.id !== id));
	}, []);

	const addToast = useCallback(
		(message: string, variant: ToastVariant) => {
			const id = crypto.randomUUID();
			setToasts((current) => [...current, { id, message, variant }]);
			window.setTimeout(() => dismiss(id), TOAST_DURATION_MS);
		},
		[dismiss],
	);

	const value = useMemo(
		() => ({
			toast: {
				success: (message: string) => addToast(message, "success"),
				error: (message: string) => addToast(message, "error"),
			},
		}),
		[addToast],
	);

	return (
		<ToastContext.Provider value={value}>
			{children}
			<div
				aria-live="polite"
				className="pointer-events-none fixed right-4 bottom-4 z-50 flex w-full max-w-sm flex-col gap-2"
			>
				{toasts.map((toast) => (
					<ToastMessage
						key={toast.id}
						onDismiss={() => dismiss(toast.id)}
						toast={toast}
					/>
				))}
			</div>
		</ToastContext.Provider>
	);
}

function ToastMessage({
	toast,
	onDismiss,
}: {
	toast: ToastItem;
	onDismiss: () => void;
}) {
	const Icon = toast.variant === "success" ? CheckCircle2 : XCircle;

	return (
		<div
			className={cn(
				"pointer-events-auto flex items-start gap-3 rounded-xl border px-4 py-3 shadow-lg backdrop-blur-xl",
				toast.variant === "success"
					? "border-emerald-500/20 bg-emerald-500/10 text-emerald-100"
					: "border-red-500/20 bg-red-500/10 text-red-100",
			)}
			role="status"
		>
			<Icon aria-hidden className="mt-0.5 h-4 w-4 shrink-0" />
			<p className="flex-1 text-sm">{toast.message}</p>
			<button
				aria-label="Dismiss notification"
				className="rounded-md p-1 text-current/70 transition-colors hover:text-current"
				onClick={onDismiss}
				type="button"
			>
				<X className="h-4 w-4" />
			</button>
		</div>
	);
}

export function useToast(): ToastContextValue {
	const context = useContext(ToastContext);
	if (!context) {
		throw new Error("useToast must be used within a ToastProvider");
	}
	return context;
}
