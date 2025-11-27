package nats

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/simonalong/gole-boot/constants"
	"github.com/simonalong/gole/global"
	"github.com/simonalong/gole/logger"
	baseTime "github.com/simonalong/gole/time"
	"github.com/simonalong/gole/util"
	"go.opentelemetry.io/otel/trace"
	"reflect"
	"time"
)

type Client struct {
	*nats.Conn
}

type MsgOfNatsHandler func(msg *MsgOfNats)

func (client *Client) AuthRequired() bool {
	// 前置
	pHookCtx, err := preHandle("AuthRequired", "", nil, trace.SpanKindProducer)
	if err != nil {
		return false
	}

	// 执行
	result := client.Conn.AuthRequired()

	// 后置
	err = postHandle(pHookCtx, result, nil)
	if err != nil {
		logger.Errorf("后置执行失败：%v", err.Error())
	}
	return result
}

func (client *Client) Barrier(f func()) error {
	// 前置
	pHookCtx, err := preHandle("Barrier", "", nil, trace.SpanKindProducer)
	if err != nil {
		return err
	}

	// 执行
	err = client.Conn.Barrier(f)

	// 后置
	err = postHandle(pHookCtx, nil, err)
	if err != nil {
		logger.Errorf("后置执行失败：%v", err.Error())
	}
	return err
}

func (client *Client) Buffered() (int, error) {
	// 前置
	pHookCtx, err := preHandle("Buffered", "", nil, trace.SpanKindProducer)
	if err != nil {
		return 0, err
	}

	// 执行
	result, err := client.Conn.Buffered()

	// 后置
	err = postHandle(pHookCtx, result, err)
	if err != nil {
		logger.Errorf("后置执行失败：%v", err.Error())
	}
	return result, err
}

func (client *Client) ChanQueueSubscribe(subj, group string, ch chan *nats.Msg) (*nats.Subscription, error) {
	// 前置
	pHookCtx, err := preHandle("ChanQueueSubscribe", subj, map[string]interface{}{"subj": subj, "group": group}, trace.SpanKindProducer)
	if err != nil {
		return nil, err
	}

	// 执行
	result, err := client.Conn.ChanQueueSubscribe(subj, group, ch)

	// 后置
	err = postHandle(pHookCtx, result, err)
	if err != nil {
		logger.Errorf("后置执行失败：%v", err.Error())
	}
	return result, err
}

func (client *Client) ChanSubscribe(subj string, ch chan *nats.Msg) (*nats.Subscription, error) {
	// 前置
	pHookCtx, err := preHandle("ChanSubscribe", subj, nil, trace.SpanKindProducer)
	if err != nil {
		return nil, err
	}

	// 执行
	result, err := client.Conn.ChanSubscribe(subj, ch)

	// 后置
	err = postHandle(pHookCtx, result, err)
	if err != nil {
		logger.Errorf("后置执行失败：%v", err.Error())
	}
	return result, err
}

func (client *Client) PublishEntity(subj string, data interface{}) error {
	var msg string
	if util.IsBaseType(reflect.TypeOf(data)) {
		msg = util.ToString(data)
	} else {
		msg = util.ToJsonString(data)
	}

	return client.Publish(subj, []byte(msg))
}

func (client *Client) Publish(subj string, data []byte) error {
	pHookCtx, err := preHandle("Publish", subj, nil, trace.SpanKindProducer)
	if err != nil {
		return err
	}

	meterIncOkCounterValueOfNats(subj, "Publish", "")
	//meterObserveByteValueOfNats(subj, "Publish", "", len(data))

	err = client.Conn.PublishMsg(&nats.Msg{
		Subject: subj,
		Header:  generateHeaderFromContext(pHookCtx.Context),
		Data:    data,
	})

	if err != nil {
		meterIncErrCounterValueOfNats(subj, "Publish", "")
	}

	err = postHandle(pHookCtx, nil, err)
	if err != nil {
		logger.Errorf("后置执行失败：%v", err.Error())
	}
	return err
}

