package nats

import (
	"context"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/simonalong/gole-boot/constants"
	_ "github.com/simonalong/gole-boot/otel"
	"github.com/simonalong/gole/bean"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/global"
	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/validate"
	"sync"
)

// CfgOfNats 配置部分使用的是 nats.Options
var CfgOfNats ConfigOfNats
var initLock sync.Mutex
var meterLoad = sync.OnceFunc(func() {
	// 支持测量指标
	initMeterOfNatsClient()
	initMeterOfNatsJsClient()
})

func init() {
	config.Load()

	LoadConfig()
}

func LoadConfig() {
	if config.Loaded && config.GetValueBoolDefault("gole.nats.enable", false) {
		err := config.GetValueObject("gole.nats", &CfgOfNats)
		if err != nil {
			logger.Warn("读取 nats 配置异常")
			return
		}
	}

	success, _, errMsg := validate.Check(CfgOfNats)
	if !success {
		logger.Errorf("nats配置异常：%v", errMsg)
	}
}

func GetClient() (*Client, error) {
	if client, ok := bean.GetBean(constants.BeanNameNats).(*Client); ok {
		return client, nil
	}
	initLock.Lock()
	defer initLock.Unlock()
	if client, ok := bean.GetBean(constants.BeanNameNats).(*Client); ok {
		return client, nil
	}
	client, err := New()
	if err != nil {
		return nil, err
	}
	bean.AddBean(constants.BeanNameNats, client)
	return client, nil
}

func GetJetStreamClient() (*Client, *JetStreamClient, error) {
	natsClient, _ := bean.GetBean(constants.BeanNameNats).(*Client)
	if natsJsClient, ok := bean.GetBean(constants.BeanNameNatsJetstream).(*JetStreamClient); ok {
		return natsClient, natsJsClient, nil
	}
	initLock.Lock()
	defer initLock.Unlock()
	natsClient, _ = bean.GetBean(constants.BeanNameNats).(*Client)
	if jsClient, ok := bean.GetBean(constants.BeanNameNatsJetstream).(*JetStreamClient); ok {
		return natsClient, jsClient, nil
	}
	natsClient, natsJsClient, err := NewJetStream()
	if err != nil {
		return nil, nil, err
	}
	bean.AddBean(constants.BeanNameNats, natsClient)
	bean.AddBean(constants.BeanNameNatsJetstream, natsJsClient)
	return natsClient, natsJsClient, nil
}

func New() (*Client, error) {
	if !config.GetValueBoolDefault("gole.nats.enable", false) {
		logger.Error("gole.nats.enable 配置为false")
		return nil, errors.New("gole.nats.enable 配置为false，无法链接 nats")
	}

	url := CfgOfNats.Url
	if url == "" {
		url = nats.DefaultURL
	}

	// 添加 认证
	var Options []nats.Option
	if accountOption := GetAccountOption(); accountOption != nil {
		Options = append(Options, accountOption)
	}
	// 添加 name
	Options = append(Options, GetNameOption())

	nc, err := nats.Connect(url, Options...)
	if err != nil {
		logger.Warnf("nats连接失败：%v", err.Error())
		return nil, err
	}

	jsType := "nats"
	if config.GetValueBoolDefault("gole.nats.jetstream.enable", false) {
		jsType = "jetstream"
	}

	// 支持opentelemetry
	AddHook(&OtelNtHook{
		JsType: jsType,
		Tracer: global.Tracer,
	})

	// 支持测量指标
	initMeterOfNatsClient()

	return &Client{
		Conn: nc,
	}, nil
}

