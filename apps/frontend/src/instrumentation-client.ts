// This file configures the initialization of Sentry on the client.
// The added config here will be used whenever a users loads a page in their browser.
// https://docs.sentry.io/platforms/javascript/guides/nextjs/

import * as Sentry from "@sentry/nextjs";

Sentry.init({
	dsn: "https://d77a82b10f43fe75ddd6d0999f3d68d6@o4511846014517248.ingest.de.sentry.io/4511847524270160",

	// Add optional integrations for additional features
	integrations: [Sentry.replayIntegration()],

	// Match server/edge config: disable performance tracing in local dev to
	// avoid MaxListenersExceededWarning noise from OTEL HTTP instrumentation.
	// See sentry.server.config.ts.
	tracesSampleRate: process.env.NODE_ENV === "production" ? 0.1 : 0,
	// Enable logs to be sent to Sentry
	enableLogs: true,

	// Define how likely Replay events are sampled.
	// This sets the sample rate to be 10%. You may want this to be 100% while
	// in development and sample at a lower rate in production
	replaysSessionSampleRate: 0.1,

	// Define how likely Replay events are sampled when an error occurs.
	replaysOnErrorSampleRate: 1.0,

	dataCollection: {
		// To disable sending user data and HTTP bodies, uncomment the lines below. For more info visit:
		// https://docs.sentry.io/platforms/javascript/guides/nextjs/configuration/options/#dataCollection
		// userInfo: false,
		// httpBodies: [],
	},
});

export const onRouterTransitionStart = Sentry.captureRouterTransitionStart;
