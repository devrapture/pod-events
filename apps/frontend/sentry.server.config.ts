// This file configures the initialization of Sentry on the server.
// The config you add here will be used whenever the server handles a request.
// https://docs.sentry.io/platforms/javascript/guides/nextjs/

import * as Sentry from "@sentry/nextjs";

const isProduction = process.env.NODE_ENV === "production";

Sentry.init({
	dsn: "https://d77a82b10f43fe75ddd6d0999f3d68d6@o4511846014517248.ingest.de.sentry.io/4511847524270160",

	// Full OTEL/HTTP instrumentation in Next 16 + Turbopack dev stacks many
	// ServerResponse "close" listeners and triggers MaxListenersExceededWarning.
	// Keep production sampling; disable performance tracing in local dev.
	// See: https://github.com/getsentry/sentry-javascript/issues/19367
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
