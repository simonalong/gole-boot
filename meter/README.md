
# meter
meter 包主要是用于统计服务的指标统计情况

## 系统指标
目前已经接入的有http、tcp、grpc、nats、orm、redis、tdengine；

说明：<br/>
- 默认开启，而每个组件包的指标默认都是跟随组件包使用而开启的，比如你使用了redis，也就是配置了gole.redis.enable，那么也就默认开启了redis对应的测量指标
- 也支持指标本身的配置
- 有些指标对应的组件包没有开关，比如http，那么这种在没有配置指标关闭的情况下就默认开启

```yaml
# application.yaml 内容
gole:
  # 测量指标
  meter:
    # 是否启用，默认：true
    enable: true
    # 输出的路径，默认为：/metrics
    path: /metrics
    http:
      # 是否启用，默认：false
      enable: false
    tcp:
      # 是否启用，默认：false
      enable: false
    grpc:
      # 是否启用，默认：false
      enable: false
    nats:
      # 是否启用，默认：false
      enable: false
    orm:
      # 是否启用，默认：false
      enable: false
    redis:
      # 是否启用，默认：false
      enable: false
    tdengine:
      # 是否启用，默认：false
      enable: false
```

## 自定义指标接入

### 计数器：只增不减
```go
import "github.com/simonalong/gole-boot/meter"

// counter计数器：只增不减：只有add和inc方法
_ = meter.AddMetric(&ginmetrics.Metric{
    Type:        ginmetrics.Counter,
    Name:        "base_demo_counter",
    Description: "这是一个描述",
    Labels:      []string{"label1"},
})

server.RegisterRoute("/counter/add", server.HmGet, func(context *gin.Context) {
    _ = meter.AddValue("base_demo_counter", []string{"label1"}, 2)
    rsp.Done(context, "ok")
})

server.RegisterRoute("/counter/inc", server.HmGet, func(context *gin.Context) {
    _ = meter.IncValue("base_demo_counter", []string{"label1"})
    rsp.Done(context, "ok")
})
```


### 测量值：可增可减
```go
import ("github.com/simonalong/gole-boot/meter")

// gauge测量：可增可减；set、inc、add三个方法
_ = meter.AddMetric(&ginmetrics.Metric{
    Type:        ginmetrics.Gauge,
    Name:        "base_demo_gauge",
    Description: "这是一个描述",
    Labels:      []string{"label1"},
})

server.RegisterRoute("/gauge/add", server.HmGet, func(context *gin.Context) {
    _ = meter.AddValue("base_demo_gauge", []string{"label1"}, 2)
    rsp.Done(context, "ok")
})

server.RegisterRoute("/gauge/inc", server.HmGet, func(context *gin.Context) {
    _ = meter.IncValue("base_demo_gauge", []string{"label1"})
    rsp.Done(context, "ok")
})

server.RegisterRoute("/gauge/set/:data", server.HmGet, func(context *gin.Context) {
    _ = meter.SetGaugeValue("base_demo_gauge", []string{"label1"}, util.ToFloat64(context.Param("data")))
    rsp.Done(context, "ok")
})
```

### 直方图
```go
import ("github.com/simonalong/gole-boot/meter")

// histogram 测量：observe一个方法
_ = meter.AddMetric(&ginmetrics.Metric{
    Type:        ginmetrics.Histogram,
    Name:        "base_demo_histogram",
    Description: "这是一个描述",
    Buckets:     []float64{0.1, 0.3, 1.2, 5, 10},
    Labels:      []string{"label1"},
})

server.RegisterRoute("/histogram/observe/:data", server.HmGet, func(context *gin.Context) {
    _ = meter.ObserveValue("base_demo_histogram", []string{"label1"}, util.ToFloat64(context.Param("data")))
    rsp.Done(context, "ok")
})
```

### 摘要
这个暂时还有点问题，后续再支持
