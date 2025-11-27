package orm

import (
	"context"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/simonalong/gole-boot/meter"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/util"
	"time"
)

var MeterGormSlowSql = "base_boot_gorm_slow_sql"
var MeterGormMaxOpenConnections = "base_boot_gorm_db_max_open_connections"
var MeterGormOpenConnections = "base_boot_gorm_db_open_connections"
var MeterGormInUseConnections = "base_boot_gorm_db_in_use_connections"
var MeterGormIdleConnections = "base_boot_gorm_db_idle_connections"
var MeterGormWaitCountConnections = "base_boot_gorm_wait_count_connections"
var MeterGormWaitDurationConnections = "base_boot_gorm_wait_duration_connections"
var MeterGormMaxIdleClosedConnections = "base_boot_gorm_max_idle_closed_connections"
var MeterGormMaxLifetimeClosedConnections = "base_boot_gorm_max_lifetime_closed_connections"
var MeterGormMaxIdleTimeClosedConnections = "base_boot_gorm_max_idle_time_closed_connections"

func init() {
	AddGormHook(&MeterGormHook{})
}

func initMeterOfGorm() {
	if !config.GetValueBoolDefault("gole.meter.orm.enable", false) {
		return
	}

	_ = meter.AddHistogram(&ginmetrics.Metric{
		Name:        MeterGormSlowSql,
		Description: "慢查询sql统计",
		Buckets:     []float64{1, 5, 10, 30, 60, 300, 1800},
		Labels:      []string{"sql"},
	})

	_ = meter.AddGauge(&ginmetrics.Metric{
		Name:        MeterGormMaxOpenConnections,
		Description: "数据库的最大打开连接数",
		Labels:      []string{"GormClient"},
	})

	_ = meter.AddGauge(&ginmetrics.Metric{
		Name:        MeterGormOpenConnections,
		Description: "正在使用和空闲的已建立连接数",
		Labels:      []string{"GormClient"},
	})

	_ = meter.AddGauge(&ginmetrics.Metric{
		Name:        MeterGormInUseConnections,
		Description: "当前正在使用的连接数",
		Labels:      []string{"GormClient"},
	})

	_ = meter.AddGauge(&ginmetrics.Metric{
		Name:        MeterGormIdleConnections,
		Description: "空闲连接数",
		Labels:      []string{"GormClient"},
	})

	_ = meter.AddGauge(&ginmetrics.Metric{
		Name:        MeterGormWaitCountConnections,
		Description: "等待的连接总数",
		Labels:      []string{"GormClient"},
	})

	_ = meter.AddGauge(&ginmetrics.Metric{
		Name:        MeterGormWaitDurationConnections,
		Description: "等待新连接的总阻塞时间",
		Labels:      []string{"GormClient"},
	})

	_ = meter.AddGauge(&ginmetrics.Metric{
		Name:        MeterGormMaxIdleClosedConnections,
		Description: "由于SetMaxIdleCon而关闭的连接总数",
		Labels:      []string{"GormClient"},
	})

	_ = meter.AddGauge(&ginmetrics.Metric{
		Name:        MeterGormMaxLifetimeClosedConnections,
		Description: "由于SetConnMaxLifetime而关闭的连接总数",
		Labels:      []string{"GormClient"},
	})

	_ = meter.AddGauge(&ginmetrics.Metric{
		Name:        MeterGormMaxIdleTimeClosedConnections,
		Description: "由于SetConnMaxIdleTime而关闭的连接总数",
		Labels:      []string{"GormClient"},
	})
}

func setMeterValueOfGorm(meterName, db string, value interface{}) {
	_ = meter.SetGaugeValue(meterName, []string{db}, util.ToFloat64(value))
}

func observeValueOfGorm(meterName, sql string, value interface{}) {
	err := meter.ObserveValue(meterName, []string{sql}, util.ToFloat64(value))
	if err != nil {
		logger.Errorf("设置指标值异常：%v", err)
	}
}

func setSlowSql(sql string, value interface{}) {
	observeValueOfGorm(MeterGormSlowSql, sql, value)
}

type MeterGormHook struct {
}

func (*MeterGormHook) Before(ctx context.Context, driverName string, parameters map[string]any) (context.Context, error) {
	nextCtx := context.WithValue(ctx, "start", time.Now())
	return nextCtx, nil
}

func (*MeterGormHook) After(ctx context.Context, driverName string, parameters map[string]any) (context.Context, error) {
	startTime := ctx.Value("start").(time.Time)
	sql := parameters["query"].(string)

	setSlowSql(sql, time.Now().Sub(startTime).Seconds())
	return ctx, nil
}

func (*MeterGormHook) Err(ctx context.Context, driverName string, err error, parameters map[string]any) error {
	return nil
}
