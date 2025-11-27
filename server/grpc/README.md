# grpc

grpc 包是用于更加方便的开发grpc的服务端的封装的服务端使用的包，这个跟client的grpc包是不一样的，client的grpc包是用于更加方便的开发grpc的客户端的封装的包
1. 连接类型：
   - 直连方式：一对一，连接某个服务端，可以是集群也可以是单节点
   - 注册中心方式：一对多，服务端是集群
2. 服务端个数：
   - 单服务端：一个项目提供一个grpc服务
   - 多服务端：一个项目提供多个grpc服务

## 一、连接类型
### 直连方式
配置
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

### 注册中心方式
配置
暂时只支持注册中心nats
```yaml
# application.yaml
gole:
  server:
    grpc:
      enable: true
      # 服务端口号
      port: 9090
      # 在注册中心注册的服务名；如果为空，则默认使用${gole.application.name}
      service-name: demo-service
  grpc:
    register:
      # 是否启用注册中心，默认关闭
      enable: true
      # 注册中心类型：默认natsJs，读取的是gole.nats下面的内容，记得一定要开启jetstream模式
      type: natsJs
```
代码（其实跟直连是一样的，代码层不感知）
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

## 二、服务端个数
### 单grpc服务端
这是最常见的场景，相关的配置和代码见上面，我们这里主要介绍一个工程多个grpc服务端的场景
### 多grpc服务端
#### 1. 直连
配置
```yaml
# application.yaml
gole:
  server:
    grpc:
      multi:
        grpc-sever-name1:
          enable: true
          port: 9090
        grpc-sever-name2:
          enable: true
          port: 9091
```

```go
func TestGrpcServer(t *testing.T) {
    // 注册grpc服务
    server1 := grpc.Server("grpc-server-name1")
    server1.Register(&xxServiceDesc1, &xxServerImpl1{})
    server1.Register(&xxServiceDesc2, &xxServerImpl2{})
    server1.Register(&xxServiceDesc3, &xxServerImpl3{})

    // 注册grpc服务
    server2 := grpc.Server("grpc-server-name2")
    server2.Register(&xxServiceDesc1, &xxServerImpl1{})
    server2.Register(&xxServiceDesc2, &xxServerImpl2{})
    server2.Register(&xxServiceDesc3, &xxServerImpl3{})

    // 启动grpc服务
    grpc.RunServer()
}
```
#### 2. 注册中心
配置
```yaml
# application.yaml
gole:
  server:
    grpc:
      multi:
        grpc-sever-xxx-name1:
          enable: true
          # 在注册中心注册的服务名；为空的话，则默认使用 grpc-sever-xxx-name1 这个为服务名
          service-name: demo-service
          port: 9090
        grpc-sever-name2:
          enable: true
          # 在注册中心注册的服务名；为空的话，则默认使用 grpc-sever-xxx-name1 这个为服务名
          service-name: demo-service
          port: 9091
      
  grpc:
     register:
        # 是否启用注册中心，默认关闭
        enable: true
        # 注册中心类型：默认natsJs，读取的是gole.nats下面的内容，记得一定要开启jetstream模式
        type: natsJs
```


#### 客户端连接
客户端连接其实就是按照多服务的方式连接，这里我们举个例子，详情可以见client部分的grpc包

## 三、grpc全部配置
服务全部配置：以下配置的值只是自己随便写的，使用时候具体请看相关文档

