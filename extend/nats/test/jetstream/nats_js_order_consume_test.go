package nats

import (
	"fmt"
	"github.com/nats-io/nats.go/jetstream"
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/logger"
	"testing"
	"time"
)

// 使用环境变量：gole.profiles.active=order-consumer1
func TestNatsJsOrderConsume(t *testing.T) {
	_, js, err := baseNats.GetJetStreamClient()
	if err != nil {
		logger.Fatal(err)
		return
	}

	consumer, err := baseNats.GetStreamConsumer(js, "stream-name1", "order-consumer1")
	if err != nil {
		logger.Fatal(err)
		return
	}

	_, err = consumer.Consume(func(msg jetstream.Msg) {
		meta, _ := msg.Metadata()
		fmt.Println(string(msg.Data()), meta.Sequence)
		//assert.Equal(t, string(msg.Data()), "hello world")
		msg.Ack()
	})

	time.Sleep(12 * time.Hour)
}

// 使用环境变量：gole.profiles.active=order-consumer1
func TestNatsJsOrderConsume2(t *testing.T) {
	_, js, err := baseNats.GetJetStreamClient()
	if err != nil {
		logger.Fatal(err)
		return
	}

	consumer, err := baseNats.GetStreamConsumer(js, "stream-name1", "order-consumer1")
	if err != nil {
		logger.Fatal(err)
		return
	}

	_, err = consumer.Consume(func(msg jetstream.Msg) {
		meta, _ := msg.Metadata()
		fmt.Println(string(msg.Data()), meta.Sequence)
		//assert.Equal(t, string(msg.Data()), "hello world")
		msg.Ack()
	})

	time.Sleep(12 * time.Hour)
}
