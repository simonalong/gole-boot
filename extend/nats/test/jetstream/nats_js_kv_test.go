package nats

import (
	"context"
	"fmt"
	jetstream "github.com/nats-io/nats.go/jetstream"
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/logger"
	"testing"
	"time"
)

// 使用环境变量：gole.profiles.active=consumer1
func TestNatsJsKv(t *testing.T) {
	_, js, err := baseNats.GetJetStreamClient()
	if err != nil {
		logger.Fatal(err)
		return
	}

	kv := GetBucket(js, "base_grpc_server_default_cbb-mid-srv-idgen", 3*time.Second)
	ctx := context.Background()

	//time.Sleep(1 * time.Second)
	//kv.Put(ctx, "key1", []byte("value1"))
	//time.Sleep(1 * time.Second)
	//valueEntity, _ := kv.Get(ctx, "key1")
	//fmt.Println(string(valueEntity.Value()))
	//
	//time.Sleep(1 * time.Second)
	//valueEntity, _ = kv.Get(ctx, "key1")
	//fmt.Println(string(valueEntity.Value()))
	//
	//kv.Put(ctx, "key1", []byte("value1"))
	//
	//time.Sleep(1 * time.Second)
	valueEntity, _ := kv.Get(ctx, "cbb-mid-srv-idgen_ad63a2")
	if valueEntity != nil {
		fmt.Println(string(valueEntity.Value()))
	} else {
		fmt.Println("value entity is nil")
	}
}

func GetBucket(js *baseNats.JetStreamClient, name string, ttl time.Duration) jetstream.KeyValue {
	ctx := context.Background()
	if kv, _ := js.KeyValue(ctx, name); nil != kv {
		return kv
	}
	kv, _ := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		// 桶名字
		Bucket: name,
		// 保存key的实效性
		TTL: ttl,
	})
	return kv
}
