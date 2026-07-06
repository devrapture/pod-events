import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";

import { TOKEN_KEY } from "@/lib/auth";

export function middleware(request: NextRequest) {
	const token = request.cookies.get(TOKEN_KEY)?.value;
	const { pathname } = request.nextUrl;

	if (pathname.startsWith("/dashboard")) {
		if (!token) {
			return NextResponse.redirect(new URL("/", request.url));
		}
		return NextResponse.next();
	}

	if (pathname === "/" && token) {
		return NextResponse.redirect(new URL("/dashboard", request.url));
	}

	return NextResponse.next();
}

export const config = {
	matcher: ["/", "/dashboard/:path*"],
};
