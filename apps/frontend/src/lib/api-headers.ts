export function getNgrokHeaders(): Record<string, string> {
	const apiUrl = process.env.NEXT_PUBLIC_API_URL ?? "";
	if (apiUrl.includes("ngrok")) {
		return { "ngrok-skip-browser-warning": "true" };
	}
	return {};
}