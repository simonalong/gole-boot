package rabbitmq

import (
	"fmt"
	"github.com/simonalong/gole/config"
)

var RbtMqCfg RbtMqConfig

type RbtMqConfig struct {
	User      string                `json:"user"`
	Password  string                `json:"password"`
	Host      string                `json:"host"`
	Port      uint                  `json:"port"`
	Vhost     string                `json:"vhost"`     // 虚拟主机
	Queues    []RbtMqQueueConfig    `json:"queues"`    // 队列
	Exchanges []RbtMqExchangeConfig `json:"exchanges"` // 交换机
	Binds     []RbtMqBindConfig     `json:"binds"`     // 绑定关系
	Producers []RbtMqProducerConfig `json:"producers"` // 生产者
	Consumers []RbtMqConsumerConfig `json:"consumers"` // 消费者
}

type RbtMqQueueConfig struct {
	Name       string `json:"name"`       // 队列名字
	Durable    bool   `json:"durable"`    // 是否持久化，默认true
	AutoDelete bool   `json:"autoDelete"` // 是否自动删除，默认false
	Exclusive  bool   `json:"exclusive"`  // 是否独占队列，默认false
	NoWait     bool   `json:"noWait"`     // 是否不等待，默认值false
}

type RbtMqExchangeConfig struct {
	Name       string `json:"name"`
	Kind       string `json:"kind"`       // 交换机类型，direct，topic，fanout；默认：direct
	Durable    bool   `json:"durable"`    // 是否持久化，默认true
	AutoDelete bool   `json:"autoDelete"` // 是否自动删除，默认false
	Internal   bool   `json:"internal"`   // 是否内部交换机（内部交换机不发布信息），默认false
	NoWait     bool   `json:"noWait"`     // 是否不等待，默认false
}

type RbtMqBindConfig struct {
	Name     string `json:"name"`
	Exchange string `json:"exchange"` // 交换机名称
	Key      string `json:"key"`      // queue绑定的key；默认空
	Queue    string `json:"queue"`    // 绑定的队列名
	NoWait   bool   `json:"noWait"`   // 是否不等待：默认false
}

type RbtMqProducerConfig struct {
	Name       string                `json:"name"`
	Exchange   string                `json:"exchange"`  // 交换机名称
	Mandatory  bool                  `json:"mandatory"` // 强制性，默认false
	Immediate  bool                  `json:"immediate"` // 立即，默认false
	Publishing RbtMqPublishingConfig `json:"publishing"`
}
type RbtMqPublishingConfig struct {
	ContentType     string `json:"contentType"`     // MIME content type
	ContentEncoding string `json:"contentEncoding"` // MIME content encoding
	DeliveryMode    uint8  `json:"deliveryMode"`    // Transient (0 or 1) or Persistent (2)
	Priority        uint8  `json:"priority"`        // 0 to 9
	ReplyTo         string `json:"replyTo"`         // address to to reply to (ex: RPC)
	Expiration      string `json:"expiration"`      // message expiration spec
	Type            string `json:"type"`            // message type name
}

type RbtMqConsumerConfig struct {
	Name      string `json:"name"`
	Queue     string `json:"queue"`     // 消费者绑定的队列名
	AutoAck   bool   `json:"autoAck"`   // 是否自动提交，默认true
	Exclusive bool   `json:"exclusive"` // 是否独占这个queue：默认false
	NoWait    bool   `json:"noWait"`    // 是否等待回复后再交付，默认false
}

func initDefaultConfig() {
	// 处理queue默认值
	for i := range RbtMqCfg.Queues {
		if config.GetValueString("gole.rabbitmq.queues["+fmt.Sprint(i)+"].durable") == "" {
			RbtMqCfg.Queues[i].Durable = true
		}
		if config.GetValueString("gole.rabbitmq.queues["+fmt.Sprint(i)+"].autoDelete") == "" {
			RbtMqCfg.Queues[i].AutoDelete = false
		}
		if config.GetValueString("gole.rabbitmq.queues["+fmt.Sprint(i)+"].exclusive") == "" {
			RbtMqCfg.Queues[i].Exclusive = false
		}
		if config.GetValueString("gole.rabbitmq.queues["+fmt.Sprint(i)+"].noWait") == "" {
			RbtMqCfg.Queues[i].NoWait = false
		}
	}

	// 处理exchange默认值
	for i := range RbtMqCfg.Exchanges {
		if config.GetValueString("gole.rabbitmq.exchanges["+fmt.Sprint(i)+"].durable") == "" {
			RbtMqCfg.Exchanges[i].Durable = true
		}
		if config.GetValueString("gole.rabbitmq.exchanges["+fmt.Sprint(i)+"].autoDelete") == "" {
			RbtMqCfg.Exchanges[i].AutoDelete = false
		}
		if config.GetValueString("gole.rabbitmq.exchanges["+fmt.Sprint(i)+"].internal") == "" {
			RbtMqCfg.Exchanges[i].Internal = false
		}
		if config.GetValueString("gole.rabbitmq.exchanges["+fmt.Sprint(i)+"].noWait") == "" {
			RbtMqCfg.Exchanges[i].NoWait = false
		}
	}

	// 处理bind默认值
	for i := range RbtMqCfg.Binds {
		if config.GetValueString("gole.rabbitmq.binds["+fmt.Sprint(i)+"].key") == "" {
			RbtMqCfg.Binds[i].Key = ""
		}
		if config.GetValueString("gole.rabbitmq.binds["+fmt.Sprint(i)+"].noWait") == "" {
			RbtMqCfg.Binds[i].NoWait = false
		}
	}

	// 处理consumer默认值
	for i := range RbtMqCfg.Consumers {
		if config.GetValueString("gole.rabbitmq.consumers["+fmt.Sprint(i)+"].autoAck") == "" {
			RbtMqCfg.Consumers[i].AutoAck = true
		}
		if config.GetValueString("gole.rabbitmq.consumers["+fmt.Sprint(i)+"].exclusive") == "" {
			RbtMqCfg.Consumers[i].Exclusive = false
		}
		if config.GetValueString("gole.rabbitmq.consumers["+fmt.Sprint(i)+"].noWait") == "" {
			RbtMqCfg.Consumers[i].NoWait = false
		}
	}
}

func getQueue(queueName string) *RbtMqQueueConfig {
	for _, queue := range RbtMqCfg.Queues {
		if queue.Name == queueName {
			return &queue
		}
	}
	return nil
}

func getExchange(exchangeName string) *RbtMqExchangeConfig {
	for _, exchange := range RbtMqCfg.Exchanges {
		if exchange.Name == exchangeName {
			return &exchange
		}
	}
	return nil
}

func getBind(bindName string) *RbtMqBindConfig {
	for _, bind := range RbtMqCfg.Binds {
		if bind.Name == bindName {
			return &bind
		}
	}
	return nil
}

func getProducer(publisherName string) *RbtMqProducerConfig {
	for _, producer := range RbtMqCfg.Producers {
		if producer.Name == publisherName {
			return &producer
		}
	}
	return nil
}

func getConsumer(consumerName string) *RbtMqConsumerConfig {
	for _, consumer := range RbtMqCfg.Consumers {
		if consumer.Name == consumerName {
			return &consumer
		}
	}
	return nil
}
