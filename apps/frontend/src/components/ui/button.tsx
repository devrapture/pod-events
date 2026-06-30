import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";
import type { ButtonHTMLAttributes } from "react";

import { cn } from "@/lib/utils";

const buttonVariants = cva(
	"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-lg text-sm font-medium transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/50 disabled:pointer-events-none disabled:opacity-50",
	{
		variants: {
			variant: {
				default:
					"bg-emerald-500 text-zinc-950 shadow-lg shadow-emerald-500/20 hover:bg-emerald-400 hover:shadow-emerald-500/30",
				secondary:
					"border border-white/10 bg-white/5 text-zinc-100 backdrop-blur-sm hover:border-white/20 hover:bg-white/10",
				ghost: "text-zinc-400 hover:bg-white/5 hover:text-zinc-100",
				outline:
					"border border-white/15 bg-transparent text-zinc-100 hover:border-emerald-500/40 hover:bg-emerald-500/5",
			},
			size: {
				default: "h-10 px-5 py-2",
				sm: "h-8 rounded-md px-3 text-xs",
				lg: "h-12 rounded-xl px-8 text-base",
				icon: "h-10 w-10",
			},
		},
		defaultVariants: {
			variant: "default",
			size: "default",
		},
	},
);

export interface ButtonProps
	extends ButtonHTMLAttributes<HTMLButtonElement>,
		VariantProps<typeof buttonVariants> {
	asChild?: boolean;
}

export function Button({
	className,
	variant,
	size,
	asChild = false,
	...props
}: ButtonProps) {
	const Comp = asChild ? Slot : "button";
	return (
		<Comp
			className={cn(buttonVariants({ variant, size, className }))}
			{...props}
		/>
	);
}