package original

// 简单模式

import (
	"log"
	"testing"

	"github.com/streadway/amqp"
)

func TestSimpleP(t *testing.T) {
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

	q, err := ch.QueueDeclare(
		"simple_queue", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	body := "Hello World!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(" [x] Sent %s\n", body)
}

func TestSimpleC(t *testing.T) {
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
		"simple_queue", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
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
			log.Printf("Received a message: %s\n", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
