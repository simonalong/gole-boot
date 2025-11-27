package nats

import (
	"github.com/nats-io/nats.go/jetstream"
	"github.com/simonalong/gole/global"
	"github.com/simonalong/gole/logger"
	"go.opentelemetry.io/otel/trace"
)

type JetStreamConsumer struct {
	jetstream.Consumer
}

func (jsConsumer *JetStreamConsumer) Consume(handler jetstream.MessageHandler, opts ...jetstream.PullConsumeOpt) (jetstream.ConsumeContext, error) {
	return jsConsumer.Consumer.Consume(func(msg jetstream.Msg) {
		global.SetGlobalContext(generateContextFromHeader(msg.Headers()))

		pHookCtx, err := preHandle("js Consume", msg.Subject(), map[string]interface{}{}, trace.SpanKindConsumer)
		if err != nil {
			MeterIncErrCounterValueOfNatsJsServer(msg.Subject())
			return
		}

		handler(msg)
		MeterIncOkCounterValueOfNatsJsServer(msg.Subject())
		err = postHandle(pHookCtx, nil, err)
		if err != nil {
			MeterIncErrCounterValueOfNatsJsServer(msg.Subject())
			logger.Errorf("js Consume 后置执行失败：%v", err.Error())
		}
	}, opts...)
}

func (jsConsumer *JetStreamConsumer) Fetch(batch int, handler jetstream.MessageHandler, opts ...jetstream.FetchOpt) error {
	messageBatch, err := jsConsumer.Consumer.Fetch(batch, opts...)
	if err != nil {
		return err
	}

	for msg := range messageBatch.Messages() {
		global.SetGlobalContext(generateContextFromHeader(msg.Headers()))

		pHookCtx, err := preHandle("js Consume", msg.Subject(), map[string]interface{}{}, trace.SpanKindConsumer)
		if err != nil {
			MeterIncErrCounterValueOfNatsJsServer(msg.Subject())
			continue
		}

		handler(msg)
		MeterIncOkCounterValueOfNatsJsServer(msg.Subject())
		err = postHandle(pHookCtx, nil, err)
		if err != nil {
			MeterIncErrCounterValueOfNatsJsServer(msg.Subject())
			logger.Errorf("js Consume 后置执行失败：%v", err.Error())
		}
	}
	return nil
}

func (jsConsumer *JetStreamConsumer) FetchBytes(maxBytes int, handler jetstream.MessageHandler, opts ...jetstream.FetchOpt) error {
	messageBatch, err := jsConsumer.Consumer.FetchBytes(maxBytes, opts...)
	if err != nil {
		return err
	}

	for msg := range messageBatch.Messages() {
		global.SetGlobalContext(generateContextFromHeader(msg.Headers()))

		pHookCtx, err := preHandle("js Consume", msg.Subject(), map[string]interface{}{}, trace.SpanKindConsumer)
		if err != nil {
			MeterIncErrCounterValueOfNatsJsServer(msg.Subject())
			continue
		}

		handler(msg)
		MeterIncOkCounterValueOfNatsJsServer(msg.Subject())
		err = postHandle(pHookCtx, nil, err)
		if err != nil {
			MeterIncErrCounterValueOfNatsJsServer(msg.Subject())
			logger.Errorf("js Consume 后置执行失败：%v", err.Error())
		}
	}
	return nil
}

func (jsConsumer *JetStreamConsumer) FetchNoWait(batch int, handler jetstream.MessageHandler) error {
	messageBatch, err := jsConsumer.Consumer.FetchNoWait(batch)
	if err != nil {
		return err
	}

	for msg := range messageBatch.Messages() {
		global.SetGlobalContext(generateContextFromHeader(msg.Headers()))

		pHookCtx, err := preHandle("js Consume", msg.Subject(), map[string]interface{}{}, trace.SpanKindConsumer)
		if err != nil {
			MeterIncErrCounterValueOfNatsJsServer(msg.Subject())
			continue
		}

		handler(msg)
		MeterIncOkCounterValueOfNatsJsServer(msg.Subject())
		err = postHandle(pHookCtx, nil, err)
		if err != nil {
			MeterIncErrCounterValueOfNatsJsServer(msg.Subject())
			logger.Errorf("js Consume 后置执行失败：%v", err.Error())
		}
	}
	return nil
}
