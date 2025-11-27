package tcp

import (
	"github.com/DmitriyVTitov/size"
	"github.com/panjf2000/gnet/v2"
	"github.com/simonalong/gole/logger"
	"time"
)

type ReceiveHandler func(msg interface{}) ([]byte, error)
type ConnectCloseHandler func(c gnet.Conn) error
type ConnectOpenHandler func(c gnet.Conn)

type BaseReceiver struct {
	gnet.BuiltinEventEngine
	eng gnet.Engine
	// 中间处理器
	hook Hook
	// 消息处理器
	receiveHandler ReceiveHandler
	// 衡量活跃的时长
	activeJudgeDuration time.Duration
	// 上一次活跃时间
	lastActiveTime time.Time
}

func (baseHandler *BaseReceiver) OnBoot(eng gnet.Engine) (action gnet.Action) {
	baseHandler.eng = eng
	return
}

func (baseHandler *BaseReceiver) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	if rootCodecGenerate != nil {
		c.SetContext(rootCodecGenerate())
	}
	if rootConnectOpenHandler != nil {
		rootConnectOpenHandler(c)
	}
	logger.Group("tcp").Debugf("开启连接：%s:%s", c.RemoteAddr().Network(), c.RemoteAddr().String())
	defer func() {
		// 更新tcp连接数
		meterUpdateTcpConnectValue()
	}()
	return
}

func (baseHandler *BaseReceiver) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	defer func() {
		// 更新tcp连接数
		meterUpdateTcpConnectValue()
	}()

	if err != nil {
		action = gnet.Close
		if rootConnectCloseHandler != nil {
			err = rootConnectCloseHandler(c)
			if err != nil {
				action = gnet.Close
				logger.Warnf("关闭连接异常：%v", err)
				return action
			}
		}
		logger.Group("tcp").Debugf("关闭连接：%s:%s", c.RemoteAddr().Network(), c.RemoteAddr().String())
		return
	}
	if rootConnectCloseHandler != nil {
		err = rootConnectCloseHandler(c)
		if err != nil {
			action = gnet.Close
			logger.Warnf("关闭连接异常：%v", err)
			return action
		}
	}
	logger.Group("tcp").Debugf("关闭连接：%s:%s", c.RemoteAddr().Network(), c.RemoteAddr().String())
	return
}

func (baseHandler *BaseReceiver) OnTraffic(c gnet.Conn) (action gnet.Action) {
	baseHandler.lastActiveTime = time.Now()
	codec := c.Context().(Decoder)
	if codec == nil {
		return gnet.Close
	}

	ctx := baseHandler.hook.Before()

	// 统计指标
	meterIncTcpReqCounterValue()

	// 解码
	req, err := codec.Decode(c)
	if err != nil {
		logger.Errorf("解码异常：%v", err)
		baseHandler.hook.After(ctx, err)
		meterIncTcpReqErrCounterValue()
		return gnet.None
	}

	// 更新统计字节数
	meterObserveByteValue(size.Of(req.([]byte)))

	// 业务处理
	rspOriginal, err := baseHandler.receiveHandler(req)
	if err != nil {
		logger.Errorf("业务处理异常：%v", err)
		baseHandler.hook.After(ctx, err)
		meterIncTcpReqErrCounterValue()
		return gnet.None
	}
	if rspOriginal == nil {
		baseHandler.hook.After(ctx, nil)
		return gnet.None
	}
	_, err = c.Write(rspOriginal)
	if err != nil {
		logger.Errorf("rsp写入异常：%v", err)
		baseHandler.hook.After(ctx, err)
		meterIncTcpReqErrCounterValue()
		return gnet.None
	}
	baseHandler.hook.After(ctx, nil)
	return gnet.None
}
