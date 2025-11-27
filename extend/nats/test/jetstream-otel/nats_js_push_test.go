package nats

import (
	"context"
	"github.com/nats-io/nats.go"
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/logger"
	"testing"
	"time"
)

// gole.profiles.active=push
func TestNatsJsMsgPush(t *testing.T) {
	_, js, err := baseNats.GetJetStreamClient()
	if err != nil {
		logger.Fatal(err)
		return
	}

	_, err = js.Publish(context.Background(), "test.pub.req", []byte("nats jetstream hello world"))
	time.Sleep(1 * time.Second)
}

func TestNatsJsMsgPushMsg(t *testing.T) {
	_, js, err := baseNats.GetJetStreamClient()
	if err != nil {
		logger.Fatal(err)
		return
	}

	pMsg := &nats.Msg{
		Subject: "test.pub.req",
		Data:    []byte("nats jetstream hello world"),
	}
	_, err = js.PublishMsg(context.Background(), pMsg)

	time.Sleep(1 * time.Second)
}
