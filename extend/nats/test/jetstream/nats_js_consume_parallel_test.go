package nats

import (
	"fmt"
	"github.com/nats-io/nats.go/jetstream"
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/logger"
	"testing"
	"time"
)

// gole.profiles.active=consumer1
func TestNatsJsConsumeParallel1(t *testing.T) {
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

// gole.profiles.active=consumer1
func TestNatsJsConsumeParallel2(t *testing.T) {
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
