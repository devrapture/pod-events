export const healthKeys = {
	all: ["health"] as const,
	check: () => [...healthKeys.all, "check"] as const,
};
