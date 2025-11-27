package original

import (
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/simonalong/gole/logger"
	"sync/atomic"
	"testing"
	"time"
)

type DemoXxxHandler struct {
	gnet.BuiltinEventEngine
	tester      *testing.T
	eng         gnet.Engine
	network     string
	addr        string
	multicore   bool
	packetBatch int
	started     int32
	aliveCount  int32
}

func (handler *DemoXxxHandler) OnBoot(eng gnet.Engine) (action gnet.Action) {
	fmt.Println("OnBoot")
	handler.eng = eng
	return
}

func (handler *DemoXxxHandler) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Println("OnOpen")
	c.SetContext(&DemoXxCodec{})
	atomic.AddInt32(&handler.aliveCount, 1)
	return
}

func (handler *DemoXxxHandler) OnShutdown(eng gnet.Engine) {
	fmt.Println("OnShutdown")
	atomic.SwapInt32(&handler.aliveCount, 0)
}

func (handler *DemoXxxHandler) OnClose(_ gnet.Conn, err error) (action gnet.Action) {
	fmt.Println("OnClose")
	if err != nil {
		logging.Debugf("error occurred on closed, %v\n", err)
	}
	atomic.AddInt32(&handler.aliveCount, -1)
	return
}

func (handler *DemoXxxHandler) OnTick() (delay time.Duration, action gnet.Action) {
	fmt.Println("OnTick")
	return 0, 0
}

func (handler *DemoXxxHandler) OnTraffic(c gnet.Conn) (action gnet.Action) {
	codec := c.Context().(*DemoXxCodec)
	data, err := codec.Decode(c)
	if err != nil {
		logger.Errorf("解码异常：%v", err)
		return gnet.None
	}

	logger.Infof("客户端在线个数：%d %d, 收到消息：%v", handler.aliveCount, handler.eng.CountConnections(), string(data))
	c.Write([]byte("ok"))
	return gnet.None
}
