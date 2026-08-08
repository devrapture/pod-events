package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/getsentry/sentry-go/attribute"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoopRecorder(t *testing.T) {
	recorder := New(context.Background(), false)

	recorder.Count("test.count", 1)
	recorder.Gauge("test.gauge", 42)
	recorder.Distribution("test.distribution", 1.5, WithUnit(UnitMillisecond), WithAttributes(attribute.String("key", "value")))
}

func TestRecorderEmitsSentryMetrics(t *testing.T) {
	transport := &sentry.MockTransport{}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:       "https://public@example.com/1",
		Transport: transport,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		sentry.CurrentHub().Client().Close()
		sentry.CurrentHub().BindClient(nil)
	})

	recorder := New(context.Background(), true)
	recorder.Count("test.count", 5, WithAttributes(attribute.String("channel_type", "slack_webhook")))
	recorder.Gauge("test.gauge", 3.5)
	recorder.Distribution("test.duration", 187.5, WithUnit(UnitMillisecond))

	require.True(t, sentry.Flush(2*time.Second))

	events := transport.Events()
	require.Len(t, events, 1)
	require.Len(t, events[0].Metrics, 3)

	metricsByName := make(map[string]sentry.Metric)
	for _, metric := range events[0].Metrics {
		metricsByName[metric.Name] = metric
	}

	assert.Equal(t, sentry.MetricTypeCounter, metricsByName["test.count"].Type)
	assert.Equal(t, "slack_webhook", metricsByName["test.count"].Attributes["channel_type"].String())

	assert.Equal(t, sentry.MetricTypeGauge, metricsByName["test.gauge"].Type)

	duration := metricsByName["test.duration"]
	assert.Equal(t, sentry.MetricTypeDistribution, duration.Type)
	assert.Equal(t, UnitMillisecond, duration.Unit)
}
