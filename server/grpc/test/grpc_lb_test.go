package test

import (
	"context"
	"github.com/simonalong/gole-boot/server/grpc"
	service2 "github.com/simonalong/gole-boot/server/grpc/test/original/service"
	"github.com/simonalong/gole/logger"
	"testing"
)

type ServerImpl1 struct {
	service2.UnimplementedGreeterServer
}
type ServerImpl2 struct {
	service2.UnimplementedGreeterServer
}
type ServerImpl3 struct {
	service2.UnimplementedGreeterServer
}

func (s *ServerImpl1) SayHello(ctx context.Context, in *service2.HelloRequest) (*service2.HelloReply, error) {
	logger.Info("Received：", in.GetName())
	return &service2.HelloReply{Message: "Hello " + in.GetName() + " from server1"}, nil
}
func (s *ServerImpl2) SayHello(ctx context.Context, in *service2.HelloRequest) (*service2.HelloReply, error) {
	logger.Info("Received：", in.GetName())
	return &service2.HelloReply{Message: "Hello " + in.GetName() + " from server2"}, nil
}
func (s *ServerImpl3) SayHello(ctx context.Context, in *service2.HelloRequest) (*service2.HelloReply, error) {
	logger.Info("Received：", in.GetName())
	return &service2.HelloReply{Message: "Hello " + in.GetName() + " from server3"}, nil
}

// gole.profiles.active=lb-server1
func TestServer1(t *testing.T) {
	// 注册服务
	grpc.Register(&service2.Greeter_ServiceDesc, &ServerImpl1{})

	// 启动grpc服务
	grpc.RunServer()
}

// gole.profiles.active=lb-server2
func TestServer2(t *testing.T) {
	// 注册服务
	grpc.Register(&service2.Greeter_ServiceDesc, &ServerImpl2{})

	// 启动grpc服务
	grpc.RunServer()
}

// gole.profiles.active=lb-server3
func TestServer3(t *testing.T) {
	// 注册服务
	grpc.Register(&service2.Greeter_ServiceDesc, &ServerImpl3{})

	// 启动grpc服务
	grpc.RunServer()
}