func (client *Client) PublishMsg(m *nats.Msg) error {
	pHookCtx, err := preHandle("PublishMsg", m.Subject, map[string]interface{}{"reply": m.Reply, "header": m.Header}, trace.SpanKindProducer)
	if err != nil {
		return err
	}

	meterIncOkCounterValueOfNats(m.Subject, "PublishMsg", m.Reply)
	//meterObserveByteValueOfNats(m.Subject, "PublishMsg", m.Reply, len(m.Data))

	m.Header = appendHeaderFromContext(m.Header, pHookCtx.Context)
	err = client.Conn.PublishMsg(m)
	if err != nil {
		meterIncErrCounterValueOfNats(m.Subject, "PublishMsg", m.Reply)
	}

	err = postHandle(pHookCtx, nil, err)
	if err != nil {
		logger.Errorf("后置执行失败：%v", err.Error())
	}
	return err
}

func (client *Client) PublishRequest(subj, reply string, data []byte) error {
	pHookCtx, err := preHandle("PublishRequest", subj, map[string]interface{}{"reply": reply}, trace.SpanKindProducer)
	if err != nil {
		return err
	}

	meterIncOkCounterValueOfNats(subj, "PublishRequest", reply)
	//meterObserveByteValueOfNats(subj, "PublishRequest", reply, len(data))

	err = client.Conn.PublishMsg(&nats.Msg{
		Subject: subj,
		Reply:   reply,
		Header:  generateHeaderFromContext(pHookCtx.Context),
		Data:    data,
	})

	if err != nil {
		meterIncErrCounterValueOfNats(subj, "PublishRequest", reply)
	}

	err = postHandle(pHookCtx, nil, err)
	if err != nil {
		logger.Errorf("后置执行失败：%v", err.Error())
	}
	return err
}

func (client *Client) QueueSubscribe(subj, queue string, cb MsgOfNatsHandler) (*nats.Subscription, error) {
	return client.Conn.QueueSubscribe(subj, queue, func(m *nats.Msg) {
		global.SetGlobalContext(generateContextFromHeader(m.Header))

		pHookCtx, err := preHandle("QueueSubscribe", subj, map[string]interface{}{"queue": queue}, trace.SpanKindConsumer)
		if err != nil {
			return
		}

		MeterIncOkCounterValueOfNatsServer(subj)

		cb(&MsgOfNats{m})

		err = postHandle(pHookCtx, nil, err)
		if err != nil {
			MeterIncErrCounterValueOfNatsServer(subj)
			logger.Errorf("后置执行失败：%v", err.Error())
			return
		}
	})
}

func (client *Client) Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	pHookCtx, err := preHandle("Request", subj, map[string]interface{}{"timeout": baseTime.ParseDurationForView(timeout)}, trace.SpanKindProducer)
	if err != nil {
		return nil, err
	}

	meterIncOkCounterValueOfNats(subj, "Request", "")
	//meterObserveByteValueOfNats(subj, "Request", "", len(data))

	result, err := client.Conn.RequestMsg(&nats.Msg{
		Subject: subj,
		Header:  generateHeaderFromContext(pHookCtx.Context),
		Data:    data,
	}, timeout)

	if err != nil {
		meterIncErrCounterValueOfNats(subj, "Request", "")
		return nil, err
	}

	err = postHandle(pHookCtx, nil, err)
	if err != nil {
		logger.Errorf("后置执行失败：%v", err.Error())
		return nil, err
	}
	return result, err
}

