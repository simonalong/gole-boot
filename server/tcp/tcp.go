package tcp

import (
	"errors"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	listener2 "github.com/simonalong/gole-boot/event"
	_ "github.com/simonalong/gole-boot/otel"
	"github.com/simonalong/gole-boot/server/sys"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/global"
	"github.com/simonalong/gole/listener"
	"github.com/simonalong/gole/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var rootHandler *BaseReceiver
var rootCodecGenerate CodecGenerate
var rootConnectCloseHandler ConnectCloseHandler
var rootConnectOpenHandler ConnectOpenHandler

var CfgOfTcpServer ConfigOfTcpServer
var ConnectIsUnAvailableErr = errors.New("当前链接已经不可用")

type ConfigOfTcpServer struct {
	// 监听的端口
	Port int
	// 连接活跃判决时长
	ActiveJudgeDuration time.Duration
}

func init() {
	config.Load()
	sys.PrintBanner()

	if config.Loaded && config.GetValueBoolDefault("gole.server.tcp.enable", false) {
		err := config.GetValueObject("gole.server.tcp", &CfgOfTcpServer)
		if err != nil {
			logger.Warn("读取server-tcp配置异常")
			return
		}
	}

	if !config.GetValueBoolDefault("gole.server.tcp.enable", false) {
		return
	}

	// 设置默认配置
	if CfgOfTcpServer.Port == 0 {
		CfgOfTcpServer.Port = 80
	}

	// 支持opentelemetry
	var hook Hook
	if config.GetValueBoolDefault("gole.opentelemetry.enable", true) {
		hook = Hook{
			Tracer: global.Tracer,
		}
	} else {
		hook = Hook{}
	}

	rootHandler = &BaseReceiver{
		activeJudgeDuration: CfgOfTcpServer.ActiveJudgeDuration,
		hook:                hook,
	}

	// 添加事件监听：配置文件变更
	listener.AddListener(config.EventOfConfigChange, configChangeListenerOfTcp)
	sys.PrintBanner()
}

func RunServer() {
	if !config.GetValueBoolDefault("gole.server.tcp.enable", false) {
		logger.Errorf("tcp服务配置为false，不允许启动")
		return
	}

	port := config.GetValueIntDefault("gole.server.tcp.port", 80)
	listener.AddListener(listener2.EventOfServerTcpRunFinish, func(event listener.BaseEvent) {
		logger.Infof("tcp服务启动完成，端口号: %d", port)
	})

	logger.Infof("开始启动tcp服务")

	listener.PublishEvent(listener2.ServerTcpRunStartEvent{})

	// 添加监控指标
	initMeter()

	// 创建实例开启
	graceRun(port)
}

func graceRun(port int) {
	go func() {
		gnetConfigs := getGnetConfigs()
		gnetConfigs = append(gnetConfigs, gnet.WithLogger(&GnetLogger{}))

		if err := gnet.Run(rootHandler, fmt.Sprintf("tcp://0.0.0.0:%d", port), gnetConfigs...); err != nil {
			logger.Errorf("tcp启动服务异常 (%v)", err)
			return
		} else {
			// 发送服务关闭事件
			listener.PublishEvent(listener2.ServerTcpStopEvent{})
		}
	}()
	// 发送服务启动事件
	listener.PublishEvent(listener2.ServerTcpRunFinishEvent{})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-quit

	logger.Warn("tcp服务端准备关闭...")
	// 发送服务端关闭事件
	listener.PublishEvent(listener2.ServerTcpStopEvent{})
	logger.Warn("tcp服务端退出")
}

func getGnetConfigs() []gnet.Option {
	var options []gnet.Option
	if !config.GetValueBoolDefault("gole.server.tcp.enable", false) {
		return options
	}

	if config.GetValueBool("gole.server.tcp.gnet.multicore") {
		options = append(options, gnet.WithMulticore(true))
	}

	if val := config.GetValueIntDefault("gole.server.tcp.gnet.num-event-loop", 0); val != 0 {
		options = append(options, gnet.WithNumEventLoop(val))
	}

	if val := config.GetValueIntDefault("gole.server.tcp.gnet.lb", -1); val != -1 {
		switch val {
		case 0:
			options = append(options, gnet.WithLoadBalancing(gnet.RoundRobin))
		case 1:
			options = append(options, gnet.WithLoadBalancing(gnet.LeastConnections))
		case 2:
			options = append(options, gnet.WithLoadBalancing(gnet.SourceAddrHash))
		}
	}

	if config.GetValueBool("gole.server.tcp.gnet.reuse-addr") {
		options = append(options, gnet.WithReuseAddr(true))
	}

	if config.GetValueBool("gole.server.tcp.gnet.reuse-port") {
		options = append(options, gnet.WithReusePort(true))
	}

	if val := config.GetValueIntDefault("gole.server.tcp.gnet.multicast-interface-index", -1); val != -1 {
		options = append(options, gnet.WithMulticastInterfaceIndex(val))
	}

	// ============================= 服务器端和客户端的选项 =============================
	if val := config.GetValueIntDefault("gole.server.tcp.gnet.read-buffer-cap", -1); val != -1 {
		options = append(options, gnet.WithReadBufferCap(val))
	}

	if val := config.GetValueIntDefault("gole.server.tcp.gnet.write-buffer-cap", -1); val != -1 {
		options = append(options, gnet.WithWriteBufferCap(val))
	}

	if config.GetValueBool("gole.server.tcp.gnet.lock-os-thread") {
		options = append(options, gnet.WithLockOSThread(true))
	}

	if config.GetValueBool("gole.server.tcp.gnet.ticker") {
		options = append(options, gnet.WithTicker(true))
	}

	if val := config.GetValueString("gole.server.tcp.gnet.tcp-keep-alive"); val != "" {
		duration, err := time.ParseDuration(val)
		if err != nil {
			logger.Error("解析配置【gole.server.tcp.gnet.tcp-keep-alive】异常", err.Error())
		} else {
			options = append(options, gnet.WithTCPKeepAlive(duration))
		}
	}

	if val := config.GetValueIntDefault("gole.server.tcp.gnet.tcp-no-delay", -1); val != -1 {
		switch val {
		case 0:
			options = append(options, gnet.WithTCPNoDelay(gnet.TCPNoDelay))
		case 1:
			options = append(options, gnet.WithTCPNoDelay(gnet.TCPDelay))
		}
	}

	if val := config.GetValueIntDefault("gole.server.tcp.gnet.socket-receive-buffer", -1); val != -1 {
		options = append(options, gnet.WithSocketRecvBuffer(val))
	}

	if val := config.GetValueIntDefault("gole.server.tcp.gnet.socket-send-buffer", -1); val != -1 {
		options = append(options, gnet.WithSocketSendBuffer(val))
	}

	if config.GetValueBool("gole.server.tcp.gnet.edge-triggered-io") {
		options = append(options, gnet.WithTicker(true))
	}
	return options
}

func SetDecoder(codeGenerateFun CodecGenerate) {
	rootCodecGenerate = codeGenerateFun
}

func ConnectOpen(connectOpenHandler ConnectOpenHandler) {
	rootConnectOpenHandler = connectOpenHandler
}

func ConnectClose(connectCloseHandler ConnectCloseHandler) {
	rootConnectCloseHandler = connectCloseHandler
}

func Receive(dataHandler ReceiveHandler) {
	if !config.GetValueBoolDefault("gole.server.tcp.enable", false) {
		logger.Error("server-tcp配置没有开启，则设置数据处理器不生效")
		return
	}
	rootHandler.receiveHandler = dataHandler
}

func Send(con gnet.Conn, data []byte) error {
	connectContext := con.Context()
	if connectContext == nil {
		return errors.New("当前链接已经不可用")
	}
	codec := connectContext.(Decoder)
	if codec == nil {
		return errors.New("当前链接没有配置编解码")
	}
	_, err := con.Write(data)
	if err != nil {
		logger.Errorf("rsp写入异常：%v", err)
		return err
	}
	return nil
}

// GetOnlineCount 获取实时的连接个数
func GetOnlineCount() int {
	if rootHandler != nil {
		return rootHandler.eng.CountConnections()
	}
	return 0
}

// GetActiveStatus 获取活跃状态
func GetActiveStatus() bool {
	if rootHandler != nil {
		return time.Now().Before(rootHandler.lastActiveTime.Add(rootHandler.activeJudgeDuration))
	}
	return false
}

func configChangeListenerOfTcp(event listener.BaseEvent) {
	ev := event.(config.EventOfChange)
	if ev.Key != "gole.server.tcp.active-judge-duration" {
		return
	}

	if ev.Value != "" {
		duration, err := time.ParseDuration(ev.Value)
		if err != nil {
			logger.Errorf("配置值不合法，变更失效：%v", err)
			return
		}
		if rootHandler != nil {
			rootHandler.activeJudgeDuration = duration
		}
	}
}
