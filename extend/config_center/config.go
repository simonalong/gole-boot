package config_center

var Cfg CenterConfig

type CenterConfig struct {
	// 配置的注册中心：目前只支持nats（其他的配置中心暂时不支持），会直接读取gole.nats的相关配置；默认：nats
	Register string
	// 服务名称；默认：读取 gole.application.name 配置
	ServiceName string
	// 分组；默认：default
	Group string
	// 配置内容类型：yaml，yml（其实也是yaml），json，properties，默认：yaml
	ConfigType string
}
