package original

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestGin(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	route := gin.Default()

	route.GET("/api/demo/data", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	s := &http.Server{
		Addr:              ":8080",
		Handler:           route,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
