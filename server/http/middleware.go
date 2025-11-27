package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/simonalong/gole-boot/errorx"
	baseOrm "github.com/simonalong/gole-boot/extend/orm"
	baseOtel "github.com/simonalong/gole-boot/otel"
	"github.com/simonalong/gole-boot/server/http/rsp"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/logger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
	"time"
)

func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer func() {
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				c.Abort()
			}
			cancel()
		}()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func addSystemMiddleware(serverOfHttp *ServerOfHttp) {
	// 接入prometheus
	if config.GetValueBoolDefault("gole.meter.enable", true) {
		registerMeterHandler(serverOfHttp)
	}

	// 注册pprof
	if config.GetValueBoolDefault("gole.debug.enable", true) {
		serviceName := serverOfHttp.ServiceName
		if serviceName == DefaultServiceName {
			if config.GetValueBoolDefault("gole.server.http.pprof.enable", false) {
				pprofHaveMap.Set(serviceName, true)
				pprof.Register(serverOfHttp.IRouter)
			}
		} else {
			if config.GetValueBoolDefault(fmt.Sprintf("gole.server.http.multi.%v.pprof.enable", serviceName), false) {
				pprofHaveMap.Set(serviceName, true)
				pprof.Register(serverOfHttp.IRouter)
			}
		}
	}

	// 注册跨域配置
	if config.GetValueBoolDefault("gole.server.http.cors.enable", true) {
		serverOfHttp.AddMiddleware(cors())
	}

	serverOfHttp.AddMiddleware(gin.Recovery())
	serverOfHttp.AddMiddleware(rsp.ResponseHandler())

	// 开启 405 Method Not Allowed 的处理程序，默认是关闭，关闭时将请求委托给 404 NotFound 处理程序
	serverOfHttp.GetEngine().HandleMethodNotAllowed = true
	// 处理 405
	serverOfHttp.GetEngine().NoRoute(noRoute)
	// 处理 405
	serverOfHttp.GetEngine().NoMethod(noMethod)

	// 支持opentelemetry埋点
	if config.GetValueBoolDefault("gole.opentelemetry.enable", true) {
		serverOfHttp.AddMiddleware(otelgin.Middleware(baseOtel.CfgOfOtel.ServiceName))
		serverOfHttp.AddMiddleware(baseOtel.GlobalContextLoad())
		serverOfHttp.AddMiddleware(baseOrm.GlobalGormContextLoad())
	}

	// 业务全局异常处理
	serverOfHttp.AddMiddleware(errHandler())
}

// NotFound
func noRoute(ctx *gin.Context) {
	rsp.Done(ctx, nil, errorx.SC_NOT_FOUND)
}

// MethodNotAllowed
func noMethod(ctx *gin.Context) {
	rsp.Done(ctx, nil, errorx.SC_METHOD_NOT_ALLOWED)
}

func errHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("服务端某个地方异常: %v", err)
				switch err.(type) {
				case *errorx.BaseError:
					rsp.Done(c, nil, err.(*errorx.BaseError))
				default:
					rsp.Done(c, nil, errorx.SC_SERVER_ERROR.WithError(err.(error)))
				}
			}
		}()
		c.Next()
	}
}

// cors 配置允许跨域请求
func cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 处理跨域请求,支持options访问
		if origin := ctx.GetHeader("Origin"); origin != "" {
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE, PATCH")
			ctx.Header("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Detail, Access-Control-Allow-Headers")
			ctx.Header("Access-Control-Max-Age", "172800")
			ctx.Header("Access-Control-Allow-Credentials", "true")
		}

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

func registerMeterHandler(serverOfHttp *ServerOfHttp) {
	m := ginmetrics.GetMonitor()
	m.SetMetricPath(config.GetValueStringDefault("gole.meter.path", "/metrics"))
	m.SetSlowTime(10)
	m.Use(serverOfHttp.IRouter)
}
