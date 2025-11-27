package nats

import (
	"github.com/nats-io/nats.go/jetstream"
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/config"
	baseTime "github.com/simonalong/gole/time"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 使用本例子之前请使用 ./deploy/xxx 文件夹中的部署进行启动即可

// 使用环境变量：gole.profiles.active=all-none-consumer
func TestNatsJsConfigWithAllNoneConsumer(t *testing.T) {
	config.LoadFile("./application-all-streams.yaml")
	config.GetValueObject("gole.nats", &baseNats.CfgOfNats)

	assert.Equal(t, "nats://127.0.0.1:4222", baseNats.CfgOfNats.Url)
	// --------------- 基本配置 ---------------
	assert.Equal(t, "stream-name1", baseNats.CfgOfNats.Jetstream.Streams[0].Name)
	assert.Equal(t, "描述1", baseNats.CfgOfNats.Jetstream.Streams[0].Description)
	// 这个顺序这边有点乱，每次发现输出的有点不一样，这个暂时不是很好解决
	//assert.Equal(t, "test1.*.req", baseNats.CfgOfNats.Jetstream.Streams[0].Subjects[0])
	//assert.Equal(t, "test2.*.req", baseNats.CfgOfNats.Jetstream.Streams[0].Subjects[1])

	// --------------- 流的消息限制和留存策略 ---------------
	assert.Equal(t, 1000, baseNats.CfgOfNats.Jetstream.Streams[0].MaxConsumers)
	assert.Equal(t, int64(1000), baseNats.CfgOfNats.Jetstream.Streams[0].MaxMsgs)
	assert.Equal(t, int64(1000), baseNats.CfgOfNats.Jetstream.Streams[0].MaxBytes)
	assert.Equal(t, "3h0m0s", baseNats.CfgOfNats.Jetstream.Streams[0].MaxAge.String())
	assert.Equal(t, int64(1000), baseNats.CfgOfNats.Jetstream.Streams[0].MaxMsgsPerSubject)
	assert.Equal(t, int32(1000), baseNats.CfgOfNats.Jetstream.Streams[0].MaxMsgSize)
	assert.Equal(t, jetstream.RetentionPolicy(2), baseNats.CfgOfNats.Jetstream.Streams[0].Retention)
	assert.Equal(t, jetstream.DiscardPolicy(1), baseNats.CfgOfNats.Jetstream.Streams[0].Discard)
	assert.Equal(t, true, baseNats.CfgOfNats.Jetstream.Streams[0].DiscardNewPerSubject)
	assert.Equal(t, jetstream.StorageType(1), baseNats.CfgOfNats.Jetstream.Streams[0].Storage)

	// --------------- 集群配置 ---------------
	assert.Equal(t, 3, baseNats.CfgOfNats.Jetstream.Streams[0].Replicas)
	assert.Equal(t, "cluster-name", baseNats.CfgOfNats.Jetstream.Streams[0].Placement.Cluster)
	assert.Equal(t, "tag1", baseNats.CfgOfNats.Jetstream.Streams[0].Placement.Tags[0])
	//assert.Equal(t, "tag2", baseNats.CfgOfNats.Jetstream.Streams[0].Placement.Tags[1])

	// --------------- 来源配置 ---------------
	assert.Equal(t, "mirror-name", baseNats.CfgOfNats.Jetstream.Streams[0].Mirror.Name)
	assert.Equal(t, "source0-name", baseNats.CfgOfNats.Jetstream.Streams[0].Sources[0].Name)

	// --------------- 操作配置 ---------------
	assert.Equal(t, true, baseNats.CfgOfNats.Jetstream.Streams[0].Sealed)
	assert.Equal(t, true, baseNats.CfgOfNats.Jetstream.Streams[0].DenyDelete)
	assert.Equal(t, true, baseNats.CfgOfNats.Jetstream.Streams[0].DenyPurge)
	assert.Equal(t, true, baseNats.CfgOfNats.Jetstream.Streams[0].AllowRollup)

	// --------------- 内容控制 ---------------
	assert.Equal(t, "3m0s", baseNats.CfgOfNats.Jetstream.Streams[0].Duplicates.String())
	assert.Equal(t, jetstream.StoreCompression(1), baseNats.CfgOfNats.Jetstream.Streams[0].Compression)
	assert.Equal(t, uint64(1), baseNats.CfgOfNats.Jetstream.Streams[0].FirstSeq)
	assert.Equal(t, "source", baseNats.CfgOfNats.Jetstream.Streams[0].SubjectTransform.Source)
	assert.Equal(t, "destination", baseNats.CfgOfNats.Jetstream.Streams[0].SubjectTransform.Destination)
	assert.Equal(t, "source", baseNats.CfgOfNats.Jetstream.Streams[0].RePublish.Source)
	assert.Equal(t, "destination", baseNats.CfgOfNats.Jetstream.Streams[0].RePublish.Destination)
	assert.Equal(t, true, baseNats.CfgOfNats.Jetstream.Streams[0].RePublish.HeadersOnly)
	assert.Equal(t, true, baseNats.CfgOfNats.Jetstream.Streams[0].AllowDirect)
	assert.Equal(t, true, baseNats.CfgOfNats.Jetstream.Streams[0].MirrorDirect)

	// --------------- 消费配置 ---------------
	assert.Equal(t, true, baseNats.CfgOfNats.Jetstream.Streams[0].NoAck)
	assert.Equal(t, "3s", baseNats.CfgOfNats.Jetstream.Streams[0].ConsumerLimits.InactiveThreshold.String())
	assert.Equal(t, 12, baseNats.CfgOfNats.Jetstream.Streams[0].ConsumerLimits.MaxAckPending)

	// --------------- 其他配置 ---------------
	dataMap := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	assert.Equal(t, dataMap, baseNats.CfgOfNats.Jetstream.Streams[0].Metadata)
}

// 使用环境变量：gole.profiles.active=all-have-consumer
func TestNatsJsConfigWithAllHaveConsumer(t *testing.T) {
	config.LoadFile("./application-all-consumers.yaml")
	config.GetValueObject("gole.nats", &baseNats.CfgOfNats)

	assert.Equal(t, "nats://127.0.0.1:4222", baseNats.CfgOfNats.Url)

	// --------------- 基本配置 ---------------
	assert.Equal(t, "consumer1", baseNats.CfgOfNats.Jetstream.Consumers[0].Name)
	assert.Equal(t, "consumer1", baseNats.CfgOfNats.Jetstream.Consumers[0].Durable)
	assert.Equal(t, "消费者描述", baseNats.CfgOfNats.Jetstream.Consumers[0].Description)
	assert.Equal(t, true, baseNats.CfgOfNats.Jetstream.Consumers[0].Order)
	assert.Equal(t, jetstream.DeliverPolicy(3), baseNats.CfgOfNats.Jetstream.Consumers[0].DeliverPolicy)
	assert.Equal(t, uint64(12), baseNats.CfgOfNats.Jetstream.Consumers[0].OptStartSeq)
	assert.Equal(t, "2024-07-30", baseTime.TimeToStringYmd(*baseNats.CfgOfNats.Jetstream.Consumers[0].OptStartTime))
	assert.Equal(t, "3s", baseNats.CfgOfNats.Jetstream.Consumers[0].BackOff[0].String())
	assert.Equal(t, "test.req", baseNats.CfgOfNats.Jetstream.Consumers[0].FilterSubject)
	assert.Equal(t, jetstream.ReplayPolicy(1), baseNats.CfgOfNats.Jetstream.Consumers[0].ReplayPolicy)

	// --------------- 确认配置 ---------------
	assert.Equal(t, jetstream.AckPolicy(2), baseNats.CfgOfNats.Jetstream.Consumers[0].AckPolicy)
	assert.Equal(t, "5m0s", baseNats.CfgOfNats.Jetstream.Consumers[0].AckWait.String())

	//--------------- 投递限制 ---------------
	assert.Equal(t, uint64(1000), baseNats.CfgOfNats.Jetstream.Consumers[0].RateLimit)
	assert.Equal(t, "10s", baseNats.CfgOfNats.Jetstream.Consumers[0].SampleFrequency)
	assert.Equal(t, 1000, baseNats.CfgOfNats.Jetstream.Consumers[0].MaxDeliver)
	assert.Equal(t, 1000, baseNats.CfgOfNats.Jetstream.Consumers[0].MaxWaiting)
	assert.Equal(t, 1000, baseNats.CfgOfNats.Jetstream.Consumers[0].MaxAckPending)
	assert.Equal(t, 1000, baseNats.CfgOfNats.Jetstream.Consumers[0].MaxRequestBatch)
	assert.Equal(t, "10s", baseNats.CfgOfNats.Jetstream.Consumers[0].MaxRequestExpires.String())
	assert.Equal(t, 1000, baseNats.CfgOfNats.Jetstream.Consumers[0].MaxRequestMaxBytes)
	assert.Equal(t, true, baseNats.CfgOfNats.Jetstream.Consumers[0].HeadersOnly)

	// --------------- 其他配置 ---------------
	assert.Equal(t, "10s", baseNats.CfgOfNats.Jetstream.Consumers[0].InactiveThreshold.String())
	assert.Equal(t, 3, baseNats.CfgOfNats.Jetstream.Consumers[0].Replicas)
	assert.Equal(t, true, baseNats.CfgOfNats.Jetstream.Consumers[0].MemoryStorage)
	assert.Equal(t, "test.req", baseNats.CfgOfNats.Jetstream.Consumers[0].FilterSubjects[0])
	dataMap := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	assert.Equal(t, dataMap, baseNats.CfgOfNats.Jetstream.Consumers[0].Metadata)
}
