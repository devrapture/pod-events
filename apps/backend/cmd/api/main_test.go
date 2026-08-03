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

	client := sentry.CurrentHub().Client()
	require.NotNil(t, client)
	t.Cleanup(func() {
		client.Close()
		sentry.CurrentHub().BindClient(nil)
	})

	withSentryLogging(zap.NewNop(), false).Info("disabled log")
	sentryLogger := withSentryLogging(zap.NewNop(), true)
	sentryLogger.Debug("debug log")
	sentryLogger.Info("info log")
	sentryLogger.Warn("warn log")
	sentryLogger.Error("error log")
	require.True(t, sentry.Flush(2*time.Second))

	events := transport.Events()
	require.Len(t, events, 1)
	require.Len(t, events[0].Logs, 3)

	expectedLogs := []struct {
		level   sentry.LogLevel
		message string
	}{
		{level: sentry.LogLevelInfo, message: "info log"},
		{level: sentry.LogLevelWarn, message: "warn log"},
		{level: sentry.LogLevelError, message: "error log"},
	}
	for i, expected := range expectedLogs {
		assert.Equal(t, expected.level, events[0].Logs[i].Level)
		assert.Equal(t, expected.message, events[0].Logs[i].Body)
	}
	for _, log := range events[0].Logs {
		assert.NotEqual(t, "debug log", log.Body)
	}
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
