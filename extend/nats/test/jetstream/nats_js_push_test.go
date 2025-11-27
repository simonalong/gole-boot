package nats

import (
	"context"
	"github.com/nats-io/nats.go"
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/logger"
	"strconv"
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

	for i := 0; i < 100; i++ {
		pMsg := &nats.Msg{
			Subject: "test.pub.req",
			Data:    []byte("hello world " + strconv.Itoa(i)),
		}
		_, err = js.PublishMsg(context.Background(), pMsg)

		time.Sleep(1 * time.Second)
	}
}

// gole.profiles.active=push
func TestNatsJsMsgPushBig(t *testing.T) {
	_, js, err := baseNats.GetJetStreamClient()
	if err != nil {
		logger.Fatal(err)
		return
	}

	for i := 0; i < 100000000; i++ {
		pMsg := &nats.Msg{
			Subject: "test.pub.req",
			Data:    []byte("hello world " + strconv.Itoa(i)),
		}
		_, err = js.PublishMsg(context.Background(), pMsg)

		time.Sleep(100 * time.Millisecond)
	}
}

// gole.profiles.active=push
func TestNatsJsMsgPushNoSleep(t *testing.T) {
	_, js, err := baseNats.GetJetStreamClient()
	if err != nil {
		logger.Fatal(err)
		return
	}

	for i := 0; i < 100000000; i++ {
		pMsg := &nats.Msg{
			Subject: "test.pub.req",
			Data:    []byte("hello world " + strconv.Itoa(i)),
		}
		_, err = js.PublishMsg(context.Background(), pMsg)
	}
}
