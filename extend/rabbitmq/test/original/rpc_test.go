package original

// 工作模式

import (
	"github.com/simonalong/gole/goid"
	"github.com/simonalong/gole/util"
	"github.com/streadway/amqp"
	"log"
	"testing"
	"time"
)

func TestRpcP(t *testing.T) {
	//conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	conn, err := amqp.Dial("amqp://admin:123456@192.168.1.75:5672/baseBoot")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_rsp_queue", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, //queue
		"",     //exchange
		false,  // auto-ack
		false,  //exclusive
		false,  //no-local
		false,  //no-wait
		nil,    //args
	)

	corrId := goid.GenerateUUID()

	go func() {
		for d := range msgs {
			if corrId == d.CorrelationId {
				res := string(d.Body)
				log.Print("收到rpc服务端的响应", res)
			}
		}
	}()

	for i := 0; i < 10; i++ {
		body := "Hello World_" + util.ToString(i)
		_ = ch.Publish(
			"",              // exchange
			"rpc_req_queue", // routing key
			false,           // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: corrId,
				ReplyTo:       q.Name,
				Body:          []byte(body),
			})
		time.Sleep(1 * time.Second)
	}

	log.Print("p1 发送数据成功")

	time.Sleep(10 * time.Second)
}

func TestRpcC1(t *testing.T) {
	//conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	conn, err := amqp.Dial("amqp://admin:123456@192.168.1.75:5672/baseBoot")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	// 声明队列
	q, err := ch.QueueDeclare(
		"rpc_req_queue", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, //global
	)

	// 消费数据
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("C1 收到请求端信息: %s", d.Body)

			err = ch.Publish(
				"",        //exchange
				d.ReplyTo, //routing key
				false,     //mandatory
				false,     //immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte("来自c1的响应：我收到了【" + string(d.Body) + "】"),
				})
		}
	}()

	<-forever
}

func TestRpcC2(t *testing.T) {
	//conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	conn, err := amqp.Dial("amqp://admin:123456@192.168.1.75:5672/baseBoot")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	// 声明队列
	q, err := ch.QueueDeclare(
		"rpc_req_queue", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, //global
	)

	// 消费数据
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("C1 收到请求端信息: %s", d.Body)

			err = ch.Publish(
				"",        //exchange
				d.ReplyTo, //routing key
				false,     //mandatory
				false,     //immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte("来自c2的响应：我收到了【" + string(d.Body) + "】"),
				})
		}
	}()

	<-forever
}