func NewJetStream() (*Client, *JetStreamClient, error) {
	if !config.GetValueBoolDefault("gole.nats.enable", false) || !config.GetValueBoolDefault("gole.nats.jetstream.enable", false) {
		logger.Error("gole.nats.enable 或者 gole.nats.jetstream.enable 配置为false，则不启动nats的js模式，请开启")
		return nil, nil, errors.New("gole.nats.enable 或者 gole.nats.jetstream.enable 配置为false，则不启动nats的js模式，请开启")
	}

	url := CfgOfNats.Url
	if url == "" {
		url = nats.DefaultURL
	}

	// 添加 认证
	var Options []nats.Option
	if accountOption := GetAccountOption(); accountOption != nil {
		Options = append(Options, accountOption)
	}
	// 添加 name
	Options = append(Options, GetNameOption())

	nc, err := nats.Connect(url, Options...)
	if err != nil {
		logger.Warnf("nats连接失败：%v", err.Error())
		return nil, nil, err
	}

	jsType := "nats"
	if config.GetValueBoolDefault("gole.nats.jetstream.enable", false) {
		jsType = "jetstream"
	}

	// 支持opentelemetry
	AddHook(&OtelNtHook{
		JsType: jsType,
		Tracer: global.Tracer,
	})

	// 加载指标
	meterLoad()

	js, err := jetstream.New(nc)
	return &Client{Conn: nc}, &JetStreamClient{JetStream: js}, err
}

func GetStream(js jetstream.JetStream, streamName string) (jetstream.Stream, error) {
	return js.CreateOrUpdateStream(context.Background(), GetStreamConfig(streamName))
}

func GetStreamConsumer(js jetstream.JetStream, streamName, consumerName string) (*JetStreamConsumer, error) {
	stream, err := js.CreateOrUpdateStream(context.Background(), GetStreamConfig(streamName))
	if err != nil {
		logger.Fatalf("创建流异常：%v", err)
		return nil, err
	}

	baseConsumerConfig, err := GetBaseNatsJsConsumerConfig(consumerName)
	if err != nil {
		return nil, err
	}
	if baseConsumerConfig.Order {
		consumer, err := stream.OrderedConsumer(context.Background(), jetStreamOrderConsumerConvert(baseConsumerConfig))
		return &JetStreamConsumer{Consumer: consumer}, err
	} else {
		consumer, err := stream.CreateOrUpdateConsumer(context.Background(), jetStreamConsumerConvert(baseConsumerConfig))
		return &JetStreamConsumer{Consumer: consumer}, err
	}
}

func GetAccountOption() nats.Option {
	if CfgOfNats.UserName != "" && CfgOfNats.Password != "" {
		return nats.UserInfo(CfgOfNats.UserName, CfgOfNats.Password)
	} else if CfgOfNats.Token != "" {
		return nats.Token(CfgOfNats.Token)
	} else if CfgOfNats.NkSeedFile != "" {
		opt, err := nats.NkeyOptionFromSeed(CfgOfNats.NkSeedFile)
		if err != nil {
			logger.Errorf("使用Nkey连接nats失败：%v", err.Error())
			return nil
		}
		return opt
	} else if CfgOfNats.CredentialsFile != "" {
		return nats.UserCredentials(CfgOfNats.CredentialsFile)
	}
	return nil
}

func GetNameOption() nats.Option {
	// 添加name
	natsName := CfgOfNats.Name
	if CfgOfNats.Name == "" {
		if val := config.BaseCfg.Application.Name; val != "" {
			natsName = val
		} else {
			logger.Fatalf("gole.application.name 不可为空")
		}
	}
	return nats.Name(natsName)
}

func GetStreamConfig(streamName string) jetstream.StreamConfig {
	for _, stream := range CfgOfNats.Jetstream.Streams {
		if stream.Name != streamName {
			continue
		}

		return jetStreamStreamConvert(stream)
	}

	return jetstream.StreamConfig{}
}

func GetBaseNatsJsConsumerConfig(consumerName string) (ConfigOfJetstreamConsumer, error) {
	for _, consumer := range CfgOfNats.Jetstream.Consumers {
		if consumer.Name != consumerName {
			continue
		}
		return consumer, nil
	}
	return ConfigOfJetstreamConsumer{}, errors.New(fmt.Sprintf("消费者：%v没找到，请检查配置", consumerName))
}
