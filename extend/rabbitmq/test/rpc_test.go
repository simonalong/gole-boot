package test

import (
	"github.com/simonalong/gole-boot/extend/rabbitmq"
	"github.com/simonalong/gole/logger"
	"github.com/streadway/amqp"
	"gorm.io/gorm/utils"
	"testing"
	"time"
)

// gole.profiles.active=rpc-p1
func TestRpcP1(t *testing.T) {
	mqClient, err := rabbitmq.GetClient()
	if err != nil {
		logger.Fatalf("获取mqClient失败:%v", err)
	}

	consumer := mqClient.GetConsumer("p1_c")
	consumer.Consume(func(d amqp.Delivery) {
		logger.Infof("收到消息: %s", d.Body)
	})

	publisher := mqClient.GetProducer("p1")

	for i := 0; i < 10; i++ {
		body := "Hello World" + "-" + utils.ToString(i)
		_ = publisher.SendRpcReq("rpc_req_queue", "rpc_rsp_queue", body)
		time.Sleep(1 * time.Second)
	}
}

// gole.profiles.active=rpc-c1
func TestRpcC1(t *testing.T) {
	mqClient, err := rabbitmq.GetClient()
	if err != nil {
		logger.Fatalf("获取mqClient失败:%v", err)
	}
	consumer := mqClient.GetConsumer("c1")
	consumer.Consume(func(d amqp.Delivery) {
		logger.Infof("收到消息: %s", d.Body)
		publisher := mqClient.GetProducer("c1_p")

		_ = publisher.SendRpcRsp(d, "来自c1：我收到了")
	})
	time.Sleep(1000 * time.Second)
}

// gole.profiles.active=rpc-c2
func TestRpcC2(t *testing.T) {
	mqClient, err := rabbitmq.GetClient()
	if err != nil {
		logger.Fatalf("获取mqClient失败:%v", err)
	}
	consumer := mqClient.GetConsumer("c2")
	consumer.Consume(func(d amqp.Delivery) {
		logger.Infof("收到消息: %s", d.Body)
		publisher := mqClient.GetProducer("c2_p")

		_ = publisher.SendRpcRsp(d, "来自c2：我收到了")
	})
	time.Sleep(1000 * time.Second)
}
