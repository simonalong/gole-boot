package test

import (
	"github.com/simonalong/gole-boot/extend/rabbitmq"
	"github.com/simonalong/gole/logger"
	"github.com/streadway/amqp"
	"testing"
	"time"
)

// gole.profiles.active=route-p1
func TestRouteP1(t *testing.T) {
	mqClient, err := rabbitmq.GetClient()
	if err != nil {
		logger.Fatalf("获取mqClient失败:%v", err)
	}
	publisher := mqClient.GetProducer("p1")
	for i := 0; i < 10; i++ {
		body := "Hello World" + "-" + indexToType(i)
		_ = publisher.Send(indexToType(i), body)
		logger.Infof(" [x] Sent %s", body)
		time.Sleep(1 * time.Second)
	}
}
func indexToType(index int) string {
	if index < 3 {
		return "debug"
	} else if index < 6 {
		return "info"
	} else if index < 9 {
		return "warn"
	} else if index < 12 {
		return "error"
	}
	return "info"
}

// gole.profiles.active=route-c1
func TestRouteC1(t *testing.T) {
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

// gole.profiles.active=route-c2
func TestRouteC2(t *testing.T) {
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
