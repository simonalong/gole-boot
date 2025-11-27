package nats

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/simonalong/gole/global"
	baseTime "github.com/simonalong/gole/time"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type MsgOfNats struct {
	*nats.Msg
}

func (m *MsgOfNats) Respond(data []byte) error {
	tracer := global.Tracer
	var ctx context.Context
	var startTime time.Time
	if tracer != nil {
		ctx, _ = tracer.Start(global.GetGlobalContext(), "nats: Respond", trace.WithSpanKind(trace.SpanKindClient))
		startTime = time.Now()
	}

	// 执行
	err := m.Msg.RespondMsg(&nats.Msg{
		Subject: m.Reply,
		Header:  generateHeaderFromContext(global.GetGlobalContext()),
		Data:    data,
	})

	if tracer != nil {
		span := trace.SpanFromContext(ctx)
		defer span.End()

		var attrs []attribute.KeyValue
		attrs = append(attrs, attribute.Key("response.Subject").String(m.Reply))
		attrs = append(attrs, attribute.Key("response.data").String(string(data)))
		attrs = append(attrs, attribute.Key("execute.time").String(baseTime.ParseDurationForView(time.Now().Sub(startTime))))

		if err != nil {
			attrs = append(attrs, attribute.Key("nats.err").String(err.Error()))
			span.SetAttributes(attrs...)

			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetAttributes(attrs...)
		}
	}
	return err
}

func (m *MsgOfNats) RespondMsg(msg *nats.Msg) error {
	tracer := global.Tracer
	var ctx context.Context
	var startTime time.Time
	if tracer != nil {
		ctx, _ = tracer.Start(global.GetGlobalContext(), "nats: RespondMsg", trace.WithSpanKind(trace.SpanKindClient))
		startTime = time.Now()

		msg.Header = appendHeaderFromContext(msg.Header, ctx)
	}

	// 执行
	err := m.Msg.RespondMsg(msg)

	if tracer != nil {
		span := trace.SpanFromContext(ctx)
		defer span.End()

		var attrs []attribute.KeyValue
		attrs = append(attrs, attribute.Key("response.Subject").String(m.Reply))
		attrs = append(attrs, attribute.Key("response.data").String(string(msg.Data)))
		attrs = append(attrs, attribute.Key("execute.time").String(baseTime.ParseDurationForView(time.Now().Sub(startTime))))

		if err != nil {
			attrs = append(attrs, attribute.Key("nats.err").String(err.Error()))
			span.SetAttributes(attrs...)

			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetAttributes(attrs...)
		}
	}
	return err
}
