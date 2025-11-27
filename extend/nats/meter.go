package nats

import (
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/simonalong/gole-boot/meter"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/logger"
)

var MeterNatsClientRequestOkCounter = "base_boot_nats_client_request_ok_counter"
var MeterNatsClientRequestErrCounter = "base_boot_nats_client_request_err_counter"

//var MeterNatsRateBytes = "base_boot_nats_client_rate_bytes"

var MeterNatsJsClientRequestOkCounter = "base_boot_nats_js_client_request_ok_counter"
var MeterNatsJsClientRequestErrCounter = "base_boot_nats_js_client_request_err_counter"

//var MeterNatsJsRateBytes = "base_boot_nats_js_client_rate_bytes"

var MeterNatsServerHandleOkCounter = "base_boot_nats_server_handle_ok_counter"
var MeterNatsServerHandleErrCounter = "base_boot_nats_server_handle_err_counter"

var MeterNatsJsServerHandleOkCounter = "base_boot_nats_js_server_handle_ok_counter"
var MeterNatsJsServerHandleErrCounter = "base_boot_nats_js_server_handle_err_counter"

func initMeterOfNatsClient() {
	if !config.GetValueBoolDefault("gole.meter.nats.enable", false) {
		return
	}
	_ = meter.AddCounter(&ginmetrics.Metric{
		Name:        MeterNatsClientRequestOkCounter,
		Description: "nats发布消息成功总数",
		Labels:      []string{"subject", "method", "reply"},
	})

	_ = meter.AddCounter(&ginmetrics.Metric{
		Name:        MeterNatsClientRequestErrCounter,
		Description: "nats发布消息失败总数",
		Labels:      []string{"subject", "method", "reply"},
	})

	// 这个subject会很多，如果在加上这个，那么这个指标就太多了，因此这边这个注释掉，暂时先不用
	//_ = meter.AddHistogram(&ginmetrics.Metric{
	//	Name:        MeterNatsRateBytes,
	//	Description: "nats发布数据的大小",
	//	Buckets:     []float64{1024, 2048, 4096, 16384, 65536, 262144, 1048576},
	//	Labels:      []string{"subject", "method", "reply"},
	//})

	_ = meter.AddCounter(&ginmetrics.Metric{
		Name:        MeterNatsServerHandleOkCounter,
		Description: "nats消息成功处理总数",
		Labels:      []string{"subject"},
	})
	_ = meter.AddCounter(&ginmetrics.Metric{
		Name:        MeterNatsServerHandleErrCounter,
		Description: "nats消息失败处理总数",
		Labels:      []string{"subject"},
	})
}

func initMeterOfNatsJsClient() {
	if !config.GetValueBoolDefault("gole.meter.nats.enable", false) {
		return
	}
	_ = meter.AddCounter(&ginmetrics.Metric{
		Name:        MeterNatsJsClientRequestOkCounter,
		Description: "nats jetstream 发布消息成功总数",
		Labels:      []string{"subject", "method", "reply"},
	})

	_ = meter.AddCounter(&ginmetrics.Metric{
		Name:        MeterNatsJsClientRequestErrCounter,
		Description: "nats jetstream 发布消息失败总数",
		Labels:      []string{"subject", "method", "reply"},
	})

	// 这个subject会很多，如果在加上这个，那么这个指标就太多了，因此这边这个注释掉，暂时先不用
	//_ = meter.AddHistogram(&ginmetrics.Metric{
	//	Name:        MeterNatsJsRateBytes,
	//	Description: "nats发布数据的大小",
	//	Buckets:     []float64{1024, 2048, 4096, 16384, 65536, 262144, 1048576},
	//	Labels:      []string{"subject", "method", "reply"},
	//})

	_ = meter.AddCounter(&ginmetrics.Metric{
		Name:        MeterNatsJsServerHandleOkCounter,
		Description: "nats jetstream 消息成功处理总数",
		Labels:      []string{"subject"},
	})
	_ = meter.AddCounter(&ginmetrics.Metric{
		Name:        MeterNatsJsServerHandleErrCounter,
		Description: "nats jetstream 消息失败处理总数",
		Labels:      []string{"subject"},
	})
}

