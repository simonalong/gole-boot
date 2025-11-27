package test

import (
	"github.com/simonalong/gole-boot/extend/rabbitmq"
	"github.com/simonalong/gole/logger"
	"github.com/streadway/amqp"
	"testing"
	"time"
)

// gole.profiles.active=work-p1
func TestWorkP1(t *testing.T) {
	mqClient, err := rabbitmq.GetClient()
	if err != nil {
		logger.Fatalf("获取mqClient失败:%v", err)
	}
	publisher := mqClient.GetProducer("p1")
	_ = publisher.Send("simple_queue", "hello")
}

// gole.profiles.active=work-c1
func TestWorkC1(t *testing.T) {
	mqClient, err := rabbitmq.GetClient()
	if err != nil {
		logger.Fatalf("获取mqClient失败:%v", err)
	}
	consumer := mqClient.GetConsumer("c1")
	consumer.Consume(func(d amqp.Delivery) {
		logger.Infof("Received a message: %s", d.Body)
	})
	time.Sleep(1000 * time.Second)
}

// gole.profiles.active=work-c2
func TestWorkC2(t *testing.T) {
	mqClient, err := rabbitmq.GetClient()
	if err != nil {
		logger.Fatalf("获取mqClient失败:%v", err)
	}
	consumer := mqClient.GetConsumer("c2")
	consumer.Consume(func(d amqp.Delivery) {
		logger.Infof("Received a message: %s", d.Body)
	})
	time.Sleep(1000 * time.Second)
}
