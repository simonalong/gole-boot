package kafka

import "github.com/IBM/sarama"

var Cfg KafkaConfig

// ---------------------------- gole.kafka ----------------------------

type KafkaConfig struct {
	// 默认 sarama
	ClientId string
	// 默认 256
	ChannelBufferSize int
	// 默认 true
	ApiVersionsRequest bool
	// 默认 V1_0_0_0，版本格式为v{x}_{x}_{x}_{x}；版本是否存在请见 Shopify/sarama 代码中的util.go包的版本
	Version string

	Admin    KafkaAdmin
	Net      KafkaNet
	Metadata KafkaMetadata
	Producer KafkaProducer
	Consumer KafkaConsumer
}

type KafkaAdmin struct {
	// 默认5
	RetryMax int
	// 默认100ms
	RetryBackoff string
	// 默认3s
	Timeout string
}

type KafkaNet struct {
	// 默认5
	MaxOpenRequests int
	// 默认3s
	DialTimeout string
	// 默认3s
	ReadTimeout string
	// 默认3s
	WriteTimeout string
	// 默认true
	SASLHandshake bool `yaml:"SASL-handshake"`
	// 默认0，只有0和1两个
	SASLVersion int16 `yaml:"SASL-version"`
}

type KafkaMetadata struct {
	// 默认 3
	RetryMax int
	// 默认250ms
	RetryBackoff string
	// 默认10分钟，即10m
	RefreshFrequency string
	// 默认 true
	Full bool
	// 默认 true
	AllowAutoTopicCreation bool
}

type KafkaProducer struct {
	// 默认1000000
	MaxMessageBytes int
	// 默认1
	RequiredAcks sarama.RequiredAcks
	// 10s
	Timeout string
	// 默认3
	RetryMax int
	// 默认100ms
	RetryBackoff string
	// 默认true
	ReturnErrors bool
	// 默认false
	ReturnSuccess bool
	// 默认-1000
	CompressionLevel int
	// 默认1分钟
	TransactionTimeout string
	// 默认50
	TransactionRetryMax int
	// 默认100毫秒
	TransactionRetryBackoff string
}

type KafkaConsumer struct {
	// 默认1
	FetchMin int32
	// 默认1024*1024
	FetchDefault int32
	// 默认2s
	RetryBackoff string
	// 默认500ms
	MaxWaitTime string
	// 默认100ms
	MaxProcessingTime string
	// 默认false
	ReturnErrors bool
	// 默认false
	OffsetsAutoCommitEnable bool
	// 默认1秒
	OffsetsAutoCommitInterval string
	// 默认-1
	OffsetsInitial int64
	// 默认3
	OffsetsRetryMax int
	Group           KafkaConsumerGroup
}

type KafkaConsumerGroup struct {
	// 默认10s
	SessionTimeout string
	// 默认3s
	HeartbeatInterval string
	// 默认60s
	RebalanceTimeout string
	// 默认4
	RebalanceRetryMax int
	// 默认2秒
	RebalanceRetryBackoff string
	// 默认true
	ResetInvalidOffsets bool
}
