import type { AxiosError, InternalAxiosRequestConfig } from "axios";
import axios from "axios";

import { TOKEN_KEY } from "./auth";

export const server = axios.create({
	baseURL: process.env.NEXT_PUBLIC_API_URL,
	responseType: "json",
	timeout: 240_000,
	withCredentials: true,
});

export const serverWithInterceptors = axios.create({
	baseURL: process.env.NEXT_PUBLIC_API_URL,
	responseType: "json",
	timeout: 240_000,
	withCredentials: true,
});

function getAuthToken(): string | null {
	if (typeof window === "undefined") return null;
	try {
		return localStorage.getItem(TOKEN_KEY);
	} catch {
		return null;
	}
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
