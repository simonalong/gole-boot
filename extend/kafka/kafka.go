package kafka

import (
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	_ "github.com/simonalong/gole-boot/otel"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/logger"
	"regexp"
	"strings"
	"time"
)

func init() {
	config.Load()

	if config.Loaded && config.GetValueBoolDefault("gole.kafka.enable", false) {
		err := config.GetValueObject("gole.kafka", &Cfg)
		if err != nil {
			logger.Warn("读取kafka配置异常")
			return
		}
	}
}

func NewClient() (sarama.Client, error) {
	cfg := GetKafkaConfig()
	if cfg == nil {
		logger.Warn("gole.kafka.enable为false，没有激活")
		return nil, errors.New("gole.kafka.enable为false，没有激活")
	}
	return sarama.NewClient(config.GetValueArrayString("gole.kafka.addrs"), cfg)
}

func NewAsyncProducer() (sarama.AsyncProducer, error) {
	cfg := GetKafkaConfig()
	if cfg == nil {
		logger.Warn("gole.kafka.enable为false，没有激活")
		return nil, errors.New("gole.kafka.enable为false，没有激活")
	}
	return sarama.NewAsyncProducer(config.GetValueArrayString("gole.kafka.addrs"), cfg)
}

func NewSyncProducer() (sarama.SyncProducer, error) {
	cfg := GetKafkaConfig()
	if cfg == nil {
		logger.Warn("gole.kafka.enable为false，没有激活")
		return nil, errors.New("gole.kafka.enable为false，没有激活")
	}
	return sarama.NewSyncProducer(config.GetValueArrayString("gole.kafka.addrs"), cfg)
}

func NewConsumer() (sarama.Consumer, error) {
	cfg := GetKafkaConfig()
	if cfg == nil {
		logger.Warn("gole.kafka.enable为false，没有激活")
		return nil, errors.New("gole.kafka.enable为false，没有激活")
	}
	return sarama.NewConsumer(config.GetValueArrayString("gole.kafka.addrs"), cfg)
}

func NewClusterAdmin() (sarama.ClusterAdmin, error) {
	cfg := GetKafkaConfig()
	if cfg == nil {
		logger.Warn("gole.kafka.enable为false，没有激活")
		return nil, errors.New("gole.kafka.enable为false，没有激活")
	}
	return sarama.NewClusterAdmin(config.GetValueArrayString("gole.kafka.addrs"), cfg)
}

func NewConsumerGroup(groupId string) (sarama.ConsumerGroup, error) {
	cfg := GetKafkaConfig()
	if cfg == nil {
		logger.Warn("gole.kafka.enable为false，没有激活")
		return nil, errors.New("gole.kafka.enable为false，没有激活")
	}
	return sarama.NewConsumerGroup(config.GetValueArrayString("gole.kafka.addrs"), groupId, cfg)
}

