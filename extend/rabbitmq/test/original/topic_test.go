package original

import (
	"github.com/streadway/amqp"
	"log"
	"testing"
	"time"
)

const topic_EXCHANGE = "topic_exchange"
const topic_QUEUE1 = "topic_queue1_1"
const topic_QUEUE2 = "topic_queue2_1"

func TestTopicP1(t *testing.T) {
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

	// 声明交换机
	err = ch.ExchangeDeclare(
		topic_EXCHANGE, // name
		"topic",        // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	// 发送数据
	for i := 0; i < 10; i++ {
		body := keyChange(i) + ":" + "Hello World"
		err = ch.Publish(
			topic_EXCHANGE, // exchange
			keyChange(i),   // routing key
			false,          // mandatory
			false,          // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(body),
			})
		if err != nil {
			log.Fatalf("Failed to publish a message: %s", err)
		}
		log.Printf(" [x] Sent %s", body)
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

func TestTopicC1(t *testing.T) {
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

	// 声明交换机
	err = ch.ExchangeDeclare(
		topic_EXCHANGE, // name
		"topic",        // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)

	// 声明队列
	q, err := ch.QueueDeclare(
		topic_QUEUE1, // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)

	// 绑定队列到交换机
	err = ch.QueueBind(
		topic_QUEUE1,   // queue name
		"low.*",        // routing key
		topic_EXCHANGE, // exchange
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	// 创建消费者
	msgs, err := ch.Consume(
		q.Name, // 引用前面的队列名
		"",     // 消费者名字，不填自动生成一个
		true,   // 自动向队列确认消息已经处理
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// 循环消费队列中的消息
	for d := range msgs {
		log.Printf("接收消息=%s", d.Body)
	}
}

func TestTopicC2(t *testing.T) {
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

	// 声明交换机
	err = ch.ExchangeDeclare(
		topic_EXCHANGE, // name
		"topic",        // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)

	// 声明队列
	q, err := ch.QueueDeclare(
		topic_QUEUE2, // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	// 绑定队列到交换机
	err = ch.QueueBind(
		q.Name,         // queue name
		"#.key",        // routing key
		topic_EXCHANGE, // exchange
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	// 创建消费者
	msgs, err := ch.Consume(
		q.Name, // 引用前面的队列名
		"",     // 消费者名字，不填自动生成一个
		true,   // 自动向队列确认消息已经处理
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// 循环消费队列中的消息
	for d := range msgs {
		log.Printf("接收消息=%s", d.Body)
	}
}
