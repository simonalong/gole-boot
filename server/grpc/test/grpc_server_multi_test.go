package test

import (
	"github.com/simonalong/gole-boot/server/grpc"
	service2 "github.com/simonalong/gole-boot/server/grpc/test/original/service"
	"testing"
)

// gole.profiles.active=server-multi
// gole.profiles.active=reg-server-multi
func TestGrpcMultiServer(t *testing.T) {
	grpcServer1 := grpc.ServerForMulti("demo1-service")
	grpcServer1.Register(&service2.Greeter_ServiceDesc, &ServerImpl{serviceName: "demo1-service"})

	grpcServer2 := grpc.ServerForMulti("demo2-service")
	grpcServer2.Register(&service2.Greeter_ServiceDesc, &ServerImpl{serviceName: "demo2-service"})

	// 启动grpc服务
	grpc.RunServer()
}
