package server

import (
	"github.com/simonalong/gole-boot/event"
	"github.com/simonalong/gole-boot/server/grpc"
	"github.com/simonalong/gole-boot/server/http"
	"github.com/simonalong/gole-boot/server/sys"
	"github.com/simonalong/gole-boot/server/tcp"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/listener"
	"github.com/simonalong/gole/logger"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func init() {
	config.Load()
	sys.PrintBanner()
}

func Run() {
	serverReady := sync.WaitGroup{}
	haveGrpcServer := false
	haveTcpServer := false
	haveHttpServer := false

	var serverTypes []string

	// grpc
	if config.GetValueBoolDefault("gole.server.grpc.enable", false) {
		haveGrpcServer = true
		serverReady.Add(1)
		listener.AddListener(event.EventOfServerGrpcRunFinish, func(e listener.BaseEvent) {
			serverReady.Done()
			serverTypes = append(serverTypes, "grpc")
		})
	}

	// tcp
	if config.GetValueBoolDefault("gole.server.tcp.enable", false) {
		haveTcpServer = true
		serverReady.Add(1)
		listener.AddListener(event.EventOfServerTcpRunFinish, func(e listener.BaseEvent) {
			serverReady.Done()
			serverTypes = append(serverTypes, "tcp")
		})
	}

	// http
	if http.ServerIsEnable() {
		haveHttpServer = true
		serverReady.Add(1)
		listener.AddListener(event.EventOfServerAllHttpRunFinish, func(e listener.BaseEvent) {
			serverReady.Done()
			serverTypes = append(serverTypes, "http")
		})
	}

	if haveHttpServer {
		go http.RunServer()
	}

	if haveGrpcServer {
		go grpc.RunServer()
	}

	if haveTcpServer {
		go tcp.RunServer()
	}

	// 服务都启动完成，则发送服务启动完成事件
	serverReady.Wait()
	logger.Infof("服务【%v】启动完成", strings.Join(serverTypes, "，"))
	listener.PublishEvent(event.ServerRunFinishEvent{})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-quit

	time.Sleep(time.Second)
}