func GetKafkaConfig() *sarama.Config {
	if !config.GetValueBoolDefault("gole.kafka.enable", false) {
		return nil
	}

	kafkaConfig := sarama.NewConfig()
	if config.GetValueStringDefault("gole.kafka.client-id", "sarama") != "sarama" {
		kafkaConfig.ClientID = Cfg.ClientId
	}

	if config.GetValueIntDefault("gole.kafka.channel-buffer-size", 256) != 256 {
		kafkaConfig.ChannelBufferSize = Cfg.ChannelBufferSize
	}

	if config.GetValueBoolDefault("gole.kafka.api-versions-request", true) != true {
		kafkaConfig.ApiVersionsRequest = Cfg.ApiVersionsRequest
	}

	if config.GetValueStringDefault("gole.kafka.version", "V1_0_0_0") != "V1_0_0_0" {
		kafkaConfig.Version = getKafkaVersion(Cfg.Version)
	}

	//============================= admin =============================
	if config.GetValueIntDefault("gole.kafka.admin.retry-max", 5) != 5 {
		kafkaConfig.Admin.Retry.Max = Cfg.Admin.RetryMax
	}

	if config.GetValueStringDefault("gole.kafka.admin.retry-backoff", "100ms") != "100ms" {
		t, err := time.ParseDuration(Cfg.Admin.RetryBackoff)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.admin.retry-backoff】异常，%v", err.Error())
		} else {
			kafkaConfig.Admin.Retry.Backoff = t
		}
	}

	if config.GetValueStringDefault("gole.kafka.admin.timeout", "3s") != "3s" {
		t, err := time.ParseDuration(Cfg.Admin.Timeout)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.admin.timeout】异常，%v", err.Error())
		} else {
			kafkaConfig.Admin.Timeout = t
		}
	}

	//============================= net =============================
	if config.GetValueIntDefault("gole.kafka.net.max-open-requests", 5) != 5 {
		kafkaConfig.Net.MaxOpenRequests = Cfg.Net.MaxOpenRequests
	}

	if config.GetValueStringDefault("gole.kafka.net.dial-timeout", "3s") != "3s" {
		t, err := time.ParseDuration(Cfg.Net.DialTimeout)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.net.dial-timeout】异常，%v", err.Error())
		} else {
			kafkaConfig.Net.DialTimeout = t
		}
	}

	if config.GetValueStringDefault("gole.kafka.net.read-timeout", "3s") != "3s" {
		t, err := time.ParseDuration(Cfg.Net.ReadTimeout)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.net.read-timeout】异常，%v", err.Error())
		} else {
			kafkaConfig.Net.ReadTimeout = t
		}
	}

	if config.GetValueStringDefault("gole.kafka.net.write-timeout", "3s") != "3s" {
		t, err := time.ParseDuration(Cfg.Net.WriteTimeout)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.net.write-timeout】异常，%v", err.Error())
		} else {
			kafkaConfig.Net.WriteTimeout = t
		}
	}

	if config.GetValueBoolDefault("gole.kafka.net.SASL-handshake", true) != true {
		kafkaConfig.Net.SASL.Handshake = false
	}

	if config.GetValueIntDefault("gole.kafka.net.SASL-version", 0) != 0 {
		kafkaConfig.Net.SASL.Version = Cfg.Net.SASLVersion
	}

	//============================= metadata =============================
	if config.GetValueIntDefault("gole.kafka.metadata.retry-max", 3) != 3 {
		kafkaConfig.Metadata.Retry.Max = Cfg.Metadata.RetryMax
	}

	if config.GetValueStringDefault("gole.kafka.metadata.retry-backoff", "250ms") != "250ms" {
		t, err := time.ParseDuration(Cfg.Metadata.RetryBackoff)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.metadata.retry-backoff】异常，%v", err.Error())
		} else {
			kafkaConfig.Metadata.Retry.Backoff = t
		}
	}

	if config.GetValueStringDefault("gole.kafka.metadata.refresh-frequency", "10m") != "10m" {
		t, err := time.ParseDuration(Cfg.Metadata.RefreshFrequency)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.metadata.refresh-frequency】异常，%v", err.Error())
		} else {
			kafkaConfig.Metadata.RefreshFrequency = t
		}
	}

	if config.GetValueBoolDefault("gole.kafka.metadata.full", true) != true {
		kafkaConfig.Metadata.Full = Cfg.Metadata.Full
	}

	if config.GetValueBoolDefault("gole.kafka.metadata.allow-auto-topic-creation", true) != true {
		kafkaConfig.Metadata.AllowAutoTopicCreation = Cfg.Metadata.AllowAutoTopicCreation
	}

	//============================= producer =============================
	if config.GetValueIntDefault("gole.kafka.producer.max-message-bytes", 1000000) != 1000000 {
		kafkaConfig.Producer.MaxMessageBytes = Cfg.Producer.MaxMessageBytes
	}

	if config.GetValueIntDefault("gole.kafka.producer.required-acks", 1) != 1 {
		kafkaConfig.Producer.RequiredAcks = Cfg.Producer.RequiredAcks
	}

	if config.GetValueStringDefault("gole.kafka.producer.timeout", "10s") != "10s" {
		t, err := time.ParseDuration(Cfg.Producer.Timeout)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.producer.timeout】异常，%v", err.Error())
		} else {
			kafkaConfig.Producer.Timeout = t
		}
	}

	if config.GetValueIntDefault("gole.kafka.producer.retry-max", 3) != 3 {
		kafkaConfig.Producer.Retry.Max = Cfg.Producer.RetryMax
	}

	if config.GetValueStringDefault("gole.kafka.producer.retry-backoff", "100ms") != "100ms" {
		t, err := time.ParseDuration(Cfg.Producer.RetryBackoff)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.producer.retry-backoff】异常，%v", err.Error())
		} else {
			kafkaConfig.Producer.Retry.Backoff = t
		}
	}

	if config.GetValueBoolDefault("gole.kafka.producer.return-errors", true) != true {
		kafkaConfig.Producer.Return.Errors = Cfg.Producer.ReturnErrors
	}

	if config.GetValueBoolDefault("gole.kafka.producer.return-success", false) != false {
		kafkaConfig.Producer.Return.Successes = Cfg.Producer.ReturnSuccess
	}

	if config.GetValueIntDefault("gole.kafka.producer.compression-level", -1000) != -1000 {
		kafkaConfig.Producer.CompressionLevel = Cfg.Producer.CompressionLevel
	}

	if config.GetValueStringDefault("gole.kafka.producer.transaction-timeout", "1m") != "1m" {
		t, err := time.ParseDuration(Cfg.Producer.TransactionTimeout)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.producer.transaction-timeout】异常，%v", err.Error())
		} else {
			kafkaConfig.Producer.Transaction.Timeout = t
		}
	}

	if config.GetValueIntDefault("gole.kafka.producer.transaction-retry-max", 50) != 50 {
		kafkaConfig.Producer.Transaction.Retry.Max = Cfg.Producer.TransactionRetryMax
	}

	if config.GetValueStringDefault("gole.kafka.producer.transaction-retry-backoff", "100ms") != "100ms" {
		t, err := time.ParseDuration(Cfg.Producer.TransactionRetryBackoff)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.producer.transaction-retry-backoff】异常，%v", err.Error())
		} else {
			kafkaConfig.Producer.Transaction.Retry.Backoff = t
		}
	}

	//============================= consumer =============================
	if config.GetValueIntDefault("gole.kafka.consumer.fetch-min", 1) != 1 {
		kafkaConfig.Consumer.Fetch.Min = Cfg.Consumer.FetchMin
	}

	if config.GetValueIntDefault("gole.kafka.consumer.fetch-default", 1048576) != 1048576 {
		kafkaConfig.Consumer.Fetch.Default = Cfg.Consumer.FetchDefault
	}

	if config.GetValueStringDefault("gole.kafka.consumer.retry-backoff", "2s") != "2s" {
		t, err := time.ParseDuration(Cfg.Consumer.RetryBackoff)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.consumer.retry-backoff】异常，%v", err.Error())
		} else {
			kafkaConfig.Consumer.Retry.Backoff = t
		}
	}

	if config.GetValueStringDefault("gole.kafka.consumer.max-wait-time", "500ms") != "500ms" {
		t, err := time.ParseDuration(Cfg.Consumer.MaxWaitTime)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.consumer.max-wait-time】异常，%v", err.Error())
		} else {
			kafkaConfig.Consumer.MaxWaitTime = t
		}
	}

	if config.GetValueStringDefault("gole.kafka.consumer.max-processing-time", "100ms") != "100ms" {
		t, err := time.ParseDuration(Cfg.Consumer.MaxProcessingTime)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.consumer.max-processing-time】异常，%v", err.Error())
		} else {
			kafkaConfig.Consumer.MaxProcessingTime = t
		}
	}

	if config.GetValueBoolDefault("gole.kafka.consumer.return-errors", false) != false {
		kafkaConfig.Consumer.Return.Errors = Cfg.Consumer.ReturnErrors
	}

	if config.GetValueBoolDefault("gole.kafka.consumer.offsets-auto-commit-enable", false) != false {
		kafkaConfig.Consumer.Offsets.AutoCommit.Enable = Cfg.Consumer.OffsetsAutoCommitEnable
	}

	if config.GetValueStringDefault("gole.kafka.consumer.offsets-auto-commit-interval", "1s") != "1s" {
		t, err := time.ParseDuration(Cfg.Consumer.OffsetsAutoCommitInterval)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.consumer.offsets-auto-commit-interval】异常，%v", err.Error())
		} else {
			kafkaConfig.Consumer.Offsets.AutoCommit.Interval = t
		}
	}

	if config.GetValueIntDefault("gole.kafka.consumer.offsets-initial", -1) != -1 {
		kafkaConfig.Consumer.Offsets.Initial = Cfg.Consumer.OffsetsInitial
	}

	if config.GetValueIntDefault("gole.kafka.consumer.offsets-retry-max", 3) != 3 {
		kafkaConfig.Consumer.Offsets.Retry.Max = Cfg.Consumer.OffsetsRetryMax
	}

	//============================= consumer.group =============================
	if config.GetValueStringDefault("gole.kafka.consumer.group.session-timeout", "10s") != "10s" {
		t, err := time.ParseDuration(Cfg.Consumer.Group.SessionTimeout)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.consumer.group.session-timeout】异常，%v", err.Error())
		} else {
			kafkaConfig.Consumer.Group.Session.Timeout = t
		}
	}

	if config.GetValueStringDefault("gole.kafka.consumer.group.heartbeat-interval", "3s") != "3s" {
		t, err := time.ParseDuration(Cfg.Consumer.Group.HeartbeatInterval)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.consumer.group.heartbeat-interval】异常，%v", err.Error())
		} else {
			kafkaConfig.Consumer.Group.Heartbeat.Interval = t
		}
	}

	if config.GetValueStringDefault("gole.kafka.consumer.group.rebalance-timeout", "60s") != "60s" {
		t, err := time.ParseDuration(Cfg.Consumer.Group.RebalanceTimeout)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.consumer.group.rebalance-timeout】异常，%v", err.Error())
		} else {
			kafkaConfig.Consumer.Group.Rebalance.Timeout = t
		}
	}

	if config.GetValueIntDefault("gole.kafka.consumer.group.rebalance-retry-max", 4) != 4 {
		kafkaConfig.Consumer.Group.Rebalance.Retry.Max = Cfg.Consumer.Group.RebalanceRetryMax
	}

	if config.GetValueStringDefault("gole.kafka.consumer.group.rebalance-retry-backoff", "2s") != "2s" {
		t, err := time.ParseDuration(Cfg.Consumer.Group.RebalanceRetryBackoff)
		if err != nil {
			logger.Warnf("读取配置【gole.kafka.consumer.group.rebalance-retry-backoff】异常，%v", err.Error())
		} else {
			kafkaConfig.Consumer.Group.Rebalance.Retry.Backoff = t
		}
	}

	if config.GetValueBoolDefault("gole.kafka.consumer.group.reset-invalid-offsets", true) != true {
		kafkaConfig.Consumer.Group.ResetInvalidOffsets = Cfg.Consumer.Group.ResetInvalidOffsets
	}
	return kafkaConfig
}

func getKafkaVersion(kafkaVersion string) sarama.KafkaVersion {
	if !regexp.MustCompile(`^[Vv]\d+_\d+_\d+_\d+$`).MatchString(kafkaVersion) {
		logger.Error("gole.kafka.version 版本不合法：" + kafkaVersion)
		return sarama.V1_0_0_0
	}

	var one, two, three, four uint
	v := [4]*uint{&one, &two, &three, &four}
	_, err := fmt.Sscanf(strings.ToUpper(kafkaVersion), "V%d_%d_%d_%d", v[0], v[1], v[2], v[3])
	if err != nil {
		logger.Errorf("%v", err.Error())
		return sarama.V1_0_0_0
	}

	var kafkaV sarama.KafkaVersion
	if one == 0 {
		kafkaV, err = sarama.ParseKafkaVersion(fmt.Sprintf("0.%d.%d.%d", two, three, four))
	} else {
		kafkaV, err = sarama.ParseKafkaVersion(fmt.Sprintf("%d.%d.%d", one, two, three))
	}
	if err != nil {
		logger.Errorf("异常：%v", err.Error())
		return sarama.V1_0_0_0
	}
	return kafkaV
}
