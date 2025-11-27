package grpc

import (
	gogrpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type ConfigOfGrpcClient struct {
	// 直连的话，配置这个；直连优先级高于负载均衡
	Host string
	Port int

	// 负载均衡算法，目前仅支持两个：
	// round_robin：（默认）连接到它看到的所有地址，并依次向每个后端发送一个RPC。例如，第一个RPC将发送到backend-1，第二个RPC将发送到backend-2，第三个RPC将再次发送到backend-1。
	// pick_first：尝试连接到第一个地址，如果连接成功，则将其用于所有RPC，如果连接失败，则尝试下一个地址（并继续这样做，直到一个连接成功），连接成功则一直保持
	LoadBalance string
	// 使用注册中心服务发现，可以使用这个
	ServiceName string

	// WithAuthority返回一个DialOption，指定用作：authority伪标头和身份验证握手中的服务器名称的值。
	Authority string
	// WithDisableServiceConfig返回一个DialOption，该选项使gRPC忽略解析器提供的任何服务配置，并向解析器提供不获取服务配置的提示。
	// 请注意，此拨号选项仅禁用解析器的服务配置。如果提供了默认服务配置，gRPC将使用默认服务配置。
	DisableServiceConfig bool
	// WithDisableRetry返回一个禁用重试的DialOption，即使服务配置启用了重试。这不会影响透明重试，如果没有数据写入线路或远程服务器未处理RPC，透明重试将自动发生。
	DisableRetry bool
	// WithDisableHealthCheck禁用此ClientConn的所有子Conn的LB通道健康检查。
	// 注意：此API是实验性的，可能会在以后的版本中更改或删除。
	DisableHealthCheck bool
	// WithDefaultServiceConfig返回一个DialOption，用于配置默认服务配置，该配置将在以下情况下使用：
	// 1.还使用WithDisableServiceConfig，或
	// 2.名称解析器不提供服务配置或提供无效的服务配置。
	// 参数s是默认服务配置的JSON表示形式。有关服务配置的更多信息，请参阅：https://github.com/grpc/grpc/blob/master/doc/service_config.md有关使用的简单示例，请参阅：examples/features/load_balancing/client/main.go
	DefaultServiceConfigRawJSON *string
	// WithIdleTimeout返回一个DialOption，用于配置通道的空闲超时。如果通道在配置的超时时间内处于空闲状态，即没有正在进行的RPC，也没有启动新的RPC，则通道将进入空闲模式，因此名称解析器和负载均衡器将关闭。当调用Connect（）方法或启动RPC时，通道将退出空闲模式。
	// 如果在拨号时未设置此拨号选项，则将使用默认超时30分钟，并且可以通过传递超时0来禁用空闲状态。
	// 注意：此API是实验性的，可能会在以后的版本中更改或删除。
	IdleTimeout time.Duration

	// ConnectParams定义了连接和重试的参数。鼓励用户使用此类型，而不是上面定义的BackoffConfig类型。请参阅此处了解更多详细信息：
	// https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md.
	//
	// #实验
	// 注意：此类型是实验性的，可能会在以后的版本中更改或删除。
	ConnectParams grpc.ConnectParams
}

func GenerateGrpcClientOptsDefault() []grpc.DialOption {
	var options []grpc.DialOption
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 添加opentelemetry的埋点配置
	if config.GetValueBoolDefault("gole.opentelemetry.enable", true) {
		options = append(options, grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	}

	var unaryList []grpc.UnaryClientInterceptor
	var streamList []grpc.StreamClientInterceptor

	// 添加meter的测量配置
	if config.GetValueBoolDefault("gole.meter.grpc.enable", false) {
		unaryList = append(unaryList, grpcprometheus.UnaryClientInterceptor)
		streamList = append(streamList, grpcprometheus.StreamClientInterceptor)
	}

	// 给客户端添加链路埋点上下文
	unaryList = append(unaryList, TraceUnaryInterceptorOfClient)
	// 异常处理的转换
	unaryList = append(unaryList, ErrorUnaryInterceptorOfClient)
	options = append(options, grpc.WithUnaryInterceptor(gogrpcmiddleware.ChainUnaryClient(unaryList...)))
	options = append(options, grpc.WithStreamInterceptor(gogrpcmiddleware.ChainStreamClient(streamList...)))

	// 负载均衡：使用随机策略
	options = append(options, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	return options
}

func generateGrpcClientOpts(serviceName string) []grpc.DialOption {
	var options []grpc.DialOption
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	cfgOfGrpcClient := ConfigOfGrpcClient{}
	err := config.GetValueObject("gole.grpc."+serviceName, &cfgOfGrpcClient)
	if err != nil {
		logger.Errorf("加载grpc的服务[%s]配置异常：%v", serviceName, err)
		return options
	}

	if cfgOfGrpcClient.Authority != "" {
		options = append(options, grpc.WithAuthority(cfgOfGrpcClient.Authority))
	}

	if cfgOfGrpcClient.DisableServiceConfig != false {
		options = append(options, grpc.WithDisableServiceConfig())
	}

	if cfgOfGrpcClient.DisableRetry != false {
		options = append(options, grpc.WithDisableRetry())
	}

	if cfgOfGrpcClient.DisableHealthCheck != false {
		options = append(options, grpc.WithDisableHealthCheck())
	}

	if cfgOfGrpcClient.DefaultServiceConfigRawJSON != nil {
		options = append(options, grpc.WithDefaultServiceConfig(*cfgOfGrpcClient.DefaultServiceConfigRawJSON))
	}

	if cfgOfGrpcClient.IdleTimeout != 0 {
		options = append(options, grpc.WithIdleTimeout(cfgOfGrpcClient.IdleTimeout))
	}

	options = append(options, grpc.WithConnectParams(cfgOfGrpcClient.ConnectParams))

	// 添加opentelemetry的埋点配置
	if config.GetValueBoolDefault("gole.opentelemetry.enable", true) {
		options = append(options, grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	}

	var unaryList []grpc.UnaryClientInterceptor
	var streamList []grpc.StreamClientInterceptor

	// 添加meter的测量配置
	if config.GetValueBoolDefault("gole.meter.grpc.enable", false) {
		unaryList = append(unaryList, grpcprometheus.UnaryClientInterceptor)
		streamList = append(streamList, grpcprometheus.StreamClientInterceptor)
	}

	// 给客户端添加链路埋点上下文
	unaryList = append(unaryList, TraceUnaryInterceptorOfClient)

	// 异常处理的转换
	unaryList = append(unaryList, ErrorUnaryInterceptorOfClient)

	options = append(options, grpc.WithUnaryInterceptor(gogrpcmiddleware.ChainUnaryClient(unaryList...)))
	options = append(options, grpc.WithStreamInterceptor(gogrpcmiddleware.ChainStreamClient(streamList...)))

	// 负载均衡：使用随机策略
	if cfgOfGrpcClient.LoadBalance == "" || cfgOfGrpcClient.LoadBalance == "round_robin" {
		options = append(options, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	}
	return options
}
