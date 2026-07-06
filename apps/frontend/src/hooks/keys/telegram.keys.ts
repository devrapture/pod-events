export const telegramKeys = {
	all: ["telegram"] as const,
	link: () => [...telegramKeys.all, "link"] as const,
};
