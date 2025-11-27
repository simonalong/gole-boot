# otel
gole-boot支持opentelemetry
## 全部配置
直接添加配置即可自动埋点和信息采集（目前信息采集暂时采用的是prometheus方案）
```yaml
gole:
  opentelemetry:
    # opentelemetry（埋点和指标采集）：默认开启
    enable: true
    # opentelemetry-collector服务的地址；支持grpc和http上报，grpc示例：localhost:4317；http示例：http://localhost:4318
    exporter-url: localhost:4317
    # 服务名，默认为 gole.application.name
    service-name: demo-service
```
#### 已经接入的组件
如下的所有gole-boot组件都已经接入otel，直接按照各自组件的方法创建即可，无需额外代码控制；请不要使用非本包内的三方组件，如果使用，请自行埋点，建议及时通知@zhouzhenyong，统一封装到gole-boot

1. gin
2. http
3. gorm
4. xorm
5. go-redis
6. tdengine-orm
7. nats
8. grpc
9. tcp

# 注意断链情况！！！以下方式会造成断链的情况

### 1. 协程不要用原生go xxx()
错误用法，示例：
```go
go callDb()
```
正确用法，示例：
```go
goid.Go(func() {
    callDb()
})
```
### 2. gorm调用
在执行sql之前请先设置上下文
```go
 // 这句话必须添加，否则链会断掉
 gormClient = gormClient.WithContext(global.GetGlobalContext())

 // 执行相关语句
 gormClient.Exec("select * from device") 
```

### 3. 不要用协程池
目前短期内协程池会造成链断开，暂时先不支持，后续考虑
