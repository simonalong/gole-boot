package test

import (
	"context"
	"github.com/simonalong/gole-boot/server/grpc"
	service2 "github.com/simonalong/gole-boot/server/grpc/test/original/service"
	"github.com/simonalong/gole/logger"
	"testing"
)

type ServerImpl struct {
	serviceName string
	service2.UnimplementedGreeterServer
}

func (s *ServerImpl) SayHello(ctx context.Context, in *service2.HelloRequest) (*service2.HelloReply, error) {
	logger.Infof("服务：%v, Received：%v", s.serviceName, in.GetName())
	return &service2.HelloReply{Message: "Hello " + in.GetName()}, nil
	//return nil, errorx.SC_DB_ERR.WithDetail("table not found")
}

// gole.profiles.active=server
// gole.profiles.active=server-otel
// gole.profiles.active=server-grpc-http
func TestGrpcServer(t *testing.T) {

	// 注册服务
	grpc.Register(&service2.Greeter_ServiceDesc, &ServerImpl{})

	// 启动grpc服务
	grpc.RunServer()
}
