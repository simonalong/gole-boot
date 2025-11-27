package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	cmap "github.com/orcaman/concurrent-map"
	"github.com/simonalong/gole-boot/event"
	"github.com/simonalong/gole-boot/server/sys"
	"github.com/simonalong/gole/maps"

	"github.com/gin-contrib/pprof"
	"github.com/simonalong/gole/listener"

	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/util"

	"github.com/gin-gonic/gin"
	"github.com/simonalong/gole/logger"
)

type Method int

const (
	HmAll Method = iota
	HmGet
	HmPost
	HmPut
	HmPatch
	HmDelete
	HmOptions
	HmHead
)

const DefaultServiceName = "default"

// key: http的server名字，value对象：*http.ServerOfHttp
var serverOfHttpMap cmap.ConcurrentMap
var serverLoadLock sync.Mutex

// key: http的server名字，value对象：bool
var pprofHaveMap cmap.ConcurrentMap

func init() {
	serverOfHttpMap = cmap.New()
	pprofHaveMap = cmap.New()

	config.Load()
	sys.PrintBanner()
}

func Server(serviceName string) *ServerOfHttp {
	if serverOfHttp, have := serverOfHttpMap.Get(serviceName); have {
		return serverOfHttp.(*ServerOfHttp)
	}
	serverLoadLock.Lock()
	defer serverLoadLock.Unlock()

	if serverOfHttp, have := serverOfHttpMap.Get(serviceName); have {
		return serverOfHttp.(*ServerOfHttp)
	}

	// 创建服务端对象
	serverOfHttp := newServer(serviceName)
	if serverOfHttp == nil {
		return nil
	}
	serverOfHttpMap.Set(serviceName, serverOfHttp)

	// 初始化服务配置
	initServerConfig(serverOfHttp)
	return serverOfHttp
}

func RunServer() {
	if config.GetValueBoolDefault("gole.server.http.enable", false) {
		runServerWithServiceName(DefaultServiceName)
		return
	}
	runMultiServer()
}

func GetRouterGroup() *gin.RouterGroup {
	return GetRouterGroupWithServiceName(DefaultServiceName)
}

func GetRouterGroupWithServiceName(serviceName string) *gin.RouterGroup {
	serverOfHttp := Server(serviceName)
	if serverOfHttp == nil {
		return nil
	}
	return serverOfHttp.IRouter.(*gin.RouterGroup)
}

func GetEngine() *gin.Engine {
	return GetEngineWithServiceName(DefaultServiceName)
}

func GetEngineWithServiceName(serviceName string) *gin.Engine {
	serverOfHttp := Server(serviceName)
	if serverOfHttp == nil {
		return nil
	}
	return serverOfHttp.IRouter.(*gin.Engine)
}

// Use 功能同 AddMiddleware，只是为了使用gin的人方便理解，多了一个函数而已
func Use(handler ...gin.HandlerFunc) *ServerOfHttp {
	return AddMiddleware(handler...)
}

// AddMiddleware 功能同 Use，只是为了使用gin的人方便理解，多了一个函数而已
func AddMiddleware(handler ...gin.HandlerFunc) *ServerOfHttp {
	return AddMiddlewareWithServiceName(DefaultServiceName, handler...)
}

func AddMiddlewareWithServiceName(serviceName string, handler ...gin.HandlerFunc) *ServerOfHttp {
	httpServer := Server(serviceName)
	httpServer.AddMiddleware(handler...)
	return httpServer
}

func AddRouteWithServiceName(serviceName string, method Method, path string, handler RouteHandler) *ServerOfHttp {
	defaultHttpServer := Server(serviceName)
	return defaultHttpServer.AddRoute(method, path, handler)
}

func AddRouteGinHandlerWithServiceName(serviceName string, method Method, path string, handler gin.HandlerFunc) *ServerOfHttp {
	defaultHttpServer := Server(serviceName)
	return defaultHttpServer.AddRouteGinHandler(method, path, handler)
}

func AddRoute(method Method, path string, handler RouteHandler) *ServerOfHttp {
	return AddRouteWithServiceName(DefaultServiceName, method, path, handler)
}

func AddRouteGinHandler(path string, method Method, handler gin.HandlerFunc) *ServerOfHttp {
	return AddRouteGinHandlerWithServiceName(DefaultServiceName, method, path, handler)
}

func Post(path string, handler RouteHandler) *ServerOfHttp {
	return AddRoute(HmPost, getPathAppendApiModel(DefaultServiceName, path), handler)
}

func Delete(path string, handler RouteHandler) *ServerOfHttp {
	return AddRoute(HmDelete, getPathAppendApiModel(DefaultServiceName, path), handler)
}

func Put(path string, handler RouteHandler) *ServerOfHttp {
	return AddRoute(HmPut, getPathAppendApiModel(DefaultServiceName, path), handler)
}

