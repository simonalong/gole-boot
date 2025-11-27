package tdengine

import (
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/simonalong/gole-boot/meter"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/util"
	"github.com/simonalong/tdorm"
)

var MeterTdRequestOkCounter = "base_boot_tdengine_request_ok_counter"
var MeterTdRequestErrCounter = "base_boot_tdengine_request_err_counter"

var MeterTdSlowSqlHistogram = "base_boot_tdengine_slow_sql_histogram"

func initMeter() {
	if !config.GetValueBoolDefault("gole.meter.tdengine.enable", false) {
		return
	}
	_ = meter.AddCounter(&ginmetrics.Metric{
		Name:        MeterTdRequestOkCounter,
		Description: "tdengine处理正常的总数",
		Labels:      []string{"db", "RunType"},
	})

	_ = meter.AddCounter(&ginmetrics.Metric{
		Name:        MeterTdRequestErrCounter,
		Description: "tdengine处理异常的总数",
		Labels:      []string{"db", "RunType"},
	})

	_ = meter.AddHistogram(&ginmetrics.Metric{
		Name:        MeterTdSlowSqlHistogram,
		Description: "tdengine慢查询的统计（单位秒）",
		Buckets:     []float64{1, 5, 10, 30, 60, 300, 1800},
		Labels:      []string{"db", "sql"},
	})
}

func incMeterOkValue(dbName, runType string) {
	err := meter.IncValue(MeterTdRequestOkCounter, []string{dbName, runType})
	if err != nil {
		logger.Errorf("递增指标异常：%v", err)
	}
}

func incMeterErrValue(dbName, runType string) {
	err := meter.IncValue(MeterTdRequestErrCounter, []string{dbName, runType})
	if err != nil {
		logger.Errorf("递增指标异常：%v", err)
	}
}

func observeMeterValue(dbName, runType string, value interface{}) {
	err := meter.ObserveValue(MeterTdSlowSqlHistogram, []string{dbName, runType}, util.ToFloat64(value))
	if err != nil {
		logger.Errorf("递增指标异常：%v", err)
	}
}

type MeterTdHook struct {
}

func (tdHook *MeterTdHook) Before(thc *tdorm.TdHookContext) (*tdorm.TdHookContext, error) {
	return thc, nil
}

func (tdHook *MeterTdHook) After(thc *tdorm.TdHookContext) error {
	incMeterOkValue(thc.DbName, thc.RunType)

	observeMeterValue(thc.DbName, thc.RunType, thc.ExecuteTime.Seconds())

	if thc.Err != nil {
		incMeterErrValue(thc.DbName, thc.RunType)
		return thc.Err
	}
	return nil
}
