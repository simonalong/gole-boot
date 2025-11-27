package original

// 工作模式

import (
	"github.com/simonalong/gole/util"
	"github.com/streadway/amqp"
	"log"
	"testing"
)

func TestWorkP(t *testing.T) {
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
		"work_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	for i := 0; i < 10; i++ {
		body := "Hello World" + "-" + util.ToString(i)
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(body),
			})
		if err != nil {
			log.Fatalf("Failed to publish a message: %s", err)
		}
		log.Printf(" [x] Sent %s", body)
		//time.Sleep(1 * time.Second)
	}
}

func TestWorkC1(t *testing.T) {
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
		"work_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

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
			log.Printf("C1 Received a message: %s\n", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func TestWorkC2(t *testing.T) {
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
		"work_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

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
			log.Printf("C2 Received a message: %s\n", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
