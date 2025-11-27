package tcp

import (
	"context"
	"github.com/simonalong/gole/global"
	baseTime "github.com/simonalong/gole/time"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type Hook struct {
	Tracer trace.Tracer
	// 开始时间
	Start time.Time
}

func (hook *Hook) Before() context.Context {
	if hook.Tracer == nil {
		return context.Background()
	}

	ctx, _ := hook.Tracer.Start(global.GetGlobalContext(), "tcp: pre", trace.WithSpanKind(trace.SpanKindServer))
	hook.Start = time.Now()
	return ctx
}

func (hook *Hook) After(ctx context.Context, err error) {
	if hook.Tracer == nil {
		return
	}
	span := trace.SpanFromContext(ctx)
	defer span.End()

	var attrs []attribute.KeyValue
	attrs = append(attrs, attribute.Key("start.time").String(baseTime.TimeToStringYmdHmsS(hook.Start)))
	attrs = append(attrs, attribute.Key("execute.time").String(baseTime.ParseDurationForView(time.Now().Sub(hook.Start))))

	if err != nil {
		attrs = append(attrs, attribute.Key("tcp.err").String(err.Error()))
		span.SetAttributes(attrs...)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetAttributes(attrs...)
	}
}
