package http

import (
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/simonalong/gole-boot/meter"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/util"
)

var MeterHttpRequestsCounter = "base_boot_http_client_requests_counter"

func InitHttpMeter() {
	if !config.GetValueBoolDefault("gole.meter.http.enable", false) {
		return
	}
	_ = meter.AddCounter(&ginmetrics.Metric{
		Name:        MeterHttpRequestsCounter,
		Description: "向外部发起的http请求总个数",
		Labels:      []string{"method", "url", "status_code"},
	})
}

func incMeterValue(method, url string, statusCode int) {
	if !config.GetValueBoolDefault("gole.meter.http.enable", false) {
		return
	}
	err := meter.IncValue(MeterHttpRequestsCounter, []string{method, url, util.ToString(statusCode)})
	if err != nil {
		logger.Errorf("指标递增失败：%v", err)
	}
}
