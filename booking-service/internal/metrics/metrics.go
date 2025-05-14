package metrics

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

const (
	contactLabel = "contact"
	handlerLabel = "handler"
	codeLabel    = "code"
)

var (
	okRespTotal      metric.Int64Counter
	clientErrorTotal metric.Int64Counter
	serverErrorTotal metric.Int64Counter

	requestLatency metric.Float64Histogram
)

func InitMetrics() {
	ctx := context.Background()

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint("otel-collector:4317"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("failed to create OTLP exporter: %v", err)
	}

	reader := sdkmetric.NewPeriodicReader(exporter)
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	otel.SetMeterProvider(provider)

	meter := otel.GetMeterProvider().Meter("booking-service-metrics")

	if err != nil {
		log.Fatalf("failed to create metric: %v", err)
	}

	clientErrorTotal, err = meter.Int64Counter(
		"booking_client_error_total",
		metric.WithDescription("Total number of clients error"),
	)
	if err != nil {
		log.Fatalf("failed to create metric: %v", err)
	}

	okRespTotal, err = meter.Int64Counter(
		"booking_ok_response_total",
		metric.WithDescription("Total number of OK responses by handler"),
	)
	if err != nil {
		log.Fatalf("failed to create metric: %v", err)
	}

	serverErrorTotal, err = meter.Int64Counter(
		"booking_server_error_total",
		metric.WithDescription("Total number of server error responses by with response code"),
	)
	if err != nil {
		log.Fatalf("failed to create metric: %v", err)
	}

	requestLatency, err = meter.Float64Histogram(
		"booking_http_request_duration_millisecond",
		metric.WithDescription("Время выполнения HTTP-запроса в миллисекундах"),
	)
	if err != nil {
		log.Fatalf("failed to create booking_http_request_duration_millisecond: %v", err)
	}
}

func IncOkRespByHandler(ctx context.Context, handler string) {
	okRespTotal.Add(ctx, 1, metric.WithAttributes(
		attribute.String(handlerLabel, handler),
	))
}

func IncServerErrorByHandler(ctx context.Context, handler string, code int) {
	serverErrorTotal.Add(ctx, 1, metric.WithAttributes(
		attribute.String(handlerLabel, handler),
		attribute.String(codeLabel, fmt.Sprint(code)),
	))
}

func IncClientErrorByHandler(ctx context.Context, handler string, code int) {
	clientErrorTotal.Add(ctx, 1, metric.WithAttributes(
		attribute.String(handlerLabel, handler),
		attribute.String(codeLabel, fmt.Sprint(code)),
	))
}

func ObserveRequestDuration(ctx context.Context, handler string, durationMillisecond float64) {
	requestLatency.Record(ctx, durationMillisecond, metric.WithAttributes(handlerLabelAttr(handler)))
}

func handlerLabelAttr(value string) attribute.KeyValue {
	return attribute.String(handlerLabel, value)
}