注意：
- 如果单服务配置和多服务配置同时配置，则只生效单服务
```yaml
gole:
  server:
    grpc:
      enable: true
      # 端口；默认9090
      port: 9090
      # 在注册中心注册的服务名；如果为空，则默认使用${gole.application.name}
      service-name: demo-service
      # MaxConcurrentStreams返回一个ServerOption，该选项将对每个ServerTransport的并发流数量施加限制。
      max-concurrent-streams: 1000
      # MaxRecvMsgSize 返回一个ServerOption，用于设置服务器可以接收的最大消息大小（以字节为单位）。如果没有设置，gRPC将使用默认的4MB。
      max-receive-message-size: 4194304
      # MaxSendMsgSize返回一个ServerOption，用于设置服务器可以发送的最大消息大小（以字节为单位）。如果没有设置，gRPC将使用默认的`math.MaxInt32
      max-send-message-size: 4194304
      # KeepaliveParams返回一个ServerOption，为服务器设置keepalive和最大年龄参数。
      keepalive-params:
        # MaxConnectionIdle是一个持续时间，在此时间之后，通过发送GoAway来关闭空闲连接。空闲持续时间是自最近一次未完成的RPC数量变为零或连接建立以来定义的。
        # 当前默认值为无穷大
        max-connection-idle: 12s
        # MaxConnectionAge是连接在发送GoAway关闭之前可能存在的最长时间的持续时间。MaxConnectionAge将添加+/-10%的随机抖动，以分散连接风暴。
        # 当前默认值为无穷大
        max-connection-age: 10h
        # MaxConnectionAgeGrace是MaxConnectionAge之后的一个附加时间段，在此时间段之后，连接将被强制关闭。
        # 当前默认值为无穷大。
        max-connection-age-grace: 12h
        # 经过一段时间后，如果服务器没有看到任何活动，它会向客户端发送ping消息，以查看传输是否仍然有效。如果设置为小于1秒，则将使用最小值1秒。
        # 当前默认值为2小时。
        time: 2h
        # 在ping了保活检查后，服务器会等待一段超时时间，如果在此之后没有看到任何活动，则关闭连接。
        # 当前默认值为20秒
        timeout: 20s
      # InitialWindowSize返回一个ServerOption，用于设置流的窗口大小。窗口大小的下限为64K，任何小于此值的值都将被忽略。
      initial-window-size: 65536
      # InitialConnWindowSize返回一个ServerOption，用于设置连接的窗口大小。窗口大小的下限为64K，任何小于此值的值都将被忽略。
      initial-conn-window-size: 65536
      # WriteBufferSize决定在线路上进行写入之前可以批处理多少数据。此缓冲区的默认值为32KB。零值或负值将禁用写缓冲区，以便每次写入都在底层连接上。注意：发送呼叫可能不会直接转换为写入。
      write-buffer-size: 32768
      # ReadBufferSize允许您设置读取缓冲区的大小，这决定了一次读取系统调用最多可以读取多少数据。此缓冲区的默认值为32KB。零值或负值将禁用连接的读取缓冲区，以便数据帧生成器可以直接访问底层连接。
      read-buffer-size: 32768
      # SharedWriteBuffer允许重用每个连接的传输写缓冲区。如果此选项设置为true，则每个连接都会在刷新线路上的数据后释放缓冲区。
      # 注意：本API是实验性的，可能会在稍后释放。
      shared-write-buffer: true
      # ConnectionTimeout返回一个ServerOption，用于设置所有新连接的连接建立超时时间（包括HTTP/2握手）。如果未设置，默认值为120秒。零值或负值将导致立即超时。
      # 注意：本API是实验性的，可能会在稍后释放。
      connection-timeout: 120s
      # MaxHeaderListSizeServerOption是一个ServerOption，用于设置服务器准备接受的标头列表的最大（未压缩）大小。
      max-header-list-size: 512
      # HeaderTableSize返回一个ServerOption，用于设置流的动态头表的大小。
      # 注意：此API是实验性的，可能会在以后的版本中更改或删除。
      header-table-size: 512
      # NumStreamWorkers返回一个ServerOption，用于设置应用于处理传入流的worker goroutines的数量。将其设置为零（默认）将禁用workers并为每个流生成一个新的goroutine。
      # 注意：此API是实验性的，可能会在以后的版本中更改或删除。
      num-server-workers: 100
      # WaitForHandlers使Stop等待所有未完成的方法处理程序退出，然后返回。如果为false，Stop将在所有连接关闭后立即返回，但方法处理程序可能仍在运行。默认情况下，Stop不会等待方法处理程序返回。
      # 注意：此API是实验性的，可能会在以后的版本中更改或删除。
      wait-for-handlers: true
    multi:
       xxxx-name1:
          enable: true
          # 端口
          port: 9090
          # 在注册中心注册的服务名
          service-name: demo-service
          # MaxConcurrentStreams返回一个ServerOption，该选项将对每个ServerTransport的并发流数量施加限制。
          max-concurrent-streams: 1000
          # MaxRecvMsgSize 返回一个ServerOption，用于设置服务器可以接收的最大消息大小（以字节为单位）。如果没有设置，gRPC将使用默认的4MB。
          # ................ （下面省略）同单个服务的配置一样 ................ 
       xxxx-name2:
          enable: true
          # 端口
          port: 9090
          # 在注册中心注册的服务名
          service-name: demo-service
          # MaxConcurrentStreams返回一个ServerOption，该选项将对每个ServerTransport的并发流数量施加限制。
          max-concurrent-streams: 1000
          # MaxRecvMsgSize 返回一个ServerOption，用于设置服务器可以接收的最大消息大小（以字节为单位）。如果没有设置，gRPC将使用默认的4MB。
          # ................ （下面省略）同单个服务的配置一样 ................
      # ........ （下面省略）更多服务的配置 ........ 
  grpc:
     enable: true
     # 注册中心配置
     register:
        # 是否启用注册中心，默认关闭
        enable: true
        # 注册中心类型：默认natsJs，读取的是gole.nats下面的内容，记得一定要开启jetstream模式
        type: natsJs
```
