package test

import (
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/simonalong/gole-boot/meter"
	httpServer "github.com/simonalong/gole-boot/server/http"
	"github.com/simonalong/gole/util"
	"testing"
)

func TestGinProm1(t *testing.T) {
	// counter计数器：只增不减：只有add和inc方法
	_ = meter.AddCounter(&ginmetrics.Metric{
		Name:        "base_demo_counter",
		Description: "这是一个描述",
		Labels:      []string{"label1"},
	})
	httpServer.AddRoute(httpServer.HmGet, "/counter/add", func(context *gin.Context) (any, error) {
		_ = meter.AddValue("base_demo_counter", []string{"label1"}, 2)
		return "ok", nil
	})

	httpServer.AddRoute(httpServer.HmGet, "/counter/inc", func(context *gin.Context) (any, error) {
		_ = meter.IncValue("base_demo_counter", []string{"label1"})
		return "ok", nil
	})

	// gauge测量：可增可减；set、inc、add三个方法
	_ = meter.AddGauge(&ginmetrics.Metric{
		Name:        "base_demo_gauge",
		Description: "这是一个描述",
		Labels:      []string{"label1"},
	})
	httpServer.AddRoute(httpServer.HmGet, "/gauge/add", func(context *gin.Context) (any, error) {
		_ = meter.AddValue("base_demo_gauge", []string{"label1"}, 2)
		return "ok", nil
	})

	httpServer.AddRoute(httpServer.HmGet, "/gauge/inc", func(context *gin.Context) (any, error) {
		_ = meter.IncValue("base_demo_gauge", []string{"label1"})
		return "ok", nil
	})

	httpServer.AddRoute(httpServer.HmGet, "/gauge/set/:data", func(context *gin.Context) (any, error) {
		_ = meter.SetGaugeValue("base_demo_gauge", []string{"label1"}, util.ToFloat64(context.Param("data")))
		return "ok", nil
	})

	// histogram 测量：observe一个方法
	_ = meter.AddHistogram(&ginmetrics.Metric{
		Name:        "base_demo_histogram",
		Description: "这是一个描述",
		Buckets:     []float64{0.1, 0.3, 1.2, 5, 10},
		Labels:      []string{"label1"},
	})
	httpServer.AddRoute(httpServer.HmGet, "/histogram/observe/:data", func(context *gin.Context) (any, error) {
		_ = meter.ObserveValue("base_demo_histogram", []string{"label1"}, util.ToFloat64(context.Param("data")))
		return "ok", nil
	})

	// summary 测量：observe一个方法；目前这个有点不好理解，暂时先不封装，后续再说
	//_ = metric.AddMetric(&ginmetrics.Metric{
	//	Type:        ginmetrics.Summary,
	//	Name:        "base_demo_summary",
	//	Description: "这是一个描述",
	//	Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	//	Labels:      []string{"label1"},
	//})
	//server.AddGinRoute("/summary/observe/:data", server.HmGet, func(context *gin.Context) {
	//	_ = metric.ObserveValue("base_demo_summary", []string{"label1"}, util.ToFloat64(context.Param("data")))
	//	rsp.Done(context, "ok")
	//})
	httpServer.RunServer()
}
