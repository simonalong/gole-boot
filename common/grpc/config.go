package grpc

var CfgOfGrpcRegisterCenter ConfigOfGrpcRegisterCenter

type ConfigOfGrpcRegisterCenter struct {
	// 是否启用，默认关闭
	Enable bool
	// 注册中心类型，默认：natsJs
	Type string `json:"type"`
}
