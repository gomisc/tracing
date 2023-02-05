package tracing

import (
	"context"
	"os"
	"time"

	"git.eth4.dev/golibs/errors"
	"git.eth4.dev/golibs/execs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

const ShutdownTimeout = time.Second * 5

type (
	traceProvider struct {
		execs.Runner
		tracingOptions
		prov *tracesdk.TracerProvider
	}
)

func ProviderRunner(ctx context.Context, opts ...Option) (context.Context, execs.Runner) {
	provider := &traceProvider{
		tracingOptions: evaluateOptions(opts...),
	}

	return context.WithValue(ctx, providerKey{}, provider), provider
}

func (tp *traceProvider) Run(signals <-chan os.Signal, ready chan<- struct{}) (err error) {
	if tp.prov, err = jaegerTraceProvider(tp.url, tp.serviceName, tp.envName); err != nil {
		return errors.Wrap(err, "prepare tracing provider")
	}

	otel.SetTracerProvider(tp.prov)
	close(ready)
	<-signals

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	if err = tp.prov.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutdown tracing provider")
	}

	return nil
}

// jaegerTraceProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter
func jaegerTraceProvider(url, service, env string) (*tracesdk.TracerProvider, error) {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, errors.Wrap(err, "create jaeger exporter")
	}

	provider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", env),
		)),
	)

	return provider, nil
}
