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
	setToken: (token: string, user?: User) => Promise<boolean>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
	const [user, setUser] = useState<User | null>(null);
	const [isLoading, setIsLoading] = useState(true);
	const [isLoggingIn, setIsLoggingIn] = useState(false);

	const redirectHome = useCallback(() => {
		if (window.location.pathname !== "/") {
			window.location.replace("/");
		}
	}, []);

	const login = useCallback(() => {
		setIsLoggingIn(true);
		window.location.href = LOGIN_URL;
	}, []);

	const logout = useCallback(() => {
		removeToken();
		setUser(null);
		window.location.href = "/";
	}, []);

	const setToken = useCallback(
		async (token: string, authenticatedUser?: User) => {
			storeToken(token);
			if (authenticatedUser) {
				setUser(authenticatedUser);
				return true;
			}

			const fetchedUser = await getCurrentUser(token);
			if (!fetchedUser) {
				removeToken();
				setUser(null);
				redirectHome();
				return false;
			}
			setUser(fetchedUser);
			return true;
		},
		[redirectHome],
	);

	useEffect(() => {
		const token = getToken();
		if (!token) {
			setIsLoading(false);
			return;
		}

		getCurrentUser(token)
			.then((fetchedUser) => {
				if (!fetchedUser) {
					removeToken();
					redirectHome();
				}
				setUser(fetchedUser);
				setIsLoading(false);
			})
			.catch(() => {
				removeToken();
				redirectHome();
				setIsLoading(false);
			});
	}, [redirectHome]);

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
