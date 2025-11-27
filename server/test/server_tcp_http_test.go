package test

import (
	"github.com/gin-gonic/gin"
	httpServer "github.com/simonalong/gole-boot/server/http"
	"github.com/simonalong/gole-boot/server/http/rsp"
	"github.com/simonalong/gole-boot/server/tcp"
	"github.com/simonalong/gole/logger"
	"testing"
)

// gole.profiles.active=server-tcp-http
func TestServerTcpAndHttp(t *testing.T) {
	// tcp：设置编码解码器
	tcp.SetDecoder(func() tcp.Decoder { return &ServerMsgDemoSaveConCodec{} })
	tcp.Receive(func(msg interface{}) ([]byte, error) {
		// 自己的业务代码...
		logger.Info("给客户端返回消息：", "收到了")
		return nil, nil
	})

	// http：启动http服务
	httpServer.AddGinRoute("/api/data", httpServer.HmGet, func(c *gin.Context) {
		rsp.Done(c, "ok")
	})

	httpServer.RunServer()
}
