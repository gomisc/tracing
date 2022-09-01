package tracing

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type tracerKey struct{}

func WithContext(ctx context.Context, tracer trace.Tracer) context.Context {
	return context.WithValue(ctx, tracerKey{}, tracer)
}

func fromContext(ctx context.Context) trace.Tracer {
	if tracer, ok := ctx.Value(tracerKey{}).(trace.Tracer); ok {
		return tracer
	}

	return trace.NewNoopTracerProvider().Tracer("tracing.NoopTracer")
}
