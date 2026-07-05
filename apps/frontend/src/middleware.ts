import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";

const TOKEN_KEY = "pod_events_token";

export function middleware(request: NextRequest) {
	const token = request.cookies.get(TOKEN_KEY)?.value;

	if (!token) {
		return NextResponse.redirect(new URL("/", request.url));
	}

	return NextResponse.next();
}

export const config = {
	matcher: "/dashboard/:path*",
};
