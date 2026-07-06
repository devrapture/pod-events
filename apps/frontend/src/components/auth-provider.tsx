"use client";

import type { ReactNode } from "react";
import {
	createContext,
	useCallback,
	useContext,
	useEffect,
	useState,
} from "react";

import {
	getCurrentUser,
	getToken,
	LOGIN_URL,
	removeToken,
	setToken as storeToken,
	type User,
} from "@/lib/auth";

interface AuthContextValue {
	user: User | null;
	isLoading: boolean;
	isLoggingIn: boolean;
	isAuthenticated: boolean;
	login: () => void;
	logout: () => void;
	setToken: (token: string) => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
	const [user, setUser] = useState<User | null>(null);
	const [isLoading, setIsLoading] = useState(true);
	const [isLoggingIn, setIsLoggingIn] = useState(false);

	const login = useCallback(() => {
		setIsLoggingIn(true);
		window.location.href = LOGIN_URL;
	}, []);

	const logout = useCallback(() => {
		removeToken();
		setUser(null);
		window.location.href = "/";
	}, []);

	const setToken = useCallback(async (token: string) => {
		storeToken(token);
		const fetchedUser = await getCurrentUser(token);
		setUser(fetchedUser);
	}, []);

	useEffect(() => {
		const token = getToken();
		if (!token) {
			setIsLoading(false);
			return;
		}

		getCurrentUser(token)
			.then((fetchedUser) => {
				setUser(fetchedUser);
				setIsLoading(false);
			})
			.catch(() => {
				removeToken();
				setIsLoading(false);
			});
	}, []);

	return (
		<AuthContext.Provider
			value={{
				user,
				isLoading,
				isLoggingIn,
				isAuthenticated: !!user,
				login,
				logout,
				setToken,
			}}
		>
			{children}
		</AuthContext.Provider>
	);
}

export function useAuth(): AuthContextValue {
	const context = useContext(AuthContext);
	if (!context) {
		throw new Error("useAuth must be used within an AuthProvider");
	}
	return context;
}
