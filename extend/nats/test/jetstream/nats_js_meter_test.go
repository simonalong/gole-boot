package nats

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	httpServer "github.com/simonalong/gole-boot/server/http"
	"github.com/simonalong/gole/logger"
	"testing"
)

// gole.profiles.active=consumer-http
func TestServerHttpNatsJs(t *testing.T) {
	_, js, err := baseNats.GetJetStreamClient()
	if err != nil {
		logger.Fatal(err)
		return
	}

	consumer, err := baseNats.GetStreamConsumer(js, "stream-name2", "consumer2")
	if err != nil {
		logger.Fatal(err)
		return
	}

	_, err = consumer.Consume(func(msg jetstream.Msg) {
		meta, _ := msg.Metadata()
		baseNats.MeterIncOkCounterValueOfNatsJsServer(msg.Subject())
		fmt.Println(string(msg.Data()), meta.Sequence)
		//assert.Equal(t, string(msg.Data()), "hello world")
		msg.Ack()
	})

	httpServer.RunServer()
}

// gole.profiles.active=push-http
func TestNatsClientHttpJs(t *testing.T) {
	_, js, err := baseNats.GetJetStreamClient()
	if err != nil {
		logger.Fatal(err)
		return
	}

	httpServer.Get("/nats/js/send", func(c *gin.Context) (any, error) {
		pMsg := &nats.Msg{
			Subject: "test.js.demo.req",
			Data:    []byte("hello world "),
		}
		_, err = js.PublishMsg(context.Background(), pMsg)

		return "", nil
	})
	httpServer.RunServer()
}
