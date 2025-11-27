package grpc

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

type ConfigOfGrpcServer struct {
	// 当前服务名
	ServiceName string
	// 监听的端口
	Port int
	// MaxConcurrentStreams返回一个ServerOption，该选项将对每个ServerTransport的并发流数量施加限制。
	MaxConcurrentStreams uint32
	// MaxRecvMsgSize 返回一个ServerOption，用于设置服务器可以接收的最大消息大小（以字节为单位）。如果没有设置，gRPC将使用默认的4MB。
	MaxReceiveMessageSize int
	// MaxSendMsgSize返回一个ServerOption，用于设置服务器可以发送的最大消息大小（以字节为单位）。如果没有设置，gRPC将使用默认的`math。MaxInt32。
	MaxSendMessageSize int
	// KeepaliveParams返回一个ServerOption，为服务器设置keepalive和最大年龄参数。
	KeepaliveParams keepalive.ServerParameters
	// InitialWindowSize返回一个ServerOption，用于设置流的窗口大小。窗口大小的下限为64K，任何小于此值的值都将被忽略。
	InitialWindowSize int32
	// InitialConnWindowSize返回一个ServerOption，用于设置连接的窗口大小。窗口大小的下限为64K，任何小于此值的值都将被忽略。
	InitialConnWindowSize int32
	// WriteBufferSize决定在线路上进行写入之前可以批处理多少数据。此缓冲区的默认值为32KB。零值或负值将禁用写缓冲区，以便每次写入都在底层连接上。注意：发送呼叫可能不会直接转换为写入。
	WriteBufferSize int
	// ReadBufferSize允许您设置读取缓冲区的大小，这决定了一次读取系统调用最多可以读取多少数据。此缓冲区的默认值为32KB。零值或负值将禁用连接的读取缓冲区，以便数据帧生成器可以直接访问底层连接。
	ReadBufferSize int
	// SharedWriteBuffer允许重用每个连接的传输写缓冲区。如果此选项设置为true，则每个连接都会在刷新线路上的数据后释放缓冲区。
	// 注意：本API是实验性的，可能会在稍后释放。
	SharedWriteBuffer bool
	// ConnectionTimeout返回一个ServerOption，用于设置所有新连接的连接建立超时时间（包括HTTP/2握手）。如果未设置，默认值为120秒。零值或负值将导致立即超时。
	// 注意：本API是实验性的，可能会在稍后释放。
	ConnectionTimeout time.Duration
	// MaxHeaderListSizeServerOption是一个ServerOption，用于设置服务器准备接受的标头列表的最大（未压缩）大小。
	MaxHeaderListSize *uint32
	// HeaderTableSize返回一个ServerOption，用于设置流的动态头表的大小。
	// 注意：此API是实验性的，可能会在以后的版本中更改或删除。
	HeaderTableSize *uint32
	// NumStreamWorkers返回一个ServerOption，用于设置应用于处理传入流的worker goroutines的数量。将其设置为零（默认）将禁用workers并为每个流生成一个新的goroutine。
	// 注意：此API是实验性的，可能会在以后的版本中更改或删除。
	NumServerWorkers uint32
	// WaitForHandlers使Stop等待所有未完成的方法处理程序退出，然后返回。如果为false，Stop将在所有连接关闭后立即返回，但方法处理程序可能仍在运行。默认情况下，Stop不会等待方法处理程序返回。
	// 注意：此API是实验性的，可能会在以后的版本中更改或删除。
	WaitForHandlers bool
}

