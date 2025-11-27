# job
分布式定时任务服务（cbb-mid-srv-job）的客户端
# 功能
- 本地调度：功能简单，就是分布式锁+定时器的组合，不依赖分布式任务cbb-mid-srv-job服务端
- 服务端调度：功能更加复杂和精细，是分布式任务cbb-mid-srv-job的客户端

说明：
两种模式如何选择？两种模式并非互斥，而是可以同时使用。
1. 本地调度：适合于简单的定时功能，不需要界面化处理，但是又需要保证业务集群里面只有一个节点执行的场景。
2. 服务端调度：适合于业务较为复杂和需要精细管理的定时任务，需要界面化处理，而且又需要支持业务集群里面所有节点都执行的场景。

## 本地调度
### 配置
必须配置redis，否则会报错；这边使用redis主要是使用redis的分布式锁功能
```yaml
gole:
  application:
    name: cbb-base-job-demo
    version: 1.0.0
  redis:
    enable: true
    standalone:
      addr: localhost:6379
```

### 代码
提供两个接口

```go
// ScheduleCron 定时任务，支持cron表达式
// bizName 业务名称，用于业务唯一识别
// cron 表达式
// fun 业务函数
ScheduleCron(bizName, cron string, bizFun func()) {}

// ScheduleFixRate 每隔一段时间执行一次，无论前次任务是否完成
// bizName 业务名称，用于业务唯一识别
// duration 固定频率，最小单位为秒
// fun 业务函数
ScheduleFixRate(bizName string, duration time.Duration, bizFun func()) {}
```
#### 示例：
```go
package test

import (
    "github.com/simonalong/gole-boot/extend/job"
    "github.com/simonalong/gole/logger"
    "testing"
    "time"
)

func TestScheduleCron(t *testing.T) {
    job.ScheduleCron("cbb-job-demo:TestScheduleCron", "* * * * * ?", job1)

    time.Sleep(12 * time.Hour)
}

func TestScheduleFixRate(t *testing.T) {
    job.ScheduleFixRate("cbb-job-demo:TestScheduleFixRate", 2*time.Second, job1)

    time.Sleep(12 * time.Hour)
}

func job1() {
    logger.Info("job1")
    time.Sleep(500 * time.Millisecond)
}
```

## 服务端调度
这个其实就是分布式任务调度中心的客户端了，与cbb-mid-srv-job进行交互，实现分布式任务调度。
### 配置
```yaml
# application.yaml
gole:
  application:
    name: cbb-base-job-demo
    version: 1.0.0
  server:
    http:
      enable: true
      port: 8080
  job:
    # 是否启用，默认：false
    enable: true
    # 分布式任务服务地址，默认：http://cbb-mid-srv-job:18080
    server-address: http://cbb-mid-srv-job:18080
    # 执行器名称（用于在分布式定时任务注册），默认：${gole.application.name}
    executor-name: cbb-base-job-demo
    # 执行器所在的业务服务名和端口（用于在分布式定时任务注册），默认：http://${gole.application.name}:${gole.server.http.port}
    executor-address: http://cbb-base-job-demo:8080
```

### 代码

提供一个任务处理接口
```go
// 添加任务处理器
AddJobHandler(jobHandlerName string, task TaskFunc) {}
```

#### 示例：
```go
package test

import (
    "context"
    "github.com/simonalong/gole-boot/extend/job"
    "github.com/simonalong/gole-boot/server"
    "testing"
)

func TestJobHandler(t *testing.T) {
    // 添加任务处理器
    job.AddJobHandler("testJob1", testJob1)
    job.AddJobHandler("testJob2", testJob2)
    job.AddJobHandler("testJob3", testJob3)

    server.Run()
}

func testJob1(cxt context.Context, param *job.RunReq) string {
    return "success"
}
func testJob2(cxt context.Context, param *job.RunReq) string {
    return "success"
}
func testJob3(cxt context.Context, param *job.RunReq) string {
    return "success"
}
```
