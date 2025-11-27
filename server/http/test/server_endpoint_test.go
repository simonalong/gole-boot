package test

import (
	"testing"

	"github.com/gin-gonic/gin"
	httpServer "github.com/simonalong/gole-boot/server/http"
	"github.com/simonalong/gole/logger"
)

// gole.profiles.active=endpoint
func TestServerEndpoint(t *testing.T) {
	httpServer.Get("ok", func(c *gin.Context) (any, error) {
		InitService()
		return nil, nil
	})
	httpServer.RunServer()
}

func InitService() {
	logger.Group("demo1").Debug("demo1 debug group test")
	logger.Group("demo1").Info("demo1 info group test")
	logger.Group("demo1").Warn("demo1 warn group test")
	logger.Group("demo1").Error("demo1 error group test")

	logger.Group("demo2").Debug("demo2 debug group test")
	logger.Group("demo2").Info("demo2 info group test")
	logger.Group("demo2").Warn("demo2 warn group test")
	logger.Group("demo2").Error("demo2 error group test")

	logger.Debug("debug test")
	logger.Info("info test")
	logger.Warn("warn test")
	logger.Error("error test")
}
