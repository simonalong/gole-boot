package nats

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	httpServer "github.com/simonalong/gole-boot/server/http"
	"github.com/simonalong/gole/logger"
)

// gole.profiles.active=user-http
func TestServerHttpNats(t *testing.T) {
	nc, err := baseNats.New()
	if err != nil {
		logger.Fatal(err)
		return
	}
	defer nc.Close()

	_, err = nc.Subscribe("test.sub.req", func(msg *baseNats.MsgOfNats) {
		logger.Info("接收到信息：", string(msg.Data))
	})

	httpServer.RunServer()
}

// gole.profiles.active=user-http-push
func TestNatsClientHttp(t *testing.T) {
	nc, err := baseNats.New()
	if err != nil {
		logger.Fatal(err)
		return
	}
	defer nc.Close()

	httpServer.Get("/nats/send", func(c *gin.Context) (any, error) {
		pMsg := &nats.Msg{
			Subject: "test.sub.req",
			Data:    []byte("nats hello world "),
		}
		_ = nc.PublishMsg(pMsg)

		return "ok", nil
	})
	httpServer.RunServer()
}
