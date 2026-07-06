export const subscriptionKeys = {
	all: ["subscriptions"] as const,
	list: () => ["subscriptions", "list"] as const,
	detail: (id: string) => ["subscriptions", "detail", id] as const,
};
