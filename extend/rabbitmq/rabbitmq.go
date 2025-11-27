package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/simonalong/gole-boot/constants"
	_ "github.com/simonalong/gole-boot/otel"
	"github.com/simonalong/gole/bean"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/goid"
	"github.com/simonalong/gole/logger"
	"github.com/streadway/amqp"
	"sync"
)

var initLock sync.Mutex

func init() {
	config.Load()

	if config.Loaded && config.GetValueBoolDefault("gole.rabbitmq.enable", false) {
		err := config.GetValueObject("gole.rabbitmq", &RbtMqCfg)
		// 配置默认值
		initDefaultConfig()
		if err != nil {
			logger.Warnf("读取rabbitmq配置异常：%v", err)
			return
		}
	}
}

type RbtMqClient struct {
	Chl *amqp.Channel
}

type RbtMqPublisher struct {
	*amqp.Channel
	ProducerConfig *RbtMqProducerConfig
}

type RbtMqConsumer struct {
	*amqp.Channel
	Msgs <-chan amqp.Delivery
}

func GetClient() (*RbtMqClient, error) {
	if bean.GetBean(constants.BeanNameRabbitmq) != nil {
		return bean.GetBean(constants.BeanNameRabbitmq).(*RbtMqClient), nil
	}
	initLock.Lock()
	defer initLock.Unlock()
	if bean.GetBean(constants.BeanNameRabbitmq) != nil {
		return bean.GetBean(constants.BeanNameRabbitmq).(*RbtMqClient), nil
	}
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	bean.AddBean(constants.BeanNameRabbitmq, client)
	return client, nil
}

func NewClient() (*RbtMqClient, error) {
	conn, err := amqp.Dial(generateUrl())
	if err != nil {
		logger.Errorf("连接rabbitmq失败:%v", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Errorf("获取Channel失败：%v", err)
		return nil, err
	}

	mqClient := &RbtMqClient{Chl: ch}

	// 加载exchange、queue、bind配置
	loadExchangeQueueBind(mqClient)

	return mqClient, nil
}

func generateUrl() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", RbtMqCfg.User, RbtMqCfg.Password, RbtMqCfg.Host, RbtMqCfg.Port, RbtMqCfg.Vhost)
}

func loadExchangeQueueBind(mqClient *RbtMqClient) {
	// 加载exchange配置
	for _, exchange := range RbtMqCfg.Exchanges {
		mqClient.UseExchange(exchange.Name)
	}

	// 加载queue配置
	for _, queue := range RbtMqCfg.Queues {
		mqClient.UseQueue(queue.Name)
	}

	// 加载bind配置
	for _, bind := range RbtMqCfg.Binds {
		mqClient.UseBind(bind.Name)
	}
}

func (r *RbtMqClient) Channel() *amqp.Channel {
	return r.Chl
}

func (r *RbtMqClient) UseQueue(queueName string) {
	r.UseQueueWithArgs(queueName, nil)
}

func (r *RbtMqClient) UseQueueWithArgs(queueName string, args map[string]interface{}) {
	rbtMqQueueCfg := getQueue(queueName)
	if rbtMqQueueCfg == nil {
		logger.Errorf("未找到队列配置，请排查配置文件是否配置:%s", queueName)
		return
	}

	_, err := r.QueueDeclare(
		rbtMqQueueCfg.Name,
		rbtMqQueueCfg.Durable,
		rbtMqQueueCfg.AutoDelete,
		rbtMqQueueCfg.Exclusive,
		rbtMqQueueCfg.NoWait,
		args,
	)
	if err != nil {
		logger.Errorf("声明队列失败:%v", err)
	}
}

func (r *RbtMqClient) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return r.Chl.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

func (r *RbtMqClient) UseBind(bindName string) {
	r.UseBindWithArgs(bindName, nil)
}

func (r *RbtMqClient) UseBindWithArgs(bindName string, args map[string]interface{}) {
	rbtMqBindCfg := getBind(bindName)
	if rbtMqBindCfg == nil {
		logger.Errorf("未找到绑定，请排查配置文件是否配置:%s", bindName)
		return
	}

	err := r.QueueBind(
		rbtMqBindCfg.Queue,
		rbtMqBindCfg.Key,
		rbtMqBindCfg.Exchange,
		rbtMqBindCfg.NoWait,
		args,
	)
	if err != nil {
		logger.Errorf("声明队列失败:%v", err)
	}
}

func (r *RbtMqClient) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	return r.Chl.QueueBind(name, key, exchange, noWait, args)
}

func (r *RbtMqClient) UseExchange(exchangeName string) {
	r.UseExchangeWithArgs(exchangeName, nil)
}

func (r *RbtMqClient) UseExchangeWithArgs(exchangeName string, args map[string]interface{}) {
	rbtMqExchangeCfg := getExchange(exchangeName)
	if rbtMqExchangeCfg == nil {
		logger.Errorf("未找到交换机配置，请排查配置文件是否配置exchanges:%s", exchangeName)
		return
	}

	err := r.ExchangeDeclare(rbtMqExchangeCfg.Name, rbtMqExchangeCfg.Kind, rbtMqExchangeCfg.Durable, rbtMqExchangeCfg.AutoDelete, rbtMqExchangeCfg.Internal, rbtMqExchangeCfg.NoWait, args)
	if err != nil {
		logger.Errorf("声明交换机失败:%v", err)
	}
}

func (r *RbtMqClient) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return r.Chl.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
}

func (r *RbtMqClient) GetProducer(publisherName string) *RbtMqPublisher {
	producerCfg := getProducer(publisherName)
	if producerCfg == nil {
		return r.DefaultPublisher()
	}
	return &RbtMqPublisher{
		Channel:        r.Chl,
		ProducerConfig: producerCfg,
	}
}

