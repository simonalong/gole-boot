package grpc

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go/jetstream"
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/util"
)

func DiscoverService() {
	if !config.GetValueBoolDefault("gole.grpc.register.enable", false) {
		logger.Warn("注册中心配置未激活，使用直连模式")
		return
	}
	// 获取服务列表
	AddDiscoverService(getServiceList())
}

func AddDiscoverService(serviceNames []string) {
	if len(serviceNames) == 0 {
		return
	}

	if registerCenter == nil {
		logger.Fatalf("注册中心为空，请查看配置")
		return
	}

	// 查询每个服务列表注册的信息
	serviceRegInfoListMap := getServiceRegInfoListMap(serviceNames)
	if serviceRegInfoListMap == nil || len(serviceRegInfoListMap) == 0 {
		logger.Errorf("获取服务列表【%v】注册信息为空", serviceNames)
		return
	}

	// 更新负载均衡里面服务列表和ip对应的关系
	if GlobalResolver == nil {
		InitResolver()
	}
	GlobalResolver.setAddressStore(serviceRegInfoListMap)
	//GlobalResolver.refreshState()

	// 添加监控对服务列表
	addWatchForServiceNameChg(serviceNames)
}

func getServiceList() []string {
	data := config.GetValue("gole.grpc")
	if data == nil {
		logger.Errorf("请检查gole.grpc配置")
		return nil
	}
	dataMap := util.ToMap(data)

	var keys []string
	for k := range dataMap {
		if k == "enable" || k == "register" {
			continue
		}
		keys = append(keys, util.ToString(k))
	}

	var serviceNameList []string
	for _, key := range keys {
		svcName := config.GetValueString("gole.grpc." + key + ".service-name")
		if svcName == "" {
			continue
		}
		serviceNameList = append(serviceNameList, svcName)
	}
	return serviceNameList
}

func getBucket(js *baseNats.JetStreamClient, name string) jetstream.KeyValue {
	ctx := context.Background()
	kv, err := js.KeyValue(ctx, name)
	if nil != err {
		logger.Errorf("获取kvBucket【%v】失败：%v", name, err)
		return nil
	}
	return kv
}

func getServiceRegInfoListMap(serviceNameRegList []string) map[string]map[string]string {
	serviceRegInfoListMap := map[string]map[string]string{}
	for _, serviceNameReg := range serviceNameRegList {
		kvBucket := getBucket(registerCenter, generateServiceRegPath(serviceNameReg))
		if kvBucket == nil {
			logger.Errorf("获取服务【%v】信息失败", serviceNameReg)
			return serviceRegInfoListMap
		}
		serviceNameUniqList, err := kvBucket.Keys(context.Background())
		if err != nil {
			logger.Errorf("注册中心服务【%v】不存在，请检查是否服务端启动否", serviceNameReg)
			return serviceRegInfoListMap
		}

		serviceRegMap := map[string]string{}
		for _, serviceNameUniq := range serviceNameUniqList {
			kvEntity, err := kvBucket.Get(context.Background(), serviceNameUniq)
			if err != nil {
				logger.Errorf("获取服务列表的存储信息get失败：%v", err.Error())
				return serviceRegInfoListMap
			}

			serviceRegConfig := ServiceRegisterConfig{}
			err = json.Unmarshal(kvEntity.Value(), &serviceRegConfig)
			if err != nil {
				logger.Errorf("解析服务注册信息失败：%v", err.Error())
				return serviceRegInfoListMap
			}
			serviceRegMap[serviceNameUniq] = fmt.Sprintf("%s:%d", serviceRegConfig.Ip, serviceRegConfig.Port)
		}
		serviceRegInfoListMap[serviceNameReg] = serviceRegMap
	}
	return serviceRegInfoListMap
}

func addWatchForServiceNameChg(serviceNameRegList []string) {
	for _, serviceNameReg := range serviceNameRegList {
		kvBucket := getBucket(registerCenter, generateServiceRegPath(serviceNameReg))
		if kvBucket == nil {
			return
		}
		keyWatcher, _ := kvBucket.WatchAll(context.Background())
		go func() {
			for {
				select {
				case kvs := <-keyWatcher.Updates():
					if nil != kvs {
						op := kvs.Operation()
						switch op {
						case jetstream.KeyValuePut:
							serviceRegConfig := ServiceRegisterConfig{}
							err := json.Unmarshal(kvs.Value(), &serviceRegConfig)
							if err != nil {
								logger.Errorf("解析服务注册信息失败：%v", err.Error())
								continue
							}
							updateFlag := GlobalResolver.addAddressStore(serviceNameReg, kvs.Key(), fmt.Sprintf("%s:%d", serviceRegConfig.Ip, serviceRegConfig.Port))
							if updateFlag {
								GlobalResolver.refreshState()
							}
						case jetstream.KeyValueDelete, jetstream.KeyValuePurge:
							GlobalResolver.deleteAddressStore(serviceNameReg, kvs.Key())
							GlobalResolver.refreshState()
						}
					}
				}
			}
		}()
	}
}
