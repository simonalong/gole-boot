package test

import (
	"context"
	"github.com/simonalong/gole-boot/client/grpc"
	service2 "github.com/simonalong/gole-boot/server/grpc/test/original/service"
	"github.com/simonalong/gole/logger"
	"strconv"
	"testing"
	"time"
)

// gole.profiles.active=lb-client
func TestClient(t *testing.T) {
	// 获取grpc客户端
	conn, err := grpc.NewClientConn("demo-client")
	if err != nil {
		return
	}

	cliengGrpc := service2.NewGreeterClient(conn)

	for i := 0; i < 1000000; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		rsp, err := cliengGrpc.SayHello(ctx, &service2.HelloRequest{Name: "name" + strconv.Itoa(i)})
		if err != nil {
			logger.Error(err)
			return
		}

		logger.Info(rsp)
		time.Sleep(2 * time.Second)
	}
}

// gole.profiles.active=lb-client-direct
func TestClientDirct(t *testing.T) {
	// 获取grpc客户端
	conn, err := grpc.NewClientConn("demo-client")
	if err != nil {
		return
	}

	cliengGrpc := service2.NewGreeterClient(conn)

	for i := 0; i < 1000000; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		rsp, err := cliengGrpc.SayHello(ctx, &service2.HelloRequest{Name: "name" + strconv.Itoa(i)})
		if err != nil {
			logger.Error(err)
			return
		}

		logger.Info(rsp)
		time.Sleep(2 * time.Second)
	}
}

// gole.profiles.active=lb-client-multi
func TestClientMulti(t *testing.T) {
	// 获取grpc客户端
	conn1, err := grpc.NewClientConn("demo1-client")
	if err != nil {
		return
	}
	conn2, err := grpc.NewClientConn("demo2-client")
	if err != nil {
		return
	}

	clientGrpc1 := service2.NewGreeterClient(conn1)
	clientGrpc2 := service2.NewGreeterClient(conn2)

	go func() {
		for i := 0; i < 1000000; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			rsp, err := clientGrpc1.SayHello(ctx, &service2.HelloRequest{Name: "name" + strconv.Itoa(i)})
			if err != nil {
				logger.Error(err)
				return
			}

			logger.Info(rsp)
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for i := 0; i < 1000000; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			rsp, err := clientGrpc2.SayHello(ctx, &service2.HelloRequest{Name: "name" + strconv.Itoa(i)})
			if err != nil {
				logger.Error(err)
				return
			}

			logger.Info(rsp)
			time.Sleep(2 * time.Second)
		}
	}()

	time.Sleep(2 * time.Hour)
}
