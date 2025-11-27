package test

import (
	"context"
	"github.com/gin-gonic/gin"
	goredis "github.com/go-redis/redis/v8"
	redis2 "github.com/simonalong/gole-boot/extend/redis"
	"github.com/simonalong/gole-boot/server/http"
	"github.com/simonalong/gole/logger"
	baseTime "github.com/simonalong/gole/time"
	"testing"
	"time"
)

func init() {
	// 客户端获取
	_rdb, err := redis2.GetClient()
	if err != nil {
		return
	}
	rdb = _rdb
}

// gole.profiles.active=meter-http
func TestMeterHttp(t *testing.T) {
	var rdb goredis.UniversalClient
	_rdb, err := redis2.GetClient()
	if err != nil {
		logger.Errorf("redis客户端创建失败：%v", err)
		return
	}
	rdb = _rdb

	http.Get("/redis/get", func(c *gin.Context) (any, error) {
		rspData := rdb.Get(context.Background(), "test_key")
		return rspData.Val(), nil
	})

	http.Get("/redis/set", func(c *gin.Context) (any, error) {
		rdb.Set(context.Background(), "test_key", baseTime.TimeToStringYmdHmsS(time.Now()), time.Hour)
		return "ok", nil
	})

	http.RunServer()
}
