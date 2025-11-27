package nats

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/logger"
	"strconv"
	"testing"
	"time"
)

// 使用环境变量：gole.profiles.active=cluster
func TestNatsJsCluster1(t *testing.T) {
	_, js, err := baseNats.GetJetStreamClient()
	if err != nil {
		logger.Fatal(err)
		return
	}

	consumer, err := baseNats.GetStreamConsumer(js, "stream-name2", "consumer1")
	if err != nil {
		logger.Fatal(err)
		return
	}

	_, err = consumer.Consume(func(msg jetstream.Msg) {
		meta, _ := msg.Metadata()
		fmt.Println(string(msg.Data()), meta.Sequence)
		msg.Ack()
	})

	time.Sleep(12 * time.Hour)
}

// 使用环境变量：gole.profiles.active=cluster
func TestNatsJsCluster2(t *testing.T) {
	_, js, err := baseNats.GetJetStreamClient()
	if err != nil {
		logger.Fatal(err)
		return
	}

	consumer, err := baseNats.GetStreamConsumer(js, "stream-name2", "consumer1")
	if err != nil {
		logger.Fatal(err)
		return
	}

	_, err = consumer.Consume(func(msg jetstream.Msg) {
		meta, _ := msg.Metadata()
		fmt.Println(string(msg.Data()), meta.Sequence)
		msg.Ack()
	})

	time.Sleep(12 * time.Hour)
}

// 使用环境变量：gole.profiles.active=cluster-one
func TestNatsJsClusterOne(t *testing.T) {
	_, js, err := baseNats.GetJetStreamClient()
	if err != nil {
		logger.Fatal(err)
		return
	}

	consumer, err := baseNats.GetStreamConsumer(js, "stream-name2", "consumer1")
	if err != nil {
		logger.Fatal(err)
		return
	}

	_, err = consumer.Consume(func(msg jetstream.Msg) {
		meta, _ := msg.Metadata()
		fmt.Println(string(msg.Data()), meta.Sequence)
		msg.Ack()
	})

	time.Sleep(12 * time.Hour)
}

// 使用环境变量：gole.profiles.active=cluster
func TestNatsJsClusterPush(t *testing.T) {
	_, js, err := baseNats.GetJetStreamClient()
	if err != nil {
		logger.Fatal(err)
		return
	}

	for i := 0; i < 100; i++ {
		pMsg := &nats.Msg{
			Subject: "cluster.pub.req",
			Data:    []byte("hello world " + strconv.Itoa(i)),
		}
		_, err = js.PublishMsg(context.Background(), pMsg)

		time.Sleep(1 * time.Second)
	}
}