func Patch(path string, handler RouteHandler) *ServerOfHttp {
	return AddRoute(HmPatch, getPathAppendApiModel(DefaultServiceName, path), handler)
}

func Head(path string, handler RouteHandler) *ServerOfHttp {
	return AddRoute(HmHead, getPathAppendApiModel(DefaultServiceName, path), handler)
}

func Get(path string, handler RouteHandler) *ServerOfHttp {
	return AddRoute(HmGet, getPathAppendApiModel(DefaultServiceName, path), handler)
}

func Options(path string, handler RouteHandler) *ServerOfHttp {
	return AddRoute(HmOptions, getPathAppendApiModel(DefaultServiceName, path), handler)
}

func All(path string, handler RouteHandler) *ServerOfHttp {
	return AddRoute(HmAll, getPathAppendApiModel(DefaultServiceName, path), handler)
}

func GroupWithServiceName(serviceName, path string, handlerFunc ...gin.HandlerFunc) *ServerOfHttp {
	defaultHttpServer := Server(serviceName)
	return defaultHttpServer.Group(path, handlerFunc...)
}

func Group(path string, handlerFunc ...gin.HandlerFunc) *ServerOfHttp {
	return GroupWithServiceName(DefaultServiceName, path, handlerFunc...)
}

func ServerIsEnable() bool {
	// 单http服务
	if config.GetValueBool("gole.server.http.enable") {
		return true
	}

	// 多http服务
	return len(getMultiServiceNames()) > 0
}

func runServerWithServiceName(serviceName string) {
	if serviceName == DefaultServiceName {
		if !config.GetValueBoolDefault("gole.server.http.enable", true) {
			logger.Warn(fmt.Sprintf("http服务未开启，请开启%v", fmt.Sprintf("gole.server.http.%v.enable", serviceName)))
			return
		}
	} else {
		if !config.GetValueBoolDefault(fmt.Sprintf("gole.server.http.%v.enable", serviceName), true) {
			logger.Warn(fmt.Sprintf("http服务【%v】未开启，请开启%v", serviceName, fmt.Sprintf("gole.server.http.%v.enable", serviceName)))
			return
		}
	}

	var port int
	if serviceName == DefaultServiceName {
		port = config.GetValueIntDefault("gole.server.http.port", 8080)
	} else {
		port = config.GetValueIntDefault(fmt.Sprintf("gole.server.http.multi.%v.port", serviceName), 8080)
	}

	listener.AddListenerWithGroup(serviceName, event.EventOfServerHttpRunFinish, func(event listener.BaseEvent) {
		if serviceName == DefaultServiceName {
			logger.Infof("http服务启动完成，端口号：%d", port)
		} else {
			logger.Infof("http服务【%v】启动完成，端口号：%d", serviceName, port)
		}
	})

	if Server(serviceName) == nil {
		return
	}

	listener.PublishEvent(event.ServerHttpRunStartEvent{ServiceName: serviceName})
	// 发送服务启动事件
	if serviceName == DefaultServiceName {
		logger.Debugf("http服务启动开始，port=%v", port)
	} else {
		logger.Debugf("http服务【%v】启动开始，port=%v", serviceName, port)
	}
	graceRun(serviceName, port)
}

func runMultiServer() {
	multiServerCnt := sync.WaitGroup{}
	// 获取serviceName的列表
	serviceNames := getMultiServiceNames()
	if len(serviceNames) == 0 {
		return
	}
	multiServerCnt.Add(len(serviceNames))

	listener.AddListenerWithGroup("*", event.EventOfServerHttpRunFinish, func(event listener.BaseEvent) {
		multiServerCnt.Done()
	})

	go func() {
		multiServerCnt.Wait()

		// 等所有的服务都启动完之后，这边发送一个默认http服务启动完成的信号
		listener.PublishEvent(event.ServerHttpAllRunFinishEvent{})
	}()

	serverNum := len(serviceNames)
	for i := range serverNum - 1 {
		go runServerWithServiceName(serviceNames[i])
	}
	runServerWithServiceName(serviceNames[serverNum-1])
}

