package test

import (
	"github.com/simonalong/gole-boot/errorx"
	"testing"

	"github.com/gin-gonic/gin"
	httpServer "github.com/simonalong/gole-boot/server/http"
)

type DemoEntity struct {
	Name string
}

// gole.profiles.active=error
func TestServerErr1(t *testing.T) {
	httpServer.Get("err", func(c *gin.Context) (any, error) {
		var demoEntity *DemoEntity
		return demoEntity.Name, nil
	})
	httpServer.RunServer()
}

func TestServerErr2(t *testing.T) {
	httpServer.Get("err2", func(c *gin.Context) (any, error) {
		return nil, errorx.New("12", "")
	})
	httpServer.RunServer()
}

func TestServerErr3(t *testing.T) {
	httpServer.Get("data", func(c *gin.Context) (any, error) {
		return "ok", errorx.SC_SERVER_ERROR.WithDetail("异常")
	})
	httpServer.RunServer()
}
