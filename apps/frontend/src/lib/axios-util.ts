import type { AxiosError, InternalAxiosRequestConfig } from "axios";
import axios from "axios";

import { getNgrokHeaders } from "./api-headers";
import { TOKEN_KEY } from "./auth";

export const server = axios.create({
	baseURL: process.env.NEXT_PUBLIC_API_URL,
	responseType: "json",
	timeout: 240_000,
});

export const serverWithInterceptors = axios.create({
	baseURL: process.env.NEXT_PUBLIC_API_URL,
	responseType: "json",
	timeout: 240_000,
});

function getAuthToken(): string | null {
	if (typeof window === "undefined") return null;
	try {
		return localStorage.getItem(TOKEN_KEY);
	} catch {
		return null;
	}
}

function applyNgrokHeaders(config: InternalAxiosRequestConfig) {
	for (const [key, value] of Object.entries(getNgrokHeaders())) {
		config.headers.set(key, value);
	}
	return config;
}

for (const instance of [server, serverWithInterceptors]) {
	instance.interceptors.request.use(applyNgrokHeaders, (error: AxiosError) =>
		Promise.reject(error),
	);
}

serverWithInterceptors.interceptors.request.use(
	(config: InternalAxiosRequestConfig) => {
		const token = getAuthToken();
		if (token) {
			config.headers.Authorization = `Bearer ${token}`;
		}
		return config;
	},
	(error: AxiosError) => Promise.reject(error),
);