func (client *Client) RequestMsg(msg *nats.Msg, timeout time.Duration) (*nats.Msg, error) {
	pHookCtx, err := preHandle("RequestMsg", msg.Subject, map[string]interface{}{"Reply": msg.Reply, "timeout": baseTime.ParseDurationForView(timeout)}, trace.SpanKindProducer)
	if err != nil {
		return nil, err
	}

	meterIncOkCounterValueOfNats(msg.Subject, "RequestMsg", msg.Reply)
	//meterObserveByteValueOfNats(msg.Subject, "RequestMsg", msg.Reply, len(msg.Data))

	msg.Header = appendHeaderFromContext(msg.Header, pHookCtx.Context)
	result, err := client.Conn.RequestMsg(msg, timeout)

	if err != nil {
		meterIncErrCounterValueOfNats(msg.Subject, "RequestMsg", msg.Reply)
	}

	err = postHandle(pHookCtx, nil, err)
	if err != nil {
		logger.Errorf("后置执行失败：%v", err.Error())
	}
	return result, err
}

func (client *Client) RequestWithContext(ctx context.Context, subj string, data []byte) (*nats.Msg, error) {
	pHookCtx, err := preHandle("RequestWithContext", subj, map[string]interface{}{"subj": subj}, trace.SpanKindProducer)
	if err != nil {
		return nil, err
	}

	meterIncOkCounterValueOfNats(subj, "RequestWithContext", "")
	//meterObserveByteValueOfNats(subj, "RequestWithContext", "", len(data))

	result, err := client.Conn.RequestMsgWithContext(ctx, &nats.Msg{
		Subject: subj,
		Header:  generateHeaderFromContext(pHookCtx.Context),
		Data:    data,
	})

	if err != nil {
		meterIncErrCounterValueOfNats(subj, "RequestWithContext", "")
	}

	err = postHandle(pHookCtx, nil, err)
	if err != nil {
		logger.Errorf("后置执行失败：%v", err.Error())
	}
	return result, err
}

func (client *Client) RequestMsgWithContext(ctx context.Context, msg *nats.Msg) (*nats.Msg, error) {
	pHookCtx, err := preHandle("RequestWithContext", msg.Subject, map[string]interface{}{"subj": msg.Subject}, trace.SpanKindProducer)
	if err != nil {
		return nil, err
	}

	meterIncOkCounterValueOfNats(msg.Subject, "RequestWithContext", msg.Reply)
	//meterObserveByteValueOfNats(msg.Subject, "RequestWithContext", msg.Reply, len(msg.Data))

	// 执行
	msg.Header = appendHeaderFromContext(msg.Header, pHookCtx.Context)
	result, err := client.Conn.RequestMsgWithContext(ctx, msg)

	if err != nil {
		meterIncErrCounterValueOfNats(msg.Subject, "RequestWithContext", msg.Reply)
	}

	// 后置
	err = postHandle(pHookCtx, nil, err)
	if err != nil {
		logger.Errorf("后置执行失败：%v", err.Error())
	}
	return result, err
}

func (client *Client) Subscribe(subj string, cb MsgOfNatsHandler) (*nats.Subscription, error) {
	return client.Conn.Subscribe(subj, func(m *nats.Msg) {
		global.SetGlobalContext(generateContextFromHeader(m.Header))

		pHookCtx, err := preHandle("Subscribe", subj, map[string]interface{}{}, trace.SpanKindConsumer)
		if err != nil {
			return
		}

		MeterIncOkCounterValueOfNatsServer(subj)
		cb(&MsgOfNats{m})

		err = postHandle(pHookCtx, nil, err)
		if err != nil {
			MeterIncErrCounterValueOfNatsServer(subj)
			logger.Errorf("后置执行失败：%v", err.Error())
		}
	})
}

func generateContextFromHeader(header nats.Header) context.Context {
	var traceId trace.TraceID
	var spanId trace.SpanID
	for key, values := range header {
		if key == constants.TraceId {
			_traceId, err := trace.TraceIDFromHex(values[0])
			if err == nil {
				traceId = _traceId
			}
		} else if key == constants.SpanId {
			_spanId, err := trace.SpanIDFromHex(values[0])
			if err == nil {
				spanId = _spanId
			}
		}
	}
	return trace.ContextWithRemoteSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceId,
		SpanID:  spanId,
	}))
}
