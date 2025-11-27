package nats

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/simonalong/gole/logger"
	"go.opentelemetry.io/otel/trace"
)

type JetStreamClient struct {
	jetstream.JetStream
}

func (jsClient *JetStreamClient) Publish(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	// 前置
	pHookCtx, err := preHandle("js Publish", subject, nil, trace.SpanKindProducer)
	if err != nil {
		return nil, err
	}

	// 执行
	result, err := jsClient.JetStream.PublishMsg(ctx, &nats.Msg{
		Subject: subject,
		Header:  generateHeaderFromContext(pHookCtx.Context),
		Data:    payload,
	}, opts...)

	if err != nil {
		meterIncErrCounterValueOfNatsJs(subject, "Publish", "")
	}

	meterIncOkCounterValueOfNatsJs(subject, "Publish", "")
	//meterObserveByteValueOfNatsJs(subject, "Publish", "", len(payload))

	// 后置
	err = postHandle(pHookCtx, result, err)
	if err != nil {
		logger.Errorf("后置执行失败：%v", err.Error())
	}
	return result, err
}

func (jsClient *JetStreamClient) PublishMsg(ctx context.Context, msg *nats.Msg, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	// 前置
	pHookCtx, err := preHandle("js PublishMsg", msg.Subject, map[string]interface{}{"reply": msg.Reply, "header": msg.Header}, trace.SpanKindProducer)
	if err != nil {
		return nil, err
	}

	meterIncOkCounterValueOfNatsJs(msg.Subject, "PublishMsg", msg.Reply)
	//meterObserveByteValueOfNatsJs(msg.Subject, "PublishMsg", msg.Reply, len(msg.Data))

	// 执行
	msg.Header = appendHeaderFromContext(msg.Header, pHookCtx.Context)
	result, err := jsClient.JetStream.PublishMsg(ctx, msg, opts...)

	if err != nil {
		meterIncErrCounterValueOfNatsJs(msg.Subject, "PublishMsg", msg.Reply)
	}

	// 后置
	err = postHandle(pHookCtx, result, err)
	if err != nil {
		logger.Errorf("后置执行失败：%v", err.Error())
	}
	return result, err
}

func (jsClient *JetStreamClient) PublishAsync(subject string, payload []byte, opts ...jetstream.PublishOpt) (jetstream.PubAckFuture, error) {
	// 前置
	pHookCtx, err := preHandle("js PublishAsync", subject, nil, trace.SpanKindProducer)
	if err != nil {
		return nil, err
	}

	meterIncOkCounterValueOfNatsJs(subject, "PublishMsg", "")
	//meterObserveByteValueOfNatsJs(subject, "PublishMsg", "", len(payload))

	// 执行
	result, err := jsClient.JetStream.PublishMsgAsync(&nats.Msg{
		Subject: subject,
		Header:  generateHeaderFromContext(pHookCtx.Context),
		Data:    payload,
	}, opts...)

	if err != nil {
		meterIncErrCounterValueOfNatsJs(subject, "PublishMsg", "")
	}

	// 后置
	err = postHandle(pHookCtx, result, err)
	if err != nil {
		logger.Errorf("后置执行失败：%v", err.Error())
	}
	return result, err
}

func (jsClient *JetStreamClient) PublishMsgAsync(msg *nats.Msg, opts ...jetstream.PublishOpt) (jetstream.PubAckFuture, error) {
	// 前置
	pHookCtx, err := preHandle("js PublishMsgAsync", msg.Subject, map[string]interface{}{"reply": msg.Reply, "header": msg.Header}, trace.SpanKindProducer)
	if err != nil {
		return nil, err
	}

	meterIncOkCounterValueOfNatsJs(msg.Subject, "PublishMsgAsync", "")
	//meterObserveByteValueOfNatsJs(msg.Subject, "PublishMsgAsync", "", len(msg.Data))

	// 执行
	msg.Header = appendHeaderFromContext(msg.Header, pHookCtx.Context)
	result, err := jsClient.JetStream.PublishMsgAsync(msg, opts...)

	if err != nil {
		meterIncErrCounterValueOfNatsJs(msg.Subject, "PublishMsgAsync", msg.Reply)
	}

	// 后置
	err = postHandle(pHookCtx, result, err)
	if err != nil {
		logger.Errorf("后置执行失败：%v", err.Error())
	}
	return result, err
}
