"use client";

import { motion, type Variants } from "framer-motion";
import type { ReactNode } from "react";

const fadeUp: Variants = {
	hidden: { opacity: 0, y: 24 },
	visible: { opacity: 1, y: 0 },
};

interface AnimatedSectionProps {
	children: ReactNode;
	className?: string;
	delay?: number;
	id?: string;
}

export function AnimatedSection({
	children,
	className,
	delay = 0,
	id,
}: AnimatedSectionProps) {
	return (
		<motion.section
			className={className}
			id={id}
			initial="hidden"
			whileInView="visible"
			viewport={{ once: true, margin: "-80px" }}
			transition={{ duration: 0.55, delay, ease: [0.21, 0.47, 0.32, 0.98] }}
			variants={fadeUp}
		>
			{children}
		</motion.section>
	);
}

interface StaggerGridProps {
	children: ReactNode;
	className?: string;
}

export function StaggerGrid({ children, className }: StaggerGridProps) {
	return (
		<motion.div
			className={className}
			initial="hidden"
			whileInView="visible"
			viewport={{ once: true, margin: "-60px" }}
			variants={{
				hidden: {},
				visible: { transition: { staggerChildren: 0.08 } },
			}}
		>
			{children}
		</motion.div>
	);
}

export function StaggerItem({
	children,
	className,
}: {
	children: ReactNode;
	className?: string;
}) {
	return (
		<motion.div
			className={className}
			variants={{
				hidden: { opacity: 0, y: 20 },
				visible: {
					opacity: 1,
					y: 0,
					transition: { duration: 0.45, ease: [0.21, 0.47, 0.32, 0.98] },
				},
			}}
		>
			{children}
		</motion.div>
	);
}