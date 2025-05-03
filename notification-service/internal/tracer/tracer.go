package tracer

import (
	"context"
	"log"
	"sync"

	"go.opentelemetry.io/otel"

	//"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func MustSetup(ctx context.Context, serviceName string) {
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("jaeger:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("failed to create OTLP trace exporter: %v", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	go func() {
		onceCloser := sync.OnceFunc(func() {
			log.Println("closing tracer")
			if err := tp.Shutdown(ctx); err != nil {
				log.Printf("error shutting down tracer provider: %v", err)
			}
		})

		<-ctx.Done()
		onceCloser()
	}()
}
