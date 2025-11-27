package test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	httpServer "github.com/simonalong/gole-boot/server/http"
)

// gole.profiles.active=pprof
func TestServerOnProfileIsPprof(t *testing.T) {
	httpServer.Get("data", func(c *gin.Context) (any, error) {
		time.Sleep(4 * time.Second)
		return "value", nil
	})

	httpServer.Get("err", func(c *gin.Context) (any, error) {
		for i := 0; i < 100; i++ {
			fmt.Println(i)
		}
		return nil, errors.New("异常")
	})

	httpServer.RunServer()
}
