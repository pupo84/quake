package tracing

import (
	"context"
	"log"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// NewTraceProvider creates a new tracer provider with a jaeger exporter
func NewTraceProvider(ctx context.Context) (func(ctx context.Context) error, error) {
	environment := viper.GetString("ENVIRONMENT")
	serviceName := viper.GetString("SERVICE_NAME")
	openTelemetryUrl := viper.GetString("OPEN_TELEMETRY_COLLECTOR_URL")

	resource, err := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", environment),
		),
	)
	if err != nil {
		log.Printf("Failed to create resource: %s", err)
		return nil, err
	}

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(openTelemetryUrl)))
	if err != nil {
		log.Printf("Failed to create exporter: %s", err.Error())
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource),
	)

	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}
