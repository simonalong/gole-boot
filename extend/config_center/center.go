package config_center

import (
	"context"
	"errors"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/simonalong/gole-boot/extend/nats"
	_ "github.com/simonalong/gole-boot/otel"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/logger"
)

var CfgCenterClient *nats.JetStreamClient
var latestVersion uint64

func init() {
	config.Load()

	if config.Loaded && config.GetValueBoolDefault("gole.config_center.enable", false) {
		err := config.GetValueObject("gole.emqx", &Cfg)
		if err != nil {
			logger.Warnf("读取emqx配置异常, %v", err.Error())
			return
		}
	}

	if config.Loaded {
		initConfigCenter()
	}
}

func initConfigCenter() {
	// 补充配置
	appendCfg()

	// 初始化配置中心
	initClient()

	// 读取全量配置到本地
	kvClient, err := loadAllConfig()
	if err != nil {
		logger.Errorf("读取全量配置到本地失败：%v", err.Error())
		return
	}

	// 监听配置中心的变化
	err = watchConfigChange(kvClient)
	if err != nil {
		logger.Errorf("监听配置中心的变化：%v", err.Error())
	}
}

func appendCfg() {
	if Cfg.ServiceName == "" {
		if val := config.BaseCfg.Application.Name; val != "" {
			Cfg.ServiceName = val
		} else {
			logger.Fatalf("gole.application.name 不可为空")
		}
	}

	if Cfg.Group == "" {
		Cfg.Group = "default"
	}

	if Cfg.ConfigType == "" {
		Cfg.ConfigType = "yaml"
	}
}

func initClient() {
	_, _js, err := nats.GetJetStreamClient()
	if err != nil {
		logger.Fatal("配置中心初始化失败：", err)
		return
	}
	CfgCenterClient = _js
}

func loadAllConfig() (jetstream.KeyValue, error) {
	bucketOfCfg, err := GetBucket(CfgCenterClient, Cfg.ServiceName)
	if err != nil {
		logger.Warnf("读取配置中心配置失败：bucketName=%v, 异常：%v", Cfg.ServiceName, err)
		return nil, err
	}
	if bucketOfCfg == nil {
		logger.Warnf("配置中心中不存在配置：%v，则读取本地的application.{yaml/yml/properties/json}（或者application-{profile}.{yaml/yml/properties/json}）的配置", Cfg.ServiceName)
		return nil, errors.New("配置中心不存在配置")
	}

	valueEntity, err := bucketOfCfg.Get(context.Background(), Cfg.Group)
	if err != nil {
		logger.Warnf("读取配置中心配置失败：%v", err)
		return nil, err
	}

	content := string(valueEntity.Value())
	switch Cfg.ConfigType {
	case "yaml", "yml":
		config.AppendYamlContent(content)
		break
	case "properties":
		config.AppendPropertyContent(content)
		break
	case "json":
		config.AppendJsonContent(content)
		break
	}
	return bucketOfCfg, nil
}

func watchConfigChange(bucketClient jetstream.KeyValue) error {
	watcher, err := bucketClient.WatchAll(context.Background(), jetstream.UpdatesOnly())
	if err != nil {
		return err
	}

	go func() {
		defer watcher.Stop()
		for {
			select {
			case v := <-watcher.Updates():
				if v.Key() != Cfg.Group {
					logger.Warnf("不是自己服务的配置，不关注：%v", v.Key())
					return
				}

				if v.Revision() <= latestVersion {
					logger.Warnf("不是最新的配置，不关注：Revision=%v", v.Revision())
					return
				}
				latestVersion = v.Revision()

				content := string(v.Value())
				switch Cfg.ConfigType {
				case "yaml", "yml":
					config.AppendYamlContent(content)
					break
				case "properties":
					config.AppendPropertyContent(content)
					break
				case "json":
					config.AppendJsonContent(content)
					break
				}
			}
		}
	}()
	return nil
}

func GetBucket(js *nats.JetStreamClient, bucketName string) (jetstream.KeyValue, error) {
	ctx := context.Background()
	if kv, err := js.KeyValue(ctx, bucketName); nil != kv {
		return kv, err
	}
	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		// 桶名字
		Bucket: bucketName,
	})
	return kv, err
}
