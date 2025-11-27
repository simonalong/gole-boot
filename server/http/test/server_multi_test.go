package test

import (
	"github.com/gin-gonic/gin"
	"github.com/simonalong/gole/logger"
	"testing"

	httpServer "github.com/simonalong/gole-boot/server/http"
)

// gole.profiles.active=multi-http
// gole.profiles.active=multi-http-api
// gole.profiles.active=multi-http-pprof
func TestServerMulti(t *testing.T) {
	serverName1 := httpServer.Server("server-name1")
	serverName1.Get("/get1", multiHandle1)
	serverName1.Get("/get2", multiHandle2)
	serverName1.Get("/get3", multiHandle3)

	serverName2 := httpServer.Server("server-name2")
	serverName2.Get("/get1", multiHandle1)
	serverName2.Get("/get2", multiHandle2)
	serverName2.Get("/get3", multiHandle3)

	serverName3 := httpServer.Server("server-name3")
	serverName3.Get("/get1", multiHandle1)
	serverName3.Get("/get2", multiHandle2)
	serverName3.Get("/get3", multiHandle3)

	httpServer.RunServer()
}

func multiHandle1(c *gin.Context) (any, error) {
	logger.Info("1")
	return 1, nil
}

func multiHandle2(c *gin.Context) (any, error) {
	logger.Info("2")
	return 2, nil
}

func multiHandle3(c *gin.Context) (any, error) {
	logger.Info("3")
	return 3, nil
}
