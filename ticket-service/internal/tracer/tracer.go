package tracer

import (
	"context"
	"log"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"

	//"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func MustSetup(ctx context.Context, serviceName string) {
	// Создаём экспортёр Jaeger
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint("http://localhost:14268/api/traces"),
	))
	if err != nil {
		log.Fatalf("failed to create Jaeger exporter: %v", err)
	}

	// Создаём ресурс с именем сервиса
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	// Создаём TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	// Устанавливаем глобальный провайдер
	otel.SetTracerProvider(tp)

	// Обработка сигнала завершения
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
