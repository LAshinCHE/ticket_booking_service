package metrics

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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

func InitMetrics() {
	meter := otel.GetMeterProvider().Meter("booking-service-metrics")

	var err error
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
