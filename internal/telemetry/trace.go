package telemetry

import (
	"context"
	"fmt"

	"management-backend/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/zap"
)

// InitTracer 初始化 OpenTelemetry + Jaeger
func InitTracer(logger *zap.Logger) (func(context.Context) error, error) {
	if !config.C.Telemetry.Enabled {
		logger.Info("Telemetry 已禁用")
		return func(ctx context.Context) error { return nil }, nil
	}

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.C.Telemetry.Endpoint)))
	if err != nil {
		return nil, fmt.Errorf("Jaeger exporter 创建失败: %w", err)
	}

	rate := config.C.Telemetry.SampleRate
	if rate <= 0 {
		rate = 0.05 // 默认 5%
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(rate)),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("management-backend"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	logger.Info("Telemetry 初始化完成", zap.String("endpoint", config.C.Telemetry.Endpoint), zap.Float64("sampleRate", rate))
	return tp.Shutdown, nil
}
