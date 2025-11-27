package nats

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/logger"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// 使用环境变量：gole.profiles.active=consumer1
func TestNatsConnectWithUser(t *testing.T) {
	_, js, err := baseNats.GetJetStreamClient()
	if err != nil {
		logger.Fatal(err)
		return
	}

	consumer, err := baseNats.GetStreamConsumer(js, "stream-name1", "consumer1")
	if err != nil {
		logger.Fatal(err)
		return
	}

	_, err = consumer.Consume(func(msg jetstream.Msg) {
		fmt.Println(string(msg.Data()))
		assert.Equal(t, string(msg.Data()), "hello world")
		msg.Ack()
	})

	pMsg := &nats.Msg{
		Subject: "test.pub.req",
		Data:    []byte("hello world"),
	}
	_, err = js.PublishMsg(context.Background(), pMsg)

	// 或者使用nats客户端进行推送都可以，效果是一样的
	//pMsg := &nats.Msg{
	//	Subject: "test.pub.req",
	//	Data:    []byte("hello world"),
	//}
	//err = nc.PublishMsg(pMsg)
	time.Sleep(5 * time.Second)
}
