# grpc

grpc 包是用于更加方便的开发grpc的客户端的封装的包，客户端连接服务端这边有两种方式：
- 直连方式（一般用来调试）：一对一，连接某个服务端，可以是集群也可以是单节点
- 注册中心方式（线上使用）：一对多，服务端是集群

## 使用直连方式

### 服务端
直连方式配置
```yaml
# application.yaml
gole:
  server:
    grpc:
      enable: true
      # 服务端口号
      port: 9090
```
代码
```go
package test

import (
    "github.com/simonalong/gole-boot/server/grpc"
    "testing"
)

func TestGrpcServer(t *testing.T) {
    // 注册grpc服务
    grpc.Register(&xxServiceDesc1, &xxServerImpl1{})
    grpc.Register(&xxServiceDesc2, &xxServerImpl2{})
    grpc.Register(&xxServiceDesc3, &xxServerImpl3{})

    // 启动grpc服务
    grpc.RunServer()
}
```
说明：
1. grpc的生成部分，自行生成
2. 对于服务使用http和grpc都使用的，可以看后面的方式

### 客户端
这里贴出来多个服务
```yaml
gole:
  grpc:
    enable: true
    xxx-name1:
      # 连接的域名，或者ip
      host: xx-name1-service
      # 连接的端口
      port: 9090
      # 连接配置
      connect-params:
        # 供连接完成的最短时间
        min-connect-timeout: 10s
    xxx-name2:
      # 连接的域名，或者ip
      host: xx-name2-service
      # 连接的端口
      port: 9091
      # 连接配置
      connect-params:
        # 供连接完成的最短时间
        min-connect-timeout: 10s
```
代码
```go
package test

import (
    "github.com/simonalong/gole-boot/client/grpc"
    "testing"
)

func TestGrpcClient(t *testing.T) {
    // 获取grpc客户端
    conn, err := grpc.NewClientConn("xxx-name1")
    if err != nil {
        return
    }
	
    // ---------------- 使用客户端 ----------------
    c := bizServiceXXX.NewGreeterClient(conn)
    // xxx
}
```
## 使用注册中心方式
### 服务端配置
暂时只支持注册中心nats
```yaml
# application.yaml
gole:
  grpc:
    register:
      # 是否启用注册中心，默认关闭
      enable: true
      # 注册中心类型：默认natsJs，读取的是gole.nats下面的内容，记得一定要开启jetstream模式
      type: natsJs
  server:
    grpc:
      enable: true
      # 服务端口号
      port: 9090
      # 在注册中心注册的服务名
      service-name: demo-service

  nats:
    enable: true
    # nats的url，可以为一个，集群情况下，可以填写多个（多个之间英文逗号区分）；也可以不填，默认：nats://127.0.0.1:4222
    url: nats://127.0.0.1:4222
    # 客户端名字，可以不填，默认为 gole.application.name
    name: xx-demo-service
    # 认证方式：用户名和密码
    user-name: admin
    # 认证方式：密码，这个密码一定要用明文（服务端加密的话，代码中要用明文）
    password: admin-demo123@xxxx.com
    jetstream:
      # 是否使用jetstream
      enable: true
```
代码同直连方式一样
### 客户端配置
```yaml
gole:
  grpc:
    enable: true
    register:
      # 是否启用注册中心，默认关闭
      enable: true
      # 注册中心类型：默认natsJs，读取的是gole.nats下面的内容，记得一定要开启jetstream模式
      type: natsJs
    xxx-name1:
      # 服务端在注册中心注册的服务名列表
      service-name: xxx-name1-service
      # 连接配置
      connect-params:
        # 供连接完成的最短时间
        min-connect-timeout: 10s
    xxx-name2:
      # 服务端在注册中心注册的服务名列表
      service-name: xxx-name2-service
      # 连接配置
      connect-params:
        # 供连接完成的最短时间
        min-connect-timeout: 10s

  nats:
    enable: true
    # nats的url，可以为一个，集群情况下，可以填写多个（多个之间英文逗号区分）；也可以不填，默认：nats://127.0.0.1:4222
    url: nats://127.0.0.1:4222
    # 客户端名字，可以不填，默认为 gole.application.name
    name: xx-demo-service
    # 认证方式：用户名和密码
    user-name: admin
    # 认证方式：密码，这个密码一定要用明文（服务端加密的话，代码中要用明文）
    password: admin-demo123@xxxx.com
    jetstream:
      # 是否使用jetstream
      enable: true
```
代码同直连方式一样
## grpc客户端全部配置
客户端全部配置
```yaml
gole:
  grpc:
    enable: true
    # 注册中心配置
    register:
      # 是否启用注册中心，默认关闭
      enable: true
      # 注册中心类型：默认natsJs，读取的是gole.nats下面的内容，记得一定要开启jetstream模式
      type: natsJs
    xx-name1:
      # 不使用注册中心：直连情况下的服务端
      host: localhost
      port: 9090
      # 使用注册中心：服务端在注册中心注册的服务名
      service-name: xx-name1-service
      # WithAuthority返回一个DialOption，指定用作：authority伪标头和身份验证握手中的服务器名称的值。
      authority: test
      # WithDisableServiceConfig返回一个DialOption，该选项使gRPC忽略解析器提供的任何服务配置，并向解析器提供不获取服务配置的提示。
      # 请注意，此拨号选项仅禁用解析器的服务配置。如果提供了默认服务配置，gRPC将使用默认服务配置。
      disable-service-config: true
      # WithDisableRetry返回一个禁用重试的DialOption，即使服务配置启用了重试。这不会影响透明重试，如果没有数据写入线路或远程服务器未处理RPC，透明重试将自动发生。
      disable-retry: true
      # WithDisableHealthCheck禁用此ClientConn的所有子Conn的LB通道健康检查。
      # 注意：此API是实验性的，可能会在以后的版本中更改或删除。
      disable-health-check: true
      # WithDefaultServiceConfig返回一个DialOption，用于配置默认服务配置，该配置将在以下情况下使用：
      # 1.还使用WithDisableServiceConfig，或
      # 2.名称解析器不提供服务配置或提供无效的服务配置。
      # 参数s是默认服务配置的JSON表示形式。有关服务配置的更多信息，请参阅：https://github.com/grpc/grpc/blob/master/doc/service_config.md有关使用的简单示例，请参阅：examples/features/load_balancing/client/main.go
      default-service-config-raw-jSON: json
      # WithIdleTimeout返回一个DialOption，用于配置通道的空闲超时。如果通道在配置的超时时间内处于空闲状态，即没有正在进行的RPC，也没有启动新的RPC，则通道将进入空闲模式，因此名称解析器和负载均衡器将关闭。当调用Connect（）方法或启动RPC时，通道将退出空闲模式。
      # 如果在拨号时未设置此拨号选项，则将使用默认超时30分钟，并且可以通过传递超时0来禁用空闲状态。
      # 注意：此API是实验性的，可能会在以后的版本中更改或删除。
      idle-timeout: 30m
      # WithConnectParams将ClientConn配置为使用提供的ConnectParams来创建和维护与服务器的连接。
      connect-params:
        # MinConnectTimeout是我们愿意提供连接完成的最短时间。
        min-connect-timeout: 10s
        # 退避指定连接退避的配置选项。
        backoff:
          # BaseDelay是第一次故障后回退的时间量。
          base-delay: 10s
          # 乘数是重试失败后用于乘以回退的因子。理想情况下应大于1。
          multiplier: 10
          # 抖动是回退随机化的因素。
          jitter: 10.0
          # MaxDelay是退避延迟的上限。
          max-delay: 10s
```
