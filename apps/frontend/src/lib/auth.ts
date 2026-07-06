export const TOKEN_KEY = "pod_events_token";

export interface User {
	id: string;
	name: string;
	email: string;
	avatar_url: string;
	spotify_user_id: string;
}

export const LOGIN_URL = `${process.env.NEXT_PUBLIC_API_URL}/auth/spotify/login`;

export function getToken(): string | null {
	if (typeof window === "undefined") return null;
	return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string): void {
	localStorage.setItem(TOKEN_KEY, token);
	document.cookie = `${TOKEN_KEY}=${token}; path=/; max-age=86400; SameSite=Lax; Secure`;
}

export function removeToken(): void {
	localStorage.removeItem(TOKEN_KEY);
	document.cookie = `${TOKEN_KEY}=; path=/; max-age=0; SameSite=Lax; Secure`;
}

export async function getCurrentUser(token: string): Promise<User | null> {
	try {
		const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/me`, {
			headers: { Authorization: `Bearer ${token}` },
		});
		if (!res.ok) return null;
		const body = await res.json();
		return body.data as User;
	} catch {
		return null;
	}
}