func generateGrpcServerOpts(serviceName string, forMulti bool) []grpc.ServerOption {
	if config.Loaded == false {
		logger.Warn("配置未加载")
		return []grpc.ServerOption{}
	}
	var cfgOfGrpcServer ConfigOfGrpcServer
	if !forMulti {
		if config.GetValueBoolDefault("gole.server.grpc.enable", false) {
			err := config.GetValueObject("gole.server.grpc", &cfgOfGrpcServer)
			if err != nil {
				logger.Warn("读取server-grpc配置异常")
				return []grpc.ServerOption{}
			}
		}
	} else {
		if config.GetValueBoolDefault(fmt.Sprintf("gole.server.grpc.multi.%v.enable", serviceName), false) {
			err := config.GetValueObject(fmt.Sprintf("gole.server.grpc.multi.%v", serviceName), &cfgOfGrpcServer)
			if err != nil {
				logger.Warn("读取server-grpc配置异常")
				return []grpc.ServerOption{}
			}
		}
	}

	var options []grpc.ServerOption
	if cfgOfGrpcServer.MaxConcurrentStreams != 0 {
		options = append(options, grpc.MaxConcurrentStreams(cfgOfGrpcServer.MaxConcurrentStreams))
	}

	if cfgOfGrpcServer.MaxReceiveMessageSize != 0 {
		options = append(options, grpc.MaxRecvMsgSize(cfgOfGrpcServer.MaxReceiveMessageSize))
	}

	if cfgOfGrpcServer.MaxSendMessageSize != 0 {
		options = append(options, grpc.MaxSendMsgSize(cfgOfGrpcServer.MaxSendMessageSize))
	}

	if cfgOfGrpcServer.KeepaliveParams.MaxConnectionIdle != 0 ||
		cfgOfGrpcServer.KeepaliveParams.MaxConnectionAge != 0 ||
		cfgOfGrpcServer.KeepaliveParams.MaxConnectionAgeGrace != 0 ||
		cfgOfGrpcServer.KeepaliveParams.Time != 0 ||
		cfgOfGrpcServer.KeepaliveParams.Timeout != 0 {
		options = append(options, grpc.KeepaliveParams(cfgOfGrpcServer.KeepaliveParams))
	}

	if cfgOfGrpcServer.InitialWindowSize != 0 {
		options = append(options, grpc.InitialWindowSize(cfgOfGrpcServer.InitialWindowSize))
	}

	if cfgOfGrpcServer.InitialConnWindowSize != 0 {
		options = append(options, grpc.InitialConnWindowSize(cfgOfGrpcServer.InitialConnWindowSize))
	}

	if cfgOfGrpcServer.WriteBufferSize != 0 {
		options = append(options, grpc.WriteBufferSize(cfgOfGrpcServer.WriteBufferSize))
	}

	if cfgOfGrpcServer.ReadBufferSize != 0 {
		options = append(options, grpc.ReadBufferSize(cfgOfGrpcServer.ReadBufferSize))
	}

	if cfgOfGrpcServer.SharedWriteBuffer != false {
		options = append(options, grpc.SharedWriteBuffer(cfgOfGrpcServer.SharedWriteBuffer))
	}

	if cfgOfGrpcServer.ConnectionTimeout != 0 {
		options = append(options, grpc.ConnectionTimeout(cfgOfGrpcServer.ConnectionTimeout))
	}

	if cfgOfGrpcServer.MaxHeaderListSize != nil {
		options = append(options, grpc.MaxHeaderListSize(*cfgOfGrpcServer.MaxHeaderListSize))
	}

	if cfgOfGrpcServer.HeaderTableSize != nil {
		options = append(options, grpc.HeaderTableSize(*cfgOfGrpcServer.HeaderTableSize))
	}

	if cfgOfGrpcServer.NumServerWorkers != 0 {
		options = append(options, grpc.NumStreamWorkers(cfgOfGrpcServer.NumServerWorkers))
	}

	if cfgOfGrpcServer.WaitForHandlers != false {
		options = append(options, grpc.WaitForHandlers(cfgOfGrpcServer.WaitForHandlers))
	}

	// 添加opentelemetry的埋点配置
	if config.GetValueBoolDefault("gole.opentelemetry.enable", true) {
		options = append(options, grpc.StatsHandler(otelgrpc.NewServerHandler()))
	}

	var unaryList []grpc.UnaryServerInterceptor
	var streamList []grpc.StreamServerInterceptor

	// 添加grpc的各种测量指标
	if !forMulti {
		if config.GetValueBoolDefault("gole.meter.grpc.enable", false) {
			streamList = append(streamList, grpc_prometheus.StreamServerInterceptor)
			unaryList = append(unaryList, grpc_prometheus.UnaryServerInterceptor)
		}
	} else {
		if config.GetValueBoolDefault(fmt.Sprintf("gole.meter.grpc.multi.%v.enable", serviceName), false) {
			streamList = append(streamList, grpc_prometheus.StreamServerInterceptor)
			unaryList = append(unaryList, grpc_prometheus.UnaryServerInterceptor)
		}
	}

	// 添加埋点处理
	unaryList = append(unaryList, TraceUnaryInterceptorOfServer())

	// 添加异常处理拦截器
	streamList = append(streamList, ErrorStreamInterceptorOfServer())
	unaryList = append(unaryList, ErrorUnaryInterceptorOfServer())

	// 添加方法调用拦截处理
	unaryList = append(unaryList, RspHandleUnaryInterceptorOfServer())

	options = append(options, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamList...)))
	options = append(options, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryList...)))
	return options
}