func newServer(serviceName string) *ServerOfHttp {
	if !config.Loaded {
		logger.Error("配置加载失败，服务启动失败")
		return nil
	}

	if serviceName == DefaultServiceName {
		if !config.GetValueBoolDefault("gole.server.http.enable", false) {
			logger.Warnf("http服务未开启")
			return nil
		}
		logger.Debugf("http服务开始初始化")
	} else {
		if !config.GetValueBoolDefault(fmt.Sprintf("gole.server.http.multi.%v.enable", serviceName), false) {
			logger.Warnf(fmt.Sprintf("http服务【%v】未开启", serviceName))
			return nil
		}
		logger.Debugf("http服务【%v】开始初始化", serviceName)
	}

	mode := config.GetValueString(fmt.Sprintf("gole.server.http.multi.%v.gin.mode", serviceName))
	if mode == "" {
		mode = config.GetValueStringDefault("gole.server.http.gin.mode", "release")
	}
	if "debug" == mode || config.GetValueStringDefault("gole.logger.level", "info") == "debug" {
		gin.SetMode(gin.DebugMode)
	} else if "test" == mode {
		gin.SetMode(gin.TestMode)
	} else if "release" == mode {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	} else {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}
	return &ServerOfHttp{ServiceName: serviceName, IRouter: gin.New()}
}

func initServerConfig(serverOfHttp *ServerOfHttp) {
	serviceName := serverOfHttp.ServiceName
	// 添加：系统内置中间件
	addSystemMiddleware(serverOfHttp)

	// 添加：系统内部路由
	addSystemRoute(serverOfHttp)

	// 发送事件：在http启动之前，业务基于这个事件可以添加中间件
	listener.PublishEvent(ServerHttpInitEvent{ServiceName: serviceName, ServerOfHttp: serverOfHttp})

	// 添加配置变更事件的监听
	listener.AddListenerWithGroup(serviceName, config.EventOfConfigChange, configChangeListenerOfPprof)
}

func graceRun(serviceName string, port int) {
	engineServer := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: GetEngineWithServiceName(serviceName)}
	go func() {
		if err := engineServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("http启动服务异常 (%v)", err)
		}
	}()

	listener.PublishEvent(event.ServerHttpAllRunFinishEvent{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-quit

	if serviceName == DefaultServiceName {
		logger.Warn("http服务端准备关闭...")
	} else {
		logger.Warn(fmt.Sprintf("http服务端【%v】准备关闭...", serviceName))
	}

	// 发送服务端关闭事件
	listener.PublishEvent(event.ServerHttpStopEvent{ServiceName: serviceName})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := engineServer.Shutdown(ctx); err != nil {
		if serviceName == DefaultServiceName {
			logger.Warnf("服务关闭异常: %v", err.Error())
		} else {
			logger.Warn(fmt.Sprintf("服务【%v】关闭异常：%v", serviceName, err.Error()))
		}
	}

	if serviceName == DefaultServiceName {
		logger.Warn("http服务端退出")
	} else {
		logger.Warn(fmt.Sprintf("http服务端【%v】退出", serviceName))
	}
}

func getMultiServiceNames() []string {
	// 获取 serviceName 的列表
	serviceNameMap := config.GetValue("gole.server.http.multi")
	if serviceNameMap == nil {
		return nil
	}
	dataMap, _ := maps.From(serviceNameMap)
	var serviceNames []string
	for _, k := range dataMap.Keys() {
		if k == "" {
			continue
		}
		if open, have := dataMap.AsDeepMap().GetBool(k + ".enable"); !have || !open {
			continue
		}
		serviceNames = append(serviceNames, k)
	}
	return serviceNames
}

func getApiPrefix(serviceName string) string {
	if serviceName == DefaultServiceName {
		apiPre := util.ISCString(config.GetValueStringDefault("gole.server.http.api.prefix", "")).Trim("/")
		if apiPre.ToString() == "" {
			return "api"
		}
		return apiPre.ToString()
	} else {
		apiPre := util.ISCString(config.GetValueStringDefault(fmt.Sprintf("gole.server.http.multi.%v.api.prefix", serviceName), "")).Trim("/")
		if apiPre.ToString() == "" {
			return "api"
		}
		return apiPre.ToString()
	}
}

func getPathAppendApiModel(serviceName, path string) string {
	apiPre := getApiPrefix(serviceName)

	pathStr := util.ISCString(path).Trim("/").ToString()
	if strings.HasPrefix(pathStr, "api") {
		return fmt.Sprintf("/%s", pathStr)
	}
	return fmt.Sprintf("/%s/%s", apiPre, pathStr)
}

func configChangeListenerOfPprof(event listener.BaseEvent) {
	ev := event.(config.EventOfChange)
	var serviceName string
	if ev.Key == "gole.server.http.pprof.enable" {
		serviceName = ev.GroupName
	} else if strings.HasPrefix(ev.Key, "gole.server.http.multi") && strings.HasSuffix(ev.Key, "pprof.enable") {
		serviceName = strings.TrimSuffix(strings.TrimPrefix(ev.Key, "gole.server.http.multi."), ".pprof.enable")
	}

	if util.ToBool(ev.Value) {
		if _, have := pprofHaveMap.Get(serviceName); !have {
			return
		}
		pprofHaveMap.Set(serviceName, true)
		pprof.Register(GetEngineWithServiceName(ev.GroupName))
	}
}
