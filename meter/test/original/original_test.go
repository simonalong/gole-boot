package original

import (
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpServer "github.com/simonalong/gole-boot/server/http"
	"testing"
)

func TestProm(t *testing.T) {
	reg := prometheus.NewRegistry()

	// Add go runtime metrics and process collectors.
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	httpServer.AddRouteGinHandler("/metrics2", httpServer.HmGet, func(c *gin.Context) {
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}).ServeHTTP(c.Writer, c.Request)
	})

	httpServer.AddRoute("/data", httpServer.HmGet, func(context *gin.Context) (any, error) {
		return "ok", nil
	})

	httpServer.RunServer()
}

func TestGinMetric(t *testing.T) {
	r := gin.Default()

	// get global Monitor object
	m := ginmetrics.GetMonitor()

	// +optional set metric path, default /debug/metrics
	m.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(10)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})

	// set middleware for gin
	m.Use(r)

	r.GET("/product/:id", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]string{
			"productId": ctx.Param("id"),
		})
	})

	_ = r.Run()
}
