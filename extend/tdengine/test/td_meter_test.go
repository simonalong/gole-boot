package test

import (
	"github.com/gin-gonic/gin"
	"github.com/simonalong/gole-boot/extend/tdengine"
	httpServer "github.com/simonalong/gole-boot/server/http"
	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/maps"
	"testing"
	"time"
)

// gole.profiles.active=meter
func TestMeterHttp(t *testing.T) {
	tdClient, err := tdengine.NewClient()
	if err != nil {
		logger.Errorf("tdengine连接失败：%v", err)
		return
	}

	httpServer.Get("/td/insert/ok", func(c *gin.Context) (any, error) {
		baseMap := maps.Of("ts", time.Now(), "name", "大牛市-boot", "age", 28, "address", "浙江杭州市")
		num, err := tdClient.Insert("td_china", baseMap)
		return num, err
	})

	httpServer.Get("/td/insert/err", func(c *gin.Context) (any, error) {
		baseMap := maps.Of("ts", time.Now(), "name", "大牛市-boot", "age", 28, "address", "浙江杭州市")
		num, err := tdClient.Insert("td_china_demo", baseMap)
		return num, err
	})

	httpServer.RunServer()
}
