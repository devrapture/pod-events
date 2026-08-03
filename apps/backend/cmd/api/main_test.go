package main

import (
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestWithSentryLogging(t *testing.T) {
	transport := &sentry.MockTransport{}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:       "https://public@example.com/1",
		Transport: transport,
	})
	require.NoError(t, err)

	withSentryLogging(zap.NewNop(), false).Info("disabled log")
	withSentryLogging(zap.NewNop(), true).Info("enabled log")
	require.True(t, sentry.Flush(2*time.Second))

	events := transport.Events()
	require.Len(t, events, 1)
	require.Len(t, events[0].Logs, 1)
	assert.Equal(t, "enabled log", events[0].Logs[0].Body)
}

func TestScrubSentryEvent(t *testing.T) {
	event := sentry.NewEvent()
	event.Request = &sentry.Request{
		Cookies: "session=secret",
		Headers: map[string]string{
			"Authorization":                   "Bearer secret",
			"Cookie":                          "session=secret",
			"Content-Type":                    "application/json",
			"X-Cron-Secret":                   "cron-secret",
			"X-Telegram-Bot-Api-Secret-Token": "telegram-secret",
		},
	}

	result := scrubSentryEvent(event, nil)

	assert.Empty(t, result.Request.Cookies)
	assert.NotContains(t, result.Request.Headers, "Authorization")
	assert.NotContains(t, result.Request.Headers, "Cookie")
	assert.NotContains(t, result.Request.Headers, "X-Cron-Secret")
	assert.NotContains(t, result.Request.Headers, "X-Telegram-Bot-Api-Secret-Token")
	assert.Equal(t, "application/json", result.Request.Headers["Content-Type"])
}
