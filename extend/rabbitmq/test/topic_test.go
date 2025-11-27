package test

import (
	"github.com/simonalong/gole-boot/extend/rabbitmq"
	"github.com/simonalong/gole/logger"
	"github.com/streadway/amqp"
	"testing"
	"time"
)

// gole.profiles.active=topic-p1
func TestTopicP1(t *testing.T) {
	mqClient, err := rabbitmq.GetClient()
	if err != nil {
		logger.Fatalf("获取mqClient失败:%v", err)
	}
	publisher := mqClient.GetProducer("p1")
	for i := 0; i < 10; i++ {
		body := keyChange(i) + ":" + "Hello World"
		_ = publisher.Send(keyChange(i), body)
		logger.Infof(" [x] Sent %s", body)
		time.Sleep(1 * time.Second)
	}
}

func keyChange(index int) string {
	if index < 3 {
		return "low.key"
	} else if index < 6 {
		return "middle.key"
	} else if index < 9 {
		return "high.key"
	} else if index < 12 {
		return "all.data.key"
	}
	return "all"
}

// gole.profiles.active=topic-c1
func TestTopicC1(t *testing.T) {
	rbtMq, err := rabbitmq.GetClient()
	if err != nil {
		logger.Fatalf("获取mqClient失败:%v", err)
	}

	consumer := rbtMq.GetConsumer("c1")
	consumer.Consume(func(d amqp.Delivery) {
		logger.Infof("Received a message: %s", d.Body)
	})
	time.Sleep(1000 * time.Second)
}

// gole.profiles.active=topic-c2
func TestTopicC2(t *testing.T) {
	rbtMq, err := rabbitmq.GetClient()
	if err != nil {
		logger.Fatalf("获取mqClient失败:%v", err)
	}

	consumer := rbtMq.GetConsumer("c2")
	consumer.Consume(func(d amqp.Delivery) {
		logger.Infof("Received a message: %s", d.Body)
	})
	time.Sleep(1000 * time.Second)
}
