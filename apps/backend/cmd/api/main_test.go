package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/devrapture/pod-events/internal/config"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRunClosesDatabaseWhenServerFails(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("PORT", "invalid")
	t.Setenv("DATABASE_URL", "postgres://unused")
	t.Setenv("TOKEN_ENCRYPTION_KEY", "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXoxMjM0NTY=")
	t.Setenv("SPOTIFY_CLIENT_SECRET", "test")
	t.Setenv("SPOTIFY_CLIENT_ID", "test")
	t.Setenv("SPOTIFY_REDIRECT_URL", "http://localhost/callback")
	t.Setenv("FRONTEND_URL", "http://localhost")
	t.Setenv("JWT_SECRET", "test")
	t.Setenv("TELEGRAM_BOT_TOKEN", "test")
	t.Setenv("TELEGRAM_WEBHOOK_SECRET", "test")
	t.Setenv("BOT_NAME", "test")
	t.Setenv("CRON_SECRET", "test")
	t.Setenv("SENTRY_DSN", "")
	t.Setenv("SENTRY_TRACES_SAMPLE_RATE", "0")
	t.Setenv("SENTRY_ENABLE_LOGS", "false")

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	mock.ExpectClose()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	require.NoError(t, err)

	originalConnectDatabase := connectDatabase
	connectDatabase = func(*config.Config) (*gorm.DB, error) {
		return gormDB, nil
	}
	t.Cleanup(func() {
		connectDatabase = originalConnectDatabase
	})

	exitStatus := run()

	assert.Equal(t, 1, exitStatus)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRunFlushesSentryBeforeExit(t *testing.T) {
	if os.Getenv("POD_EVENTS_SENTRY_EXIT_HELPER") == "1" {
		os.Exit(run())
	}

	envelopes := make(chan []byte, 4)
	ingestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		envelopes <- body
		w.WriteHeader(http.StatusOK)
	}))
	defer ingestServer.Close()

	dsn, err := url.Parse(ingestServer.URL)
	require.NoError(t, err)
	dsn.User = url.User("public")
	dsn.Path = "/1"

	cmd := exec.Command(os.Args[0], "-test.run=^TestRunFlushesSentryBeforeExit$")
	cmd.Env = append(os.Environ(),
		"POD_EVENTS_SENTRY_EXIT_HELPER=1",
		"APP_ENV=test",
		"DATABASE_URL=postgres://postgres:postgres@127.0.0.1:1/podevents?connect_timeout=1",
		"TOKEN_ENCRYPTION_KEY=YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXoxMjM0NTY=",
		"SPOTIFY_CLIENT_SECRET=test",
		"SPOTIFY_CLIENT_ID=test",
		"SPOTIFY_REDIRECT_URL=http://localhost/callback",
		"FRONTEND_URL=http://localhost",
		"JWT_SECRET=test",
		"TELEGRAM_BOT_TOKEN=test",
		"TELEGRAM_WEBHOOK_SECRET=test",
		"BOT_NAME=test",
		"CRON_SECRET=test",
		"SENTRY_DSN="+dsn.String(),
		"SENTRY_TRACES_SAMPLE_RATE=0",
		"SENTRY_ENABLE_LOGS=true",
	)
	output, err := cmd.CombinedOutput()
	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr, "subprocess output: %s", output)
	assert.Equal(t, 1, exitErr.ExitCode(), "subprocess output: %s", output)

	timeout := time.NewTimer(2 * time.Second)
	defer timeout.Stop()
	for {
		select {
		case envelope := <-envelopes:
			if bytes.Contains(envelope, []byte("Failed to initialize database")) {
				return
			}
		case <-timeout.C:
			t.Fatalf("Sentry did not receive the failure log before exit; subprocess output: %s", output)
		}
	}
}

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
