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
	notifiedPositionsByContactTotal metric.Int64Counter
	okRespByHandlerTotal            metric.Int64Counter
	badRespByHandlerTotal           metric.Int64Counter
)

// func StartMetricsEndpoint() {
// 	go func() {
// 		http.Handle("/metrics", promhttp.Handler())
// 		if err := http.ListenAndServe(":2112", nil); err != nil {
// 			log.Fatalf("failed to start /metrics: %v", err)
// 		}
// 	}()
// }

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

	notifiedPositionsByContactTotal, err = meter.Int64Counter(
		"bakerbot_notified_positions_by_contact_total",
		metric.WithDescription("Total number of notified positions by contact"),
	)
	if err != nil {
		log.Fatalf("failed to create metric: %v", err)
	}

	okRespByHandlerTotal, err = meter.Int64Counter(
		"bakerbot_ok_response_by_handler_total",
		metric.WithDescription("Total number of OK responses by handler"),
	)
	if err != nil {
		log.Fatalf("failed to create metric: %v", err)
	}

	badRespByHandlerTotal, err = meter.Int64Counter(
		"bakerbot_bad_response_by_handler_total",
		metric.WithDescription("Total number of BAD responses by handler with response code"),
	)
	if err != nil {
		log.Fatalf("failed to create metric: %v", err)
	}
}

func AddNotifiedPositionsByContactTotal(ctx context.Context, cnt int64, contact string) {
	notifiedPositionsByContactTotal.Add(ctx, cnt, metric.WithAttributes(
		attribute.String(contactLabel, contact),
	))
}

func IncOkRespByHandler(ctx context.Context, handler string) {
	okRespByHandlerTotal.Add(ctx, 1, metric.WithAttributes(
		attribute.String(handlerLabel, handler),
	))
}

func IncBadRespByHandler(ctx context.Context, handler string, code int) {
	badRespByHandlerTotal.Add(ctx, 1, metric.WithAttributes(
		attribute.String(handlerLabel, handler),
		attribute.String(codeLabel, fmt.Sprint(code)),
	))
}
