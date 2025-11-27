package grpc

import (
	"github.com/simonalong/gole/logger"
	"google.golang.org/grpc"
)

type ServerOfGrpc struct {
	ServiceName    string
	Port           int
	ServerInstance *grpc.Server
}

func (serverOfGrpc *ServerOfGrpc) Register(sd *grpc.ServiceDesc, ss any) {
	if serverOfGrpc == nil {
		return
	}
	serverInstance := serverOfGrpc.ServerInstance
	if serverInstance == nil {
		logger.Errorf("服务【%v】grpcServer为空，请检查配置", serverOfGrpc.ServiceName)
		return
	}
	serverInstance.RegisterService(sd, ss)
}
