package original

import (
	"github.com/simonalong/gole/util"
	"github.com/streadway/amqp"
	"log"
	"testing"
	"time"
)

const PubSub_EXCHANGE = "pubSub_exchange1"
const PubSub_QUEUE1 = "pubSub_queue1"
const PubSub_QUEUE2 = "pubSub_queue2"

// 发布订阅模式测试
// 这里需要注意的是，队列名不能相同

func TestPubSubP1(t *testing.T) {
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
		PubSub_EXCHANGE, // name
		"fanout",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	// 发送数据
	for i := 0; i < 10; i++ {
		body := "Hello World" + "-" + util.ToString(i)
		err = ch.Publish(
			PubSub_EXCHANGE, // exchange
			"",              // routing key
			false,           // mandatory
			false,           // immediate
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

func TestPubSubC1(t *testing.T) {
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
		PubSub_EXCHANGE, // name
		"fanout",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)

	// 声明队列
	q, err := ch.QueueDeclare(
		PubSub_QUEUE1, // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)

	// 绑定队列到交换机
	err = ch.QueueBind(
		q.Name,          // queue name
		"",              // routing key
		PubSub_EXCHANGE, // exchange
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

func TestPubSubC2(t *testing.T) {
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
		PubSub_EXCHANGE, // name
		"fanout",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)

	// 声明队列
	q, err := ch.QueueDeclare(
		PubSub_QUEUE1, // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)

	// 绑定队列到交换机
	err = ch.QueueBind(
		q.Name,          // queue name
		"",              // routing key
		PubSub_EXCHANGE, // exchange
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

func TestPubSubC3(t *testing.T) {
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
		PubSub_EXCHANGE, // name
		"fanout",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)

	// 声明队列
	q, err := ch.QueueDeclare(
		PubSub_QUEUE2, // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)

	// 绑定队列到交换机
	err = ch.QueueBind(
		q.Name,          // queue name
		"",              // routing key
		PubSub_EXCHANGE, // exchange
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

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
