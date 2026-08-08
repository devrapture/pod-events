// This file configures the initialization of Sentry for edge features (middleware, edge routes, and so on).
// The config you add here will be used whenever one of the edge features is loaded.
// Note that this config is unrelated to the Vercel Edge Runtime and is also required when running locally.
// https://docs.sentry.io/platforms/javascript/guides/nextjs/

import * as Sentry from "@sentry/nextjs";

const isProduction = process.env.NODE_ENV === "production";

Sentry.init({
	dsn: "https://d77a82b10f43fe75ddd6d0999f3d68d6@o4511846014517248.ingest.de.sentry.io/4511847524270160",

	// Match server config: disable performance tracing in local dev to avoid
	// MaxListenersExceededWarning noise from OTEL HTTP instrumentation under
	// Next 16 + Turbopack. See sentry.server.config.ts.
	tracesSampleRate: isProduction ? 0.1 : 0,

	// Enable logs to be sent to Sentry
	enableLogs: true,

	dataCollection: {
		// To disable sending user data and HTTP bodies, uncomment the lines below. For more info visit:
		// https://docs.sentry.io/platforms/javascript/guides/nextjs/configuration/options/#dataCollection
		// userInfo: false,
		// httpBodies: false,
	},
});
