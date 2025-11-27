package nats

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/simonalong/gole-boot/constants"
	"github.com/simonalong/gole/global"
	"github.com/simonalong/gole/logger"
	baseTime "github.com/simonalong/gole/time"
	"github.com/simonalong/gole/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var Hooks []NtHook

type NtHook interface {
	Before(c *NtHookContext, spanKind trace.SpanKind) (*NtHookContext, error)
	After(c *NtHookContext) error
}

type NtHookContext struct {
	Context context.Context
	// 开始时间
	Start time.Time
	// 主题：如果有主题的话，则添加，没有则空
	Subject string
	// 函数名
	FuncName string
	// 请求参数
	Args map[string]interface{}
	// 结果，有些结果很少，则可以记录下
	Result interface{}

	// 执行耗时
	ExecuteTime time.Duration
	Err         error
}

func AddHook(hook NtHook) {
	Hooks = append(Hooks, hook)
}

// ---------------------

type OtelNtHook struct {
	// 使用js还是nats原生；nats 或者 jetstream
	JsType string
	Tracer trace.Tracer
}

func (otlHook *OtelNtHook) Before(hookContext *NtHookContext, spanKind trace.SpanKind) (*NtHookContext, error) {
	if otlHook.Tracer != nil {
		ctx, _ := otlHook.Tracer.Start(hookContext.Context, "nats: "+hookContext.FuncName, trace.WithSpanKind(spanKind))
		hookContext.Context = ctx
	}
	return hookContext, nil
}

func (otlHook *OtelNtHook) After(nhc *NtHookContext) error {
	if otlHook.Tracer == nil {
		return nil
	}
	span := trace.SpanFromContext(nhc.Context)
	defer span.End()

	var attrs []attribute.KeyValue
	attrs = append(attrs, attribute.Key("nats.Subject").String(nhc.Subject))
	attrs = append(attrs, attribute.Key("nats.func").String(nhc.FuncName))
	attrs = append(attrs, attribute.Key("nats.func.result").String(util.ToString(nhc.Result)))
	attrs = append(attrs, attribute.Key("start.time").String(baseTime.TimeToStringYmdHmsS(nhc.Start)))
	attrs = append(attrs, attribute.Key("execute.time").String(baseTime.ParseDurationForView(nhc.ExecuteTime)))
	attrs = append(attrs, attribute.Key("nats.func.args").String(util.ToJsonString(nhc.Args)))

	if nhc.Err != nil {
		attrs = append(attrs, attribute.Key("nats.err").String(nhc.Err.Error()))
		span.SetAttributes(attrs...)

		span.RecordError(nhc.Err)
		span.SetStatus(codes.Error, nhc.Err.Error())
		return nhc.Err
	} else {
		span.SetAttributes(attrs...)
	}
	return nil
}

func preHandle(funcName, subject string, args map[string]interface{}, spanKind trace.SpanKind) (*NtHookContext, error) {
	pHookCtx := &NtHookContext{
		Context:  global.GetGlobalContext(),
		Start:    time.Now(),
		Subject:  subject,
		FuncName: funcName,
		Args:     args,
	}

	for _, hook := range Hooks {
		hookContext, err := hook.Before(pHookCtx, spanKind)
		if err != nil {
			logger.Errorf("前置执行失败：%v", err.Error())
			return hookContext, err
		}
	}
	return pHookCtx, nil
}

func postHandle(hookContext *NtHookContext, result interface{}, err error) error {
	hookContext.Result = result
	hookContext.Err = err
	hookContext.ExecuteTime = time.Now().Sub(hookContext.Start)
	for _, hook := range Hooks {
		err := hook.After(hookContext)
		if err != nil {
			return err
		}
	}
	return nil
}

func postHandleWithoutExecuteTime(hookContext *NtHookContext, result interface{}, err error) error {
	hookContext.Result = result
	hookContext.Err = err
	for _, hook := range Hooks {
		err := hook.After(hookContext)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateHeaderFromContext(ctx context.Context) nats.Header {
	spanCtx := trace.SpanContextFromContext(ctx)

	header := nats.Header{}
	header.Set(constants.TraceId, spanCtx.TraceID().String())
	header.Set(constants.SpanId, spanCtx.SpanID().String())
	return header
}

func appendHeaderFromContext(srcHeader nats.Header, ctx context.Context) nats.Header {
	spanCtx := trace.SpanContextFromContext(ctx)

	header := srcHeader
	if srcHeader == nil {
		header = nats.Header{}
	}
	header.Set(constants.TraceId, spanCtx.TraceID().String())
	header.Set(constants.SpanId, spanCtx.SpanID().String())
	return header
}
