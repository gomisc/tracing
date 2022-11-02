package tracing

import (
	"context"
	"runtime"
	"strings"

	"git.eth4.dev/golibs/errors"
	"git.eth4.dev/golibs/fields"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Trace struct {
	span trace.Span
	ctx  context.Context
}

func SetTrace(ctx context.Context, args ...any) *Trace {
	tracer := fromContext(ctx)
	name := getOperationName(2)

	if len(args) != 0 {
		name = args[0].(string)
	}

	next, span := tracer.Start(ctx, name)

	return &Trace{
		span: span,
		ctx:  next,
	}
}

func (t *Trace) Context() context.Context {
	return t.ctx
}

func (t *Trace) TraceID() string {
	return t.span.SpanContext().TraceID().String()
}

func (t *Trace) WithFields(flds ...fields.Field) *Trace {
	attrs := attributes{}

	for _, f := range flds {
		f.Value().Extract(f.Key(), &attrs)
	}

	t.span.SetAttributes(attrs...)

	return t
}

func (t *Trace) WithError(err error, args ...any) (*Trace, error) {
	if l := len(args); l != 0 {
		if msg, ok := args[0].(string); ok && l == 1 {
			err = errors.Ctx().Pos(2).Wrap(err, msg)
		} else {
			err = errors.Ctx().Pos(2).Wrapf(err, msg, args[1:]...)
		}
	}

	t.span.SetStatus(codes.Error, err.Error())

	return t, err
}

func (t *Trace) WithFormattedError(err error, args ...any) (*Trace, error) {
	err = errors.Formatted(err, args...)

	t.span.SetStatus(codes.Error, err.Error())

	return t, err
}

func (t *Trace) End() {
	t.span.End()
}

// Получает имя вызывающей операции трейсинга по стеку вызовов
// nolint
func getOperationName(depth int) (fname string) {
	pc, _, _, _ := runtime.Caller(depth) // nolint: gomnd, dogsled
	fname = runtime.FuncForPC(pc).Name()

	if oi := strings.LastIndex(fname, "/"); oi > 0 {
		return fname[oi+1:]
	}

	return fname
}
