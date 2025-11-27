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
func TestNatsJsMsgFetch(t *testing.T) {
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

	for {
		err := consumer.Fetch(100, func(msg jetstream.Msg) {
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(msg.Data()))
			msg.Ack()
		}, jetstream.FetchMaxWait(1*time.Second))
		if err != nil {
			fmt.Println("Error fetching messages: ", err)
		}
	}
}

// gole.profiles.active=consumer1
func TestNatsJsMsgFetch2(t *testing.T) {
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

	for {
		err := consumer.Fetch(100, func(msg jetstream.Msg) {
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(msg.Data()))
			msg.Ack()
		}, jetstream.FetchMaxWait(1*time.Second))
		if err != nil {
			fmt.Println("Error fetching messages: ", err)
		}
	}

}

// gole.profiles.active=consumer1
func TestNatsJsMsgFetch3(t *testing.T) {
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

	for {
		err := consumer.Fetch(100, func(msg jetstream.Msg) {
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(msg.Data()))
			msg.Ack()
		}, jetstream.FetchMaxWait(1*time.Second))
		if err != nil {
			fmt.Println("Error fetching messages: ", err)
		}
	}

}
