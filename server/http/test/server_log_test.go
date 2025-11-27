package test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/simonalong/gole-boot/errorx"
	httpServer "github.com/simonalong/gole-boot/server/http"
)

// gole.profiles.active=endpoint
func TestLoggerServer(t *testing.T) {
	httpServer.Get("/api/get", GetTest)
	httpServer.Get("/api/get/:test", GetTestParam)
	httpServer.Get("/api/get/query", GetTestQuery)
	httpServer.Post("/api/post", PostTestBody)
	httpServer.RunServer()
}

func GetTest(c *gin.Context) (any, error) {
	return "ok", nil
}

func GetTestParam(c *gin.Context) (any, error) {
	return c.Param("test"), nil
}

func GetTestQuery(c *gin.Context) (any, error) {
	return c.Query("test"), nil
}

func PostTestBody(c *gin.Context) (any, error) {
	var dataReq PostBodyReqTest
	if err := c.ShouldBindJSON(&dataReq); err != nil {
		return nil, errorx.SC_BAD_REQUEST.WithError(err)
	}
	return dataReq, nil
}

type PostBodyReqTest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
