package test

import (
	"context"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	baseHttp "github.com/simonalong/gole-boot/client/http"
	httpServer "github.com/simonalong/gole-boot/server/http"
)

type DemoHttpHook struct {
}

func (*DemoHttpHook) Before(ctx context.Context, req *http.Request) context.Context {
	return ctx
}

func (*DemoHttpHook) After(ctx context.Context, rsp *http.Response, rspCode int, rspData any, err error) {

}

// gole.profiles.active=client
func TestGetSimple(t *testing.T) {
	httpServer.Get("call/ok", func(c *gin.Context) (any, error) {
		_, _, _, _ = baseHttp.GetSimple("http://localhost:8080/api/ok")
		return "ok", nil
	})
	httpServer.Get("call/err", func(c *gin.Context) (any, error) {
		_, _, _, _ = baseHttp.GetSimple("http://localhost:8080/api/err")
		return "ok", nil
	})
	httpServer.RunServer()
}
