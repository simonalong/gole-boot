package meter

import (
	"errors"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/simonalong/gole/logger"
)

func AddCounter(metric *ginmetrics.Metric) error {
	metric.Type = ginmetrics.Counter
	return AddMetric(metric)
}

func AddGauge(metric *ginmetrics.Metric) error {
	metric.Type = ginmetrics.Gauge
	return AddMetric(metric)
}

func AddHistogram(metric *ginmetrics.Metric) error {
	metric.Type = ginmetrics.Histogram
	return AddMetric(metric)
}

func AddSummary(metric *ginmetrics.Metric) error {
	metric.Type = ginmetrics.Summary
	return AddMetric(metric)
}

func AddMetric(metric *ginmetrics.Metric) error {
	err := ginmetrics.GetMonitor().AddMetric(metric)
	if err != nil {
		logger.Errorf("添加metric【%v】失败：%v", metric.Name, err)
		return err
	}
	return nil
}

func GetMetric(metricName string) *ginmetrics.Metric {
	return ginmetrics.GetMonitor().GetMetric(metricName)
}

func AddValue(metricName string, labelValues []string, value float64) error {
	metricObj := GetMetric(metricName)
	if metricObj == nil {
		return errors.New("metric not found")
	}
	return metricObj.Add(labelValues, value)
}

func IncValue(metricName string, labelValues []string) error {
	metricObj := GetMetric(metricName)
	if metricObj == nil {
		return errors.New("metric not found")
	}
	return metricObj.Inc(labelValues)
}

func SetGaugeValue(metricName string, labelValues []string, value float64) error {
	metricObj := GetMetric(metricName)
	if metricObj == nil {
		return errors.New("metric not found")
	}
	return metricObj.SetGaugeValue(labelValues, value)
}

func ObserveValue(metricName string, labelValues []string, value float64) error {
	metricObj := GetMetric(metricName)
	if metricObj == nil {
		return errors.New("metric not found")
	}
	return metricObj.Observe(labelValues, value)
}
