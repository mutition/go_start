package tracing

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("default_tracer")

func InitJaegerProvider(jaegerURL string, serviceName string) (func(ctx context.Context) error, error) {
	if jaegerURL == "" {
		panic("jaegerURL is required")
	}
	tracer = otel.Tracer(serviceName)
	ctx := context.Background()
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(jaegerURL),
		otlptracehttp.WithInsecure(), // 如果使用 HTTP 而不是 HTTPS
	)
	if err != nil {
		return nil, err
	}
	//orivuder带上的data
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)
	b3propagator := b3.New(b3.WithInjectEncoding(b3.B3SingleHeader))
	p := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}, b3propagator)
	otel.SetTextMapPropagator(p)

	logrus.Infof("✅ TracerProvider initialized for service: %s", serviceName)
	logrus.Infof("✅ Jaeger endpoint: %s", jaegerURL)

	return tp.Shutdown, nil
}

func StartSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName)
}

func TraceID(ctx context.Context) string {
	return trace.SpanContextFromContext(ctx).TraceID().String()
}
