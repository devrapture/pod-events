import axios from "axios";

import type { APIResponse } from "@/services/types";

export function getAPIErrorMessage(error: unknown, fallback: string): string {
	if (!axios.isAxiosError<APIResponse>(error)) {
		return fallback;
	}

	const apiError = error.response?.data?.error;
	const detail = Object.values(apiError?.details ?? {})[0];

	if (detail) {
		return detail;
	}

	if (apiError?.message) {
		return apiError.message;
	}

	if (!error.response) {
		return "Could not reach the server. Check your connection and try again.";
	}

	return fallback;
}
