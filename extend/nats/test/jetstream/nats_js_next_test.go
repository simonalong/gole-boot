package nats

import (
	"fmt"
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/logger"
	"testing"
)

// 使用环境变量：gole.profiles.active=user
func TestNatsJsNext(t *testing.T) {
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

	// 目前这个每次启动还是会从0开始消费所有的数据
	for {
		msg, err := consumer.Next()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(string(msg.Data()))
		msg.Ack()
	}
}

func TestNatsJsNext2(t *testing.T) {
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
		msg, err := consumer.Next()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(string(msg.Data()))
		msg.Ack()
	}
}
