package test

import (
	"github.com/gin-gonic/gin"
	"testing"

	httpServer "github.com/simonalong/gole-boot/server/http"
)

// gole.profiles.active=local
func TestServerGroup(t *testing.T) {
	routeGroup1 := httpServer.Group("/api/v1")
	{
		routeGroup1.Get("/path/xxx1", groupH1)
		routeGroup1.Get("/path/xxx2", groupH1)
		routeGroup1.Get("/path/xxx3", groupH1)
		routeGroup1.Get("/path/xxx4", groupH1)
	}

	routeGroup2 := httpServer.Group("/api/v2")
	{
		routeGroup2.Get("/path/xxx1", groupH1)
		routeGroup2.Get("/path/xxx2", groupH1)
		routeGroup2.Get("/path/xxx3", groupH1)
		routeGroup2.Get("/path/xxx4", groupH1)
	}
	httpServer.RunServer()
}

func groupH1(c *gin.Context) (any, error) {
	return 1, nil
}
