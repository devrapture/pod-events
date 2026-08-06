// Package metrics provides a thin, swappable wrapper around Sentry metrics.
// When metrics are disabled it falls back to a no-op recorder, so call sites
// can record metrics unconditionally.
package metrics

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/getsentry/sentry-go/attribute"
)

// Units supported by Sentry metrics.
const (
	UnitMillisecond = sentry.UnitMillisecond
	UnitSecond      = sentry.UnitSecond
)

// Option configures a single metric recording.
type Option struct {
	unit       string
	attributes []attribute.Builder
}

// WithUnit sets the unit for the metric.
func WithUnit(unit string) Option {
	return Option{unit: unit}
}

// WithAttributes attaches parameters to the metric.
func WithAttributes(attrs ...attribute.Builder) Option {
	return Option{attributes: attrs}
}

// Recorder records metrics. All methods are safe to call when metrics are
// disabled; the no-op recorder discards everything.
type Recorder interface {
	Count(name string, count int64, opts ...Option)
	Gauge(name string, value float64, opts ...Option)
	Distribution(name string, sample float64, opts ...Option)
}

// New returns a Sentry-backed recorder when enabled is true, otherwise a
// no-op recorder.
func New(ctx context.Context, enabled bool) Recorder {
	if !enabled {
		return noopRecorder{}
	}
	return sentryRecorder{meter: sentry.NewMeter(ctx)}
}

type sentryRecorder struct {
	meter sentry.Meter
}

func (r sentryRecorder) Count(name string, count int64, opts ...Option) {
	r.meter.Count(name, count, r.sentryOptions(opts)...)
}

func (r sentryRecorder) Gauge(name string, value float64, opts ...Option) {
	r.meter.Gauge(name, value, r.sentryOptions(opts)...)
}

func (r sentryRecorder) Distribution(name string, sample float64, opts ...Option) {
	r.meter.Distribution(name, sample, r.sentryOptions(opts)...)
}

func (r sentryRecorder) sentryOptions(opts []Option) []sentry.MeterOption {
	var meterOpts []sentry.MeterOption
	for _, opt := range opts {
		if opt.unit != "" {
			meterOpts = append(meterOpts, sentry.WithUnit(opt.unit))
		}
		if len(opt.attributes) > 0 {
			meterOpts = append(meterOpts, sentry.WithAttributes(opt.attributes...))
		}
	}
	return meterOpts
}

type noopRecorder struct{}

func (noopRecorder) Count(string, int64, ...Option)          {}
func (noopRecorder) Gauge(string, float64, ...Option)        {}
func (noopRecorder) Distribution(string, float64, ...Option) {}
