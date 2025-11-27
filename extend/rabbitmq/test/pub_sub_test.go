package test

import (
	"github.com/simonalong/gole-boot/extend/rabbitmq"
	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/util"
	"github.com/streadway/amqp"
	"testing"
	"time"
)

// gole.profiles.active=pubsub-p1
func TestPubSubP1(t *testing.T) {
	rbtMq, err := rabbitmq.GetClient()
	if err != nil {
		logger.Fatalf("获取mqClient失败:%v", err)
	}
	publisher := rbtMq.GetProducer("p1")

	for i := 0; i < 10; i++ {
		body := "Hello World" + "-" + util.ToString(i)
		_ = publisher.Send("", body)
		logger.Infof(" [x] Sent %s", body)
		time.Sleep(1 * time.Second)
	}
}

// gole.profiles.active=pubsub-c1
func TestPubSubC1(t *testing.T) {
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

// gole.profiles.active=pubsub-c2
func TestPubSubC2(t *testing.T) {
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

// gole.profiles.active=pubsub-c3
func TestPubSubC3(t *testing.T) {
	rbtMq, err := rabbitmq.GetClient()
	if err != nil {
		logger.Fatalf("获取mqClient失败:%v", err)
	}

	consumer := rbtMq.GetConsumer("c3")
	consumer.Consume(func(d amqp.Delivery) {
		logger.Infof("Received a message: %s", d.Body)
	})
	time.Sleep(1000 * time.Second)
}
