export const dashboardKeys = {
	all: ["dashboard"] as const,
	summary: () => ["dashboard", "summary"] as const,
};
