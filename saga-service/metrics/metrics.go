package metrics

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// ActivityType is a logical group for all Saga activities that we want to track.
type ActivityType string

const (
	ActivityBooking      ActivityType = "booking"
	ActivityTicket       ActivityType = "ticket"
	ActivityPayment      ActivityType = "payment"
	ActivityNotification ActivityType = "notification"
)

var (
	meter metric.Meter

	SagaStarted   metric.Int64Counter
	SagaSucceeded metric.Int64Counter
	SagaFailed    metric.Int64Counter

	activityStarted   map[ActivityType]metric.Int64Counter
	activitySucceeded map[ActivityType]metric.Int64Counter
	activityFailed    map[ActivityType]metric.Int64Counter

	ActivityLatencyMs metric.Float64Histogram

	HTTPCallLatencyMs metric.Float64Histogram

	provider *sdkmetric.MeterProvider
)

func Init(ctx context.Context, collectorEndpoint string, serviceName string) error {

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(collectorEndpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("create otlp exporter: %w", err)
	}

	reader := sdkmetric.NewPeriodicReader(exporter)
	res, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		),
	)
	provider = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(provider)

	meter = otel.Meter("saga‑service/metrics")

	if SagaStarted, err = meter.Int64Counter("saga.started", metric.WithDescription("Total saga executions started")); err != nil {
		return err
	}
	if SagaSucceeded, err = meter.Int64Counter("saga.succeeded", metric.WithDescription("Total saga executions succeeded")); err != nil {
		return err
	}
	if SagaFailed, err = meter.Int64Counter("saga.failed", metric.WithDescription("Total saga executions failed")); err != nil {
		return err
	}

	activityStarted = make(map[ActivityType]metric.Int64Counter)
	activitySucceeded = make(map[ActivityType]metric.Int64Counter)
	activityFailed = make(map[ActivityType]metric.Int64Counter)

	for _, a := range []ActivityType{ActivityBooking, ActivityTicket, ActivityPayment, ActivityNotification} {
		if activityStarted[a], err = meter.Int64Counter(
			fmt.Sprintf("saga.%s.started", a),
			metric.WithDescription(fmt.Sprintf("%s activities started", a)),
		); err != nil {
			return err
		}
		if activitySucceeded[a], err = meter.Int64Counter(
			fmt.Sprintf("saga.%s.succeeded", a),
			metric.WithDescription(fmt.Sprintf("%s activities succeeded", a)),
		); err != nil {
			return err
		}
		if activityFailed[a], err = meter.Int64Counter(
			fmt.Sprintf("saga.%s.failed", a),
			metric.WithDescription(fmt.Sprintf("%s activities failed", a)),
		); err != nil {
			return err
		}
	}

	if ActivityLatencyMs, err = meter.Float64Histogram("saga.activity.latency_ms",
		metric.WithDescription("Activity latency (ms)"),
		metric.WithUnit("ms")); err != nil {
		return err
	}

	log.Printf("[metrics] OpenTelemetry exporter configured for %s (service=%s)\n", collectorEndpoint, serviceName)
	return nil
}

// Shutdown flushes metric data and shuts down the provider. Call from main()'s GracefulStop.
func Shutdown(ctx context.Context) error {
	if provider == nil {
		return nil
	}
	return provider.Shutdown(ctx)
}

// Saga helpers --------------------------------------------------------------

func IncSagaStarted(ctx context.Context) {
	SagaStarted.Add(ctx, 1)
}

func IncSagaSucceeded(ctx context.Context) {
	SagaSucceeded.Add(ctx, 1)
}

func IncSagaFailed(ctx context.Context, err error) {
	SagaFailed.Add(ctx, 1, metric.WithAttributes(attribute.String("reason", err.Error())))
}

// Activity helpers ----------------------------------------------------------

func IncActivityStarted(ctx context.Context, a ActivityType) {
	activityStarted[a].Add(ctx, 1)
}

func IncActivitySucceeded(ctx context.Context, a ActivityType) {
	activitySucceeded[a].Add(ctx, 1)
}

func IncActivityFailed(ctx context.Context, a ActivityType, err error) {
	activityFailed[a].Add(ctx, 1, metric.WithAttributes(attribute.String("reason", err.Error())))
}

// RecordActivityLatency records a duration for the given activity.
func RecordActivityLatency(ctx context.Context, a ActivityType, d time.Duration) {
	ActivityLatencyMs.Record(ctx, float64(d.Milliseconds()), metric.WithAttributes(attribute.String("activity", string(a))))
}
