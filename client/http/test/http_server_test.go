package test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/simonalong/gole-boot/errorx"
	"github.com/simonalong/gole-boot/server/http"
)

// gole.profiles.active=server
func TestHttpServer(t *testing.T) {
	http.Get("ok", func(c *gin.Context) (any, error) {
		return "ok", nil
	})
	http.Get("err", func(c *gin.Context) (any, error) {
		return "", &errorx.BaseError{Code: "500", Msg: "err"}
	})
	http.RunServer()
}
