package tracing

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type (
	// tracer key
	tracerKey struct{}
	// provider key
	providerKey struct{}
)

func WithContext(ctx context.Context, tracer trace.Tracer) context.Context {
	return context.WithValue(ctx, tracerKey{}, tracer)
}

func fromContext(ctx context.Context) trace.Tracer {
	if tracer, ok := ctx.Value(tracerKey{}).(trace.Tracer); ok {
		return tracer
	}

	if provider, ok := ctx.Value(providerKey{}).(*traceProvider); ok {
		return provider.prov.Tracer(provider.serviceName)
	}

	return trace.NewNoopTracerProvider().Tracer("tracing.NoopTracer")
}
