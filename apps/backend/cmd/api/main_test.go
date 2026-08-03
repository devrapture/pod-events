package main

import (
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
)

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
