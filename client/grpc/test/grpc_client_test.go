package test

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/simonalong/gole-boot/client/grpc"
	service2 "github.com/simonalong/gole-boot/server/grpc/test/original/service"
	httpServer "github.com/simonalong/gole-boot/server/http"
	"github.com/simonalong/gole/logger"
	"testing"
	"time"
)

// gole.profiles.active=client
// gole.profiles.active=client-otel
func TestGrpcClient(t *testing.T) {
	// 获取grpc客户端
	conn, err := grpc.NewClientConn("demo1")
	if err != nil {
		return
	}

	c := service2.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rspData, err := c.SayHello(ctx, &service2.HelloRequest{Name: "name1"})
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Info(rspData)
	//time.Sleep(5 * time.Second)
}

// gole.profiles.active=client
// gole.profiles.active=client-otel
// gole.profiles.active=client-http
func TestGrpcClientWithHttp(t *testing.T) {

	// 获取grpc客户端
	conn, err := grpc.NewClientConn("demo1")
	if err != nil {
		return
	}

	grpcClient := service2.NewGreeterClient(conn)

	httpServer.Get("api/call/grpc", func(c *gin.Context) (any, error) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		rspData, err := grpcClient.SayHello(ctx, &service2.HelloRequest{Name: "name1"})
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		logger.Info(rspData)
		return rspData, nil
	})

	httpServer.RunServer()
}
