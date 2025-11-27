# gole-boot

gole-boot 框架是从Java向golang转型的过程中总结的一套至简化的工具框架，借鉴spring-boot的思想，遵从大道至简原则，脱胎与gole项目，旨在将三方框架集成放到这里

## 下载
```shell
go get github.com/simonalong/gole-boot
```
### 提示：
1. 更新相关依赖
```shell
go mod tidy
```
2. 下载指定版本
```shell
go get github.com/simonalong/gole-boot@<version>
```
目前最新版本：v1.2.0

## 快速入门
gole-boot定位是工具框架，包含各种各样的工具，并对开发中的各种常用的方法进行封装。也包括web方面的工具
### web项目
创建`main.go`文件和同目录的`application.yml` 文件

```text
├── application.yaml
├── go.mod
└── main.go
```

```yaml
# application.yaml 内容
gole:
  server:
    # http服务
    http:
      # 是否启用，默认：false
      enable: true
```

```go
// main.go 文件
package main

import (
    "errors"
    "github.com/gin-gonic/gin"
	httpServer "github.com/simonalong/gole-boot/server/http"
    "github.com/simonalong/gole-boot/server/http/rsp"
)

func main() {
    httpServer.Get("api/get", GetData)
    httpServer.Get("api/err", GetDataErr)
    httpServer.RunServer()
}

func GetData(c *gin.Context) (any, error) {
    return "value", nil
}

func GetDataErr(c *gin.Context) (any, error) {
    return "value", errors.New("异常")
}
```
运行如下
```shell
// 正常
root@user ~> curl http://localhost:8080/api/get
{
    "code": 0,
    "data": "value",
    "message": "success"
}

// 异常
root@uer ~> curl http://localhost:8180/api/err
{
    "code": 500,
    "message": "服务器错误",
    "error": "异常"
}
```

### 包列表
| 包名                                            |          简介          |
|-----------------------------------------------|:--------------------:|
| [debug](/debug)                               |     线上调试工具统一介绍文档     |
| [errorx](/errorx)                             |         异常分类         |
| [event](/event)                               |         事件分类         |
| [otel](otel)                                  | opentelemetry埋点客户端封装 |
| [meter](meter)                                |       自定义埋点封装        |
| [server/grpc](/server/grpc)                   |        grpc服务端        |
| [server/http](/server/http)                   |       http服务端        |
| [server/tcp](/server/tcp)                     |        tcp服务端        |
| [server/winsrv](/server/winsrv)               |      windows服务端      |
| [client/grpc](client/grpc)                    |       grpc客户端        |
| [client/http](client/http)                    |       http客户端        |
| [extend/orm](/extend/orm)                     |     gorm、xorm的封装     |
| [extend/etcd](/extend/etcd)                   |        etcd封装        |
| [extend/redis](/extend/redis)                 |     go-redis的封装      |
| [extend/emqx](/extend/emqx)                   |      emqx客户端的封装      |
| [extend/kafka](/extend/kafka)                 |     kafka客户端的封装      |
| [extend/tdengine](/extend/kafka)              |    tdengine客户端的封装    |
| [extend/nats](/extend/nats)                   |      nats客户端的封装      |
| [extend/rabbitmq](/extend/rabbitmq)           |    rabbitmq客户端的封装    |
| [extend/config_center](/extend/config_center) |      配置中心的客户端      |
| [extend/job](/extend/job)                     |    分布式任务调度客户端     |

### gole-boot 框架测试
根目录提供go_test.sh文件，统一执行所有gole-boot中包的测试模块
```shell
sh go_test.sh
```
