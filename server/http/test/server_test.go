package test

import (
	"errors"
	"github.com/simonalong/gole-boot/errorx"
	"github.com/simonalong/gole-boot/event"
	"github.com/simonalong/gole/listener"
	"testing"

	"github.com/simonalong/gole-boot/server/http/test/pojo"
	"github.com/simonalong/gole/util"

	"github.com/gin-gonic/gin"
	httpServer "github.com/simonalong/gole-boot/server/http"
	"github.com/simonalong/gole/logger"
)

// gole.profiles.active=local
func TestServerDemo(t *testing.T) {
	httpServer.Get("/api/get", GetData)
	httpServer.Get("/api/err", GetErr)
	httpServer.RunServer()
}

func GetData(c *gin.Context) (any, error) {
	logger.Info("get")
	return "hello", nil
}

func GetErr(c *gin.Context) (any, error) {
	return nil, errorx.SC_SERVER_ERROR.WithDetail("异常信息...")
}

func TestServerGet(t *testing.T) {
	httpServer.Get("/info", func(c *gin.Context) (any, error) {
		logger.Debug("debug的日志")
		logger.Info("info的日志")
		logger.Warn("warn的日志")
		logger.Error("error的日志")
		return "hello", nil
	})

	// 测试事件监听机制
	//listener.AddListener(listener.EventOfServerHttpRunFinish, func(event listener.BaseEvent) {
	//	logger.Info("应用启动完成")
	//})

	httpServer.RunServer()
}

func TestServer2(t *testing.T) {
	httpServer.Get("/test/req1", func(c *gin.Context) (any, error) {
		return "ok", nil
	})

	httpServer.Get("/test/req2", func(c *gin.Context) (any, error) {
		return "value", nil
	})

	httpServer.Get("/test/req3/:key", func(c *gin.Context) (any, error) {
		return "value", nil
	})

	httpServer.Post("/test/rsp1", func(c *gin.Context) (any, error) {
		testReq := pojo.TestReq{}
		_, _ = util.DataToEntity(c.Request.Body, &testReq)
		return testReq, nil
	})

	httpServer.Get("/test/err", func(c *gin.Context) (any, error) {
		return nil, errorx.SC_SERVER_ERROR
	})

	httpServer.Get("/test/ok", func(c *gin.Context) (any, error) {
		return "value", nil
	})

	httpServer.Get("/test/err", func(c *gin.Context) (any, error) {
		return nil, errorx.SC_SERVER_ERROR
	})

	httpServer.RunServer()
}

func TestServerError(t *testing.T) {
	httpServer.Get("data", func(c *gin.Context) (any, error) {
		return nil, errorx.SC_SERVER_ERROR
	})
	httpServer.RunServer()
}

// gole.profiles.active=otel
func TestServerOtelLogger(t *testing.T) {
	httpServer.Get("debug", func(c *gin.Context) (any, error) {
		logger.Debug("测试")
		return "value", nil
	})

	httpServer.Get("info", func(c *gin.Context) (any, error) {
		logger.Info("测试数据")
		return "value", nil
	})

	httpServer.Get("warn", func(c *gin.Context) (any, error) {
		logger.Warn("告警")
		return "value", nil
	})

	httpServer.Get("err", func(c *gin.Context) (any, error) {
		logger.Errorf("异常")
		return nil, errors.New("异常")
	})

	httpServer.Get("fatal", func(c *gin.Context) (any, error) {
		logger.Fatal("异常")
		return nil, errors.New("异常")
	})
	httpServer.RunServer()
}

func TestEvent(t *testing.T) {
	listener.AddListener(event.EventOfServerHttpRunFinish, func(event listener.BaseEvent) {
		logger.Info("应用启动完成")
	})

	httpServer.RunServer()
}
