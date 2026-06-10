package app

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// OpenTelemetry tracer provider 的关闭函数。
type ShutdownFunc func(context.Context) error

// 初始化 OpenTelemetry trace 导出和上下文传播器。
func InitTracer(ctx context.Context, cfg *Config, logger *slog.Logger) (ShutdownFunc, error) {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	if !cfg.Observ.Enabled {
		logger.Info("opentelemetry disabled")
		return func(context.Context) error { return nil }, nil
	}

	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(cfg.Observ.OTLPEndpoint))
	if err != nil {
		return nil, err
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(cfg.AppName), semconv.DeploymentEnvironmentName(cfg.Env)),
	)
	if err != nil {
		return nil, err
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.Observ.SampleRatio)),
	)
	otel.SetTracerProvider(provider)
	logger.Info("opentelemetry enabled", slog.String("endpoint", cfg.Observ.OTLPEndpoint), slog.Float64("sampleRatio", cfg.Observ.SampleRatio))
	return provider.Shutdown, nil
}
