package grpc

import (
	"context"
	"errors"
	"github.com/nats-io/nats.go/jetstream"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/simonalong/gole-boot/constants"
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/goid"
	"github.com/simonalong/gole/logger"
	baseTime "github.com/simonalong/gole/time"
	"github.com/simonalong/gole/util"
	"time"
)

var registerCenter *baseNats.JetStreamClient

// var currentUniqueServiceName string
var uniqueServiceNameMap = cmap.New()

const heartbeatIntervalSecond = 60

const RegisterPrePath = "base_grpc_server_default_"

type ServiceRegisterConfig struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

func InitRegisterCenter() {
	err := config.GetValueObject("gole.grpc.register", &CfgOfGrpcRegisterCenter)
	if err != nil {
		logger.Warn("读取 nats 配置异常")
		return
	}

	if CfgOfGrpcRegisterCenter.Type == constants.RegisterTypeNats {
		_, js, err := baseNats.GetJetStreamClient()
		if err != nil {
			logger.Warnf("注册中心连接失败：%v", err.Error())
			return
		}
		registerCenter = &baseNats.JetStreamClient{JetStream: js}
	}
}

func RegisterService(serviceName string, port int) {
	if registerCenter == nil {
		logger.Fatalf("注册中心实例为空")
		return
	}

	kv, err := getAndUpdateBucket(registerCenter, generateServiceRegPath(serviceName), 3*heartbeatIntervalSecond*time.Second)
	if err != nil {
		logger.Fatalf("创建kvBueckt失败：%v", err)
		return
	}

	internalIp, err := util.GetIntranetIp()
	if err != nil {
		logger.Fatalf("获取ip失败：%v", err)
		return
	}
	registerConfig := ServiceRegisterConfig{
		Ip:   internalIp,
		Port: port,
	}

	// 生成唯一服务标识
	uniqueServiceName, err := generateUniqueServiceName(kv, serviceName)
	if err != nil {
		logger.Fatalf("生成服务唯一标识失败：%v", err)
		return
	}
	_, err = kv.PutString(context.Background(), uniqueServiceName, util.ToJsonString(registerConfig))
	if err != nil {
		logger.Errorf("服务【%v】注册到注册中心失败：%v", uniqueServiceName, err)
		return
	}

	logger.Infof("服务【%v】注册到注册中心成功，唯一id【%v】", serviceName, uniqueServiceName)

	uniqueServiceNameMap.Set(serviceName, uniqueServiceName)

	// 开启心跳
	go func() {
		timer := baseTime.NewTimerWithFire(heartbeatIntervalSecond, func(t *baseTime.Timer) {
			_, err = kv.PutString(context.Background(), uniqueServiceName, util.ToJsonString(registerConfig))
			if err != nil {
				logger.Errorf("服务【%v】心跳刷新失败：%v", serviceName, err)
				return
			}
		})
		timer.Start()
	}()
}

func UnRegisterService(serviceName string) {
	if registerCenter == nil {
		logger.Fatalf("注册中心实例为空")
		return
	}

	kv, err := getAndUpdateBucket(registerCenter, generateServiceRegPath(serviceName), 15*time.Second)
	if err != nil {
		logger.Errorf("创建kvBueckt失败：%v", err)
		return
	}
	currentUniqueServiceName, have := uniqueServiceNameMap.Get(serviceName)
	if !have {
		logger.Errorf("服务【%v】未注册", serviceName)
		return
	}
	err = kv.Purge(context.Background(), util.ToString(currentUniqueServiceName))
	if err != nil {
		logger.Errorf("删除服务【%v】唯一id【%v】失败：%v", serviceName, currentUniqueServiceName, err)
		return
	}
}

func getAndUpdateBucket(js *baseNats.JetStreamClient, name string, ttl time.Duration) (jetstream.KeyValue, error) {
	ctx := context.Background()
	if kv, err := js.KeyValue(ctx, name); nil != kv {
		return kv, err
	}
	kv, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		// 桶名字
		Bucket: name,
		// 保存key的实效性
		TTL: ttl,
	})
	if err != nil {
		logger.Errorf("创建kvBueckt【name=%v】失败：%v", name, err)
		return nil, err
	}
	return kv, nil
}

func generateUniqueServiceName(kvBucket jetstream.KeyValue, serviceName string) (string, error) {
	serviceUniqueName := serviceName + "_" + goid.GenerateUUIDFullString()[:6]
	_, err := kvBucket.Get(context.Background(), serviceUniqueName)
	if nil != err {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			return serviceUniqueName, nil
		}
		logger.Errorf("获取服务的信息【%v】失败，%v", serviceUniqueName, err.Error())
		return "", err
	}
	return generateUniqueServiceName(kvBucket, serviceName)
}

func generateServiceRegPath(serviceName string) string {
	return RegisterPrePath + serviceName
}
