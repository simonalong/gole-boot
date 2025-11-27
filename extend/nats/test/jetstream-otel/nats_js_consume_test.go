package nats

import (
	"fmt"
	"github.com/nats-io/nats.go/jetstream"
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/logger"
	"testing"
	"time"
)

// 使用环境变量：gole.profiles.active=consumer1
func TestNatsJsConsume(t *testing.T) {
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
		meta, _ := msg.Metadata()
		fmt.Println(string(msg.Data()), meta.Sequence)
		//assert.Equal(t, string(msg.Data()), "hello world")
		msg.Ack()
	})

	time.Sleep(12 * time.Hour)
}

func TestNatsJsConsume2(t *testing.T) {
	_, js, err := baseNats.GetJetStreamClient()
	if err != nil {
		logger.Fatal(err)
		return
	}

	consumer, err := baseNats.GetStreamConsumer(js, "stream-name1", "consumer2")
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

func TestNatsJsConsume3(t *testing.T) {
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

	for i := 0; i < 5; i++ {
		handler := func(consumeID int) jetstream.MessageHandler {
			return func(msg jetstream.Msg) {
				fmt.Printf("Received msg 【%v】on consume %d\n", string(msg.Data()), consumeID)
				msg.Ack()
			}
		}(i)

		_, err = consumer.Consume(handler)
	}
	time.Sleep(12 * time.Hour)
}
