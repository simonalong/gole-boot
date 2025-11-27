package tcp

import (
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/simonalong/gole-boot/meter"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/util"
)

var MeterTcpRequestCounter = "base_boot_tcp_request_counter"
var MeterTcpRequestErrorCounter = "base_boot_tcp_request_error_counter"
var MeterTcpConnectOpen = "base_boot_tcp_connections_open"
var MeterTcpTransferRateBytes = "base_boot_tcp_transfer_rate_bytes"

func initMeter() {
	if !config.GetValueBoolDefault("gole.meter.tcp.enable", false) {
		return
	}
	_ = meter.AddMetric(&ginmetrics.Metric{
		Type:        ginmetrics.Counter,
		Name:        MeterTcpRequestCounter,
		Description: "tcp的请求总数",
		Labels:      nil,
	})

	_ = meter.AddMetric(&ginmetrics.Metric{
		Type:        ginmetrics.Counter,
		Name:        MeterTcpRequestErrorCounter,
		Description: "tcp的请求异常总数",
		Labels:      nil,
	})

	_ = meter.AddMetric(&ginmetrics.Metric{
		Type:        ginmetrics.Gauge,
		Name:        MeterTcpConnectOpen,
		Description: "tcp的连接数",
		Labels:      nil,
	})

	_ = meter.AddMetric(&ginmetrics.Metric{
		Type:        ginmetrics.Histogram,
		Name:        MeterTcpTransferRateBytes,
		Description: "TCP传输速率",
		Buckets:     []float64{1024, 2048, 4096, 16384, 65536, 262144, 1048576},
		Labels:      nil,
	})
}

func meterIncTcpReqCounterValue() {
	if !config.GetValueBoolDefault("gole.meter.tcp.enable", false) {
		return
	}
	err := meter.IncValue(MeterTcpRequestCounter, nil)
	if err != nil {
		logger.Errorf("指标【%v】递增异常%v", MeterTcpRequestCounter, err)
	}
}

func meterIncTcpReqErrCounterValue() {
	if !config.GetValueBoolDefault("gole.meter.tcp.enable", false) {
		return
	}
	err := meter.IncValue(MeterTcpRequestErrorCounter, nil)
	if err != nil {
		logger.Errorf("指标【%v】递增异常%v", MeterTcpRequestErrorCounter, err)
	}
}

func meterUpdateTcpConnectValue() {
	if !config.GetValueBoolDefault("gole.meter.tcp.enable", false) {
		return
	}

	var connect int
	if rootHandler != nil {
		connect = rootHandler.eng.CountConnections()
	}
	err := meter.SetGaugeValue(MeterTcpConnectOpen, nil, util.ToFloat64(connect))
	if err != nil {
		logger.Errorf("指标【%v】设置异常%v", MeterTcpConnectOpen, err)
	}
}

func meterObserveByteValue(byteLength int) {
	if !config.GetValueBoolDefault("gole.meter.tcp.enable", false) {
		return
	}
	err := meter.ObserveValue(MeterTcpTransferRateBytes, nil, util.ToFloat64(byteLength))
	if err != nil {
		logger.Errorf("指标【%v】设置异常%v", MeterTcpTransferRateBytes, err)
	}
}