func meterIncOkCounterValueOfNats(subject, method, reply string) {
	if !config.GetValueBoolDefault("gole.meter.nats.enable", false) {
		return
	}
	err := meter.IncValue(MeterNatsClientRequestOkCounter, []string{subject, method, reply})
	if err != nil {
		logger.Errorf("指标【%v】递增异常：%v", MeterNatsClientRequestOkCounter, err)
	}
}

func meterIncOkCounterValueOfNatsJs(subject, method, reply string) {
	if !config.GetValueBoolDefault("gole.meter.nats.enable", false) {
		return
	}
	err := meter.IncValue(MeterNatsJsClientRequestOkCounter, []string{subject, method, reply})
	if err != nil {
		logger.Errorf("指标【%v】递增异常：%v", MeterNatsJsClientRequestOkCounter, err)
	}
}

func meterIncErrCounterValueOfNats(subject, method, reply string) {
	if !config.GetValueBoolDefault("gole.meter.nats.enable", false) {
		return
	}
	err := meter.IncValue(MeterNatsClientRequestErrCounter, []string{subject, method, reply})
	if err != nil {
		logger.Errorf("指标【%v】递增异常：%v", MeterNatsClientRequestErrCounter, err)
	}
}

func meterIncErrCounterValueOfNatsJs(subject, method, reply string) {
	if !config.GetValueBoolDefault("gole.meter.nats.enable", false) {
		return
	}
	err := meter.IncValue(MeterNatsJsClientRequestErrCounter, []string{subject, method, reply})
	if err != nil {
		logger.Errorf("指标【%v】递增异常：%v", MeterNatsJsClientRequestErrCounter, err)
	}
}

//func meterObserveByteValueOfNats(subject, method, reply string, byteLength int) {
//	if !config.GetValueBoolDefault("gole.meter.nats.enable", false) {
//		return
//	}
//	err := meter.ObserveValue(MeterNatsRateBytes, []string{subject, method, reply}, util.ToFloat64(byteLength))
//	if err != nil {
//		logger.Errorf("指标【%v】设置异常：%v", MeterNatsRateBytes, err)
//	}
//}

//func meterObserveByteValueOfNatsJs(subject, method, reply string, byteLength int) {
//	if !config.GetValueBoolDefault("gole.meter.nats.enable", false) {
//		return
//	}
//	err := meter.ObserveValue(MeterNatsJsRateBytes, []string{subject, method, reply}, util.ToFloat64(byteLength))
//	if err != nil {
//		logger.Errorf("指标【%v】设置异常：%v", MeterNatsJsRateBytes, err)
//	}
//}

func MeterIncOkCounterValueOfNatsServer(subject string) {
	if !config.GetValueBoolDefault("gole.meter.nats.enable", false) {
		return
	}
	err := meter.IncValue(MeterNatsServerHandleOkCounter, []string{subject})
	if err != nil {
		logger.Errorf("指标【%v】递增异常：%v", MeterNatsServerHandleOkCounter, err)
	}
}

func MeterIncErrCounterValueOfNatsServer(subject string) {
	if !config.GetValueBoolDefault("gole.meter.nats.enable", false) {
		return
	}
	err := meter.IncValue(MeterNatsServerHandleErrCounter, []string{subject})
	if err != nil {
		logger.Errorf("指标【%v】递增异常：%v", MeterNatsServerHandleErrCounter, err)
	}
}

func MeterIncOkCounterValueOfNatsJsServer(subject string) {
	if !config.GetValueBoolDefault("gole.meter.nats.enable", false) {
		return
	}
	err := meter.IncValue(MeterNatsJsServerHandleOkCounter, []string{subject})
	if err != nil {
		logger.Errorf("指标【%v】递增异常：%v", MeterNatsJsServerHandleOkCounter, err)
	}
}

func MeterIncErrCounterValueOfNatsJsServer(subject string) {
	if !config.GetValueBoolDefault("gole.meter.nats.enable", false) {
		return
	}
	err := meter.IncValue(MeterNatsJsServerHandleErrCounter, []string{subject})
	if err != nil {
		logger.Errorf("指标【%v】递增异常：%v", MeterNatsJsServerHandleErrCounter, err)
	}
}
