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

func ProviderRunner(opts ...Option) execs.Runner {
	options := evaluateOptions(opts...)

	return execs.RunFunc(func(signals <-chan os.Signal, ready chan<- struct{}) error {
		// todo: validate options

		provider, err := jaegerTraceProvider(
			options.url,
			options.serviceName,
			options.envName,
		)
		if err != nil {
			return errors.Wrap(err, "prepare tracing provider")
		}

		otel.SetTracerProvider(provider)
		close(ready)
		<-signals

		ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
		defer cancel()

		if err = provider.Shutdown(ctx); err != nil {
			return errors.Wrap(err, "shutdown tracing provider")
		}

		return nil
	})
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