func (r *RbtMqClient) GetConsumer(consumer string) *RbtMqConsumer {
	consumerCfg := getConsumer(consumer)
	if consumerCfg == nil {
		return nil
	}
	msgs, err := r.Chl.Consume(consumerCfg.Queue, consumerCfg.Name, consumerCfg.AutoAck, consumerCfg.Exclusive, false, consumerCfg.NoWait, nil)
	if err != nil {
		logger.Errorf("声明队列失败:%v", err)
		return nil
	}
	return &RbtMqConsumer{
		Channel: r.Chl,
		Msgs:    msgs,
	}
}

func (r *RbtMqClient) DefaultPublisher() *RbtMqPublisher {
	return &RbtMqPublisher{
		Channel: r.Chl,
		ProducerConfig: &RbtMqProducerConfig{
			Exchange:  "",
			Mandatory: false,
			Immediate: false,
		},
	}
}

func (p *RbtMqPublisher) SimpleSend(queueName, msg string) error {
	return p.Send(queueName, msg)
}

func (p *RbtMqPublisher) SimpleSendJson(queueName, msg string) error {
	return p.SendJson(queueName, msg)
}

func (p *RbtMqPublisher) SimpleSendJsonData(queueName string, msg interface{}) error {
	return p.SendJsonData(queueName, msg)
}

func (p *RbtMqPublisher) Send(routeKey, msg string) error {
	return p.Channel.Publish(p.ProducerConfig.Exchange, routeKey, p.ProducerConfig.Mandatory, p.ProducerConfig.Immediate,
		amqp.Publishing{
			ContentType:     p.ProducerConfig.Publishing.ContentType,
			ContentEncoding: p.ProducerConfig.Publishing.ContentEncoding,
			DeliveryMode:    p.ProducerConfig.Publishing.DeliveryMode,
			Priority:        p.ProducerConfig.Publishing.Priority,
			ReplyTo:         p.ProducerConfig.Publishing.ReplyTo,
			Expiration:      p.ProducerConfig.Publishing.Expiration,
			Type:            p.ProducerConfig.Publishing.Type,
			Body:            []byte(msg),
		})
}

func (p *RbtMqPublisher) SendRpcReq(routeKey, replyToQueue, msg string) error {
	return p.Channel.Publish(p.ProducerConfig.Exchange, routeKey, p.ProducerConfig.Mandatory, p.ProducerConfig.Immediate,
		amqp.Publishing{
			ContentType:     p.ProducerConfig.Publishing.ContentType,
			ContentEncoding: p.ProducerConfig.Publishing.ContentEncoding,
			DeliveryMode:    p.ProducerConfig.Publishing.DeliveryMode,
			Priority:        p.ProducerConfig.Publishing.Priority,
			ReplyTo:         replyToQueue,
			Expiration:      p.ProducerConfig.Publishing.Expiration,
			CorrelationId:   goid.GenerateUUID(),
			Type:            p.ProducerConfig.Publishing.Type,
			Body:            []byte(msg),
		})
}

func (p *RbtMqPublisher) SendRpcRsp(d amqp.Delivery, msg string) error {
	return p.Channel.Publish(p.ProducerConfig.Exchange, d.ReplyTo, p.ProducerConfig.Mandatory, p.ProducerConfig.Immediate,
		amqp.Publishing{
			ContentType:     p.ProducerConfig.Publishing.ContentType,
			ContentEncoding: p.ProducerConfig.Publishing.ContentEncoding,
			DeliveryMode:    p.ProducerConfig.Publishing.DeliveryMode,
			Priority:        p.ProducerConfig.Publishing.Priority,
			ReplyTo:         p.ProducerConfig.Publishing.ReplyTo,
			Expiration:      p.ProducerConfig.Publishing.Expiration,
			CorrelationId:   d.CorrelationId,
			Type:            p.ProducerConfig.Publishing.Type,
			Body:            []byte(msg),
		})
}

func (p *RbtMqPublisher) SendPublishing(routeKey string, msg amqp.Publishing) error {
	return p.Channel.Publish(p.ProducerConfig.Exchange, routeKey, false, false, msg)
}

func (p *RbtMqPublisher) SendJson(routeKey, msg string) error {
	return p.Channel.Publish(p.ProducerConfig.Exchange, routeKey, false, false,
		amqp.Publishing{
			ContentType:     "application/json",
			ContentEncoding: p.ProducerConfig.Publishing.ContentEncoding,
			DeliveryMode:    p.ProducerConfig.Publishing.DeliveryMode,
			Priority:        p.ProducerConfig.Publishing.Priority,
			ReplyTo:         p.ProducerConfig.Publishing.ReplyTo,
			Expiration:      p.ProducerConfig.Publishing.Expiration,
			Type:            p.ProducerConfig.Publishing.Type,
			Body:            []byte(msg),
		})
}

func (p *RbtMqPublisher) SendJsonData(routeKey string, msg interface{}) error {
	marshalByte, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.Channel.Publish(p.ProducerConfig.Exchange, routeKey, false, false,
		amqp.Publishing{
			ContentType:     "application/json",
			ContentEncoding: p.ProducerConfig.Publishing.ContentEncoding,
			DeliveryMode:    p.ProducerConfig.Publishing.DeliveryMode,
			Priority:        p.ProducerConfig.Publishing.Priority,
			ReplyTo:         p.ProducerConfig.Publishing.ReplyTo,
			Expiration:      p.ProducerConfig.Publishing.Expiration,
			Type:            p.ProducerConfig.Publishing.Type,
			Body:            marshalByte,
		})
}

func (c *RbtMqConsumer) Consume(msg func(amqp.Delivery)) {
	goid.Go(func() {
		for d := range c.Msgs {
			msg(d)
		}
	})
}
