package grpc

import (
	"errors"
	"fmt"
	common "github.com/simonalong/gole-boot/common/grpc"
	"github.com/simonalong/gole/bean"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/logger"
	"google.golang.org/grpc"
	"sync"
)

const BeanNameOfGrpcPrefix = "baseGrpc_"

var initLock sync.Mutex

func init() {
	config.Load()

	if !config.Loaded {
		logger.Fatalf("配置文件加载失败")
		return
	}

	if !config.GetValueBoolDefault("gole.grpc.enable", false) {
		return
	}

	if !config.GetValueBoolDefault("gole.grpc.register.enable", false) {
		logger.Warn("注册中心配置未激活，使用直连模式")
		return
	}

	// 创建注册中心
	common.InitRegisterCenter()

	// 获取监听的服务数据
	common.DiscoverService()
}

func GetClientConn(serviceName string) (*grpc.ClientConn, error) {
	grpcSrvName := fmt.Sprintf("%v%v", BeanNameOfGrpcPrefix, serviceName)
	if grpcClientOfSrv, ok := bean.GetBean(grpcSrvName).(*grpc.ClientConn); ok {
		return grpcClientOfSrv, nil
	}
	initLock.Lock()
	defer initLock.Unlock()
	if grpcClientOfSrv, ok := bean.GetBean(grpcSrvName).(*grpc.ClientConn); ok {
		return grpcClientOfSrv, nil
	}
	grpcClientOfSrv, err := NewClientConn(serviceName)
	if err != nil {
		return nil, err
	}
	bean.AddBean(grpcSrvName, grpcClientOfSrv)
	return grpcClientOfSrv, nil
}

func NewClientConn(serviceName string) (*grpc.ClientConn, error) {
	if serviceName == "" {
		return nil, errors.New("服务名不可为空")
	}

	targetUrl := generateClientTarget(serviceName)
	if targetUrl == "" {
		logger.Fatalf("服务【%v】grpc的目标url创建失败", serviceName)
		return nil, errors.New("grpc的目标url创建失败")
	}
	grpcServer, err := grpc.NewClient(targetUrl, generateGrpcClientOpts(serviceName)...)
	if err != nil {
		logger.Fatalf("创建grpc客户端异常：%v", err)
		return nil, err
	}
	return grpcServer, err
}

func generateClientTarget(serviceName string) string {
	host := config.GetValueString("gole.grpc." + serviceName + ".host")
	port := config.GetValueInt("gole.grpc." + serviceName + ".port")
	serviceNameOfReg := config.GetValueString("gole.grpc." + serviceName + ".service-name")
	if host != "" && port != 0 {
		return fmt.Sprintf("%s:%d", host, port)
	} else if serviceNameOfReg != "" {
		return fmt.Sprintf("gole:///%v", serviceNameOfReg)
	} else {
		return ""
	}
}
