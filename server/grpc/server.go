package grpc

import (
	"fmt"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	cmap "github.com/orcaman/concurrent-map"
	common "github.com/simonalong/gole-boot/common/grpc"
	"github.com/simonalong/gole-boot/event"
	"github.com/simonalong/gole-boot/server/sys"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/listener"
	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/maps"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// key: grpc的server名字，value对象：*ServerOfGrpc
var serverOfGrpcMap cmap.ConcurrentMap
var serverLoadLock sync.Mutex

//const DefaultServiceName = "default"

func init() {
	config.Load()
	sys.PrintBanner()
	serverOfGrpcMap = cmap.New()

	if !config.Loaded {
		logger.Fatalf("配置加载失败")
		return
	}

	// 初始化注册中心
	initRegister()
}

func ServerForSingle() *ServerOfGrpc {
	serviceName := getSingleGrpcServiceName()
	if serverOfGrpc, have := serverOfGrpcMap.Get(serviceName); have {
		return serverOfGrpc.(*ServerOfGrpc)
	}
	serverLoadLock.Lock()
	defer serverLoadLock.Unlock()

	if serverOfGrpc, have := serverOfGrpcMap.Get(serviceName); have {
		return serverOfGrpc.(*ServerOfGrpc)
	}

	// 创建服务端对象
	serverOfGrpc := newServer(serviceName, false)
	if serverOfGrpc == nil {
		return nil
	}
	serverOfGrpcMap.Set(serviceName, serverOfGrpc)
	return serverOfGrpc
}

func ServerForMulti(serviceName string) *ServerOfGrpc {
	if serverOfGrpc, have := serverOfGrpcMap.Get(serviceName); have {
		return serverOfGrpc.(*ServerOfGrpc)
	}
	serverLoadLock.Lock()
	defer serverLoadLock.Unlock()

	if serverOfGrpc, have := serverOfGrpcMap.Get(serviceName); have {
		return serverOfGrpc.(*ServerOfGrpc)
	}

	// 创建服务端对象
	serverOfGrpc := newServer(serviceName, true)
	if serverOfGrpc == nil {
		return nil
	}
	serverOfGrpcMap.Set(serviceName, serverOfGrpc)
	return serverOfGrpc
}

func Register(sd *grpc.ServiceDesc, ss any) {
	grpcServer := ServerForSingle()
	grpcServer.Register(sd, ss)
}

func newServer(serviceName string, forMulti bool) *ServerOfGrpc {
	if !config.Loaded {
		logger.Warn("配置未加载")
		return nil
	}
	if !forMulti {
		if !config.GetValueBoolDefault("gole.server.grpc.enable", false) {
			logger.Warn("grpc服务未开启，请开启gole.server.grpc.enable")
			return nil
		}
		port := config.GetValueIntDefault("gole.server.grpc.port", 9090)
		return &ServerOfGrpc{ServiceName: serviceName, Port: port, ServerInstance: grpc.NewServer(generateGrpcServerOpts(serviceName, forMulti)...)}
	} else {
		if !config.GetValueBoolDefault(fmt.Sprintf("gole.server.grpc.multi.%v.enable", serviceName), false) {
			logger.Warnf("grpc服务【%v】未开启，请开启gole.server.grpc.multi.%v.enable", serviceName, serviceName)
			return nil
		}
		port := config.GetValueIntDefault(fmt.Sprintf("gole.server.grpc.multi.%v.port", serviceName), 9090)
		return &ServerOfGrpc{ServiceName: serviceName, Port: port, ServerInstance: grpc.NewServer(generateGrpcServerOpts(serviceName, forMulti)...)}
	}
}

func RunServer() {
	if config.GetValueBoolDefault("gole.server.grpc.enable", false) {
		runServerForSingle()
		return
	}
	runMultiServer()
}

func initRegister() {
	if !config.GetValueBoolDefault("gole.grpc.register.enable", false) {
		return
	}
	// 创建注册中心
	common.InitRegisterCenter()

	// grpc服务端处理
	if config.GetValueBoolDefault("gole.server.grpc.enable", false) {
		serviceName := getSingleGrpcServiceName()
		if serviceName == "" {
			logger.Errorf("grpc服务未配置服务名，请配置gole.server.grpc.service-name")
			return
		}
		grpcServer := ServerForSingle()
		if grpcServer == nil {
			return
		}
		// 服务注册
		common.RegisterService(grpcServer.ServiceName, grpcServer.Port)

		// 添加事件监听机制
		listener.AddListener(event.EventOfServerGrpcStop, func(event listener.BaseEvent) {
			logger.Infof(fmt.Sprintf("监听服务【%v】退出事件，进行清理grpc的服务注册", grpcServer.ServiceName))
			common.UnRegisterService(grpcServer.ServiceName)
		})
	} else {
		serviceNames := getMultiServiceNames()
		for _, serviceNameTem := range serviceNames {
			if !config.GetValueBoolDefault(fmt.Sprintf("gole.server.grpc.multi.%v.enable", serviceNameTem), false) {
				logger.Errorf("grpc服务【%v】未配置启动 gole.server.grpc.multi.%v.enable", serviceNameTem, serviceNameTem)
				continue
			}
			grpcServer := ServerForMulti(serviceNameTem)
			if grpcServer == nil {
				return
			}
			// 服务注册
			common.RegisterService(serviceNameTem, grpcServer.Port)

			// 添加事件监听机制
			listener.AddListenerWithGroup(serviceNameTem, event.EventOfServerGrpcStop, func(event listener.BaseEvent) {
				logger.Infof(fmt.Sprintf("监听服务【%v】退出事件，进行清理grpc的服务注册", serviceNameTem))
				common.UnRegisterService(serviceNameTem)
			})
		}
	}
}

func getSingleGrpcServiceName() string {
	newServiceName := config.GetValueString("gole.server.grpc.service-name")
	if newServiceName != "" {
		return newServiceName
	}
	if val := config.BaseCfg.Application.Name; val != "" {
		return val
	} else {
		logger.Fatalf("gole.application.name 不可为空")
	}
	return ""
}

func runServerForSingle() {
	if !config.GetValueBoolDefault("gole.server.grpc.enable", false) {
		logger.Warn("grpc服务未开启，请开启 gole.server.grpc.enable")
		return
	}

	port := config.GetValueIntDefault("gole.server.grpc.port", 9090)

	listener.AddListener(event.EventOfServerGrpcRunFinish, func(ev listener.BaseEvent) {
		logger.Infof("grpc服务启动完成，端口号：%d", port)
	})

	grpcServer := ServerForSingle()
	if grpcServer == nil {
		return
	}

	if config.GetValueBoolDefault("gole.meter.grpc.enable", false) {
		grpcPrometheus.Register(grpcServer.ServerInstance)
	}

	logger.Debugf("grpc服务启动开始，port=%v", port)

	// 发送服务开启
	listener.PublishEvent(event.ServerGrpcRunStartEvent{})

	graceRunForSingle(grpcServer, port)
}

func runServerForMulti(serviceName string) {
	if !config.GetValueBoolDefault(fmt.Sprintf("gole.server.http.%v.enable", serviceName), true) {
		logger.Warn(fmt.Sprintf("http服务【%v】未开启，请开启%v", serviceName, fmt.Sprintf("gole.server.http.%v.enable", serviceName)))
		return
	}

	port := config.GetValueIntDefault(fmt.Sprintf("gole.server.grpc.multi.%v.port", serviceName), 9090)

	listener.AddListenerWithGroup("*", event.EventOfServerGrpcRunFinish, func(ev listener.BaseEvent) {
		finishEvent, ok := ev.(event.ServerGrpcRunFinishEvent)
		if ok && finishEvent.ServiceName == serviceName {
			logger.Infof("grpc服务【%v】启动完成，端口号：%d", serviceName, port)
		}
	})

	grpcServer := ServerForMulti(serviceName)
	if grpcServer == nil {
		return
	}

	// 注意：这个prometheus只支持单个服务，因此这边会只有最后一个注册的生效，这个也是一个问题，后续再考虑解决
	if config.GetValueBoolDefault(fmt.Sprintf("gole.meter.grpc.multi.%v.enable", serviceName), false) {
		grpcPrometheus.Register(grpcServer.ServerInstance)
	}

	logger.Debugf("grpc服务【%v】启动开始，port=%v", serviceName, port)

	// 发送服务开启
	listener.PublishEvent(event.ServerGrpcRunStartEvent{ServiceName: serviceName})

	graceRunForMulti(grpcServer, port)
}

func runMultiServer() {
	multiServerCnt := sync.WaitGroup{}
	// 获取 serviceName 的列表
	serviceNames := getMultiServiceNames()
	if len(serviceNames) == 0 {
		return
	}
	multiServerCnt.Add(len(serviceNames))

	listener.AddListenerWithGroup("*", event.EventOfServerGrpcRunFinish, func(event listener.BaseEvent) {
		multiServerCnt.Done()
	})

	go func() {
		multiServerCnt.Wait()

		// 等所有的服务都启动完之后，这边发送一个默认http服务启动完成的信号
		listener.PublishEvent(event.ServerGrpcAllRunFinishEvent{})
	}()

	serverNum := len(serviceNames)
	for i := range serverNum - 1 {
		go runServerForMulti(serviceNames[i])
	}
	runServerForMulti(serviceNames[serverNum-1])
}

func getMultiServiceNames() []string {
	// 获取 serviceName 的列表
	serviceNames := config.GetValue("gole.server.grpc.multi")
	if serviceNames == nil {
		return nil
	}
	dataMap, _ := maps.From(serviceNames)
	var serviceNams []string
	for _, k := range dataMap.Keys() {
		if k == "" {
			continue
		}
		if open, have := dataMap.AsDeepMap().GetBool(k + ".enable"); !have || !open {
			continue
		}
		serviceNams = append(serviceNams, k)
	}
	return serviceNams
}

func graceRunForSingle(serverOfGrpc *ServerOfGrpc, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Errorf("listen 配置端口异常 (%v)", err)
		return
	}

	go func() {
		if err := serverOfGrpc.ServerInstance.Serve(lis); err != nil {
			logger.Errorf("grpc服务启动服务异常 (%v)", err)
			// 发送服务关闭事件
			listener.PublishEvent(event.ServerGrpcStopEvent{ServiceName: listener.DefaultGroup})
			return
		}
	}()
	// 发送服务启动事件
	listener.PublishEvent(event.ServerGrpcRunFinishEvent{ServiceName: listener.DefaultGroup})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-quit

	logger.Warn("grpc服务端准备关闭...")

	// 发送服务端关闭事件
	listener.PublishEvent(event.ServerGrpcStopEvent{ServiceName: listener.DefaultGroup})

	serverOfGrpc.ServerInstance.GracefulStop()
	logger.Warn("grpc服务端退出")
}

func graceRunForMulti(serverOfGrpc *ServerOfGrpc, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Errorf("listen 配置端口异常 (%v)", err)
		return
	}

	go func() {
		if err := serverOfGrpc.ServerInstance.Serve(lis); err != nil {
			logger.Errorf("grpc服务【%v】启动服务异常 (%v)", serverOfGrpc.ServiceName, err)
			// 发送服务关闭事件
			listener.PublishEvent(event.ServerGrpcStopEvent{ServiceName: serverOfGrpc.ServiceName})
			return
		}
	}()
	// 发送服务启动事件
	listener.PublishEvent(event.ServerGrpcRunFinishEvent{ServiceName: serverOfGrpc.ServiceName})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-quit

	logger.Warnf("grpc服务端【%v】准备关闭...", serverOfGrpc.ServiceName)

	// 发送服务端关闭事件
	listener.PublishEvent(event.ServerGrpcStopEvent{ServiceName: serverOfGrpc.ServiceName})

	serverOfGrpc.ServerInstance.GracefulStop()
	logger.Warnf("grpc服务端【%v】退出", serverOfGrpc.ServiceName)
}
