# 版本说明
x.y.z[.s]
- x：大版本：重大功能重构，不兼容低位大版本，合入低位中版本
- y：中版本：新增大特性，兼容近期中版本，远期不支持，合入低位小版本
- z：小版本：小特性以及中等问题修复，合入低位微版本
- s：微版本（可能不存在）：小问题修复

## v1.6.0
### 新增
1. 新增：增加各种客户端的单例模式：nats、gorm、tdengine、redis、rabbitmq
2. 新增：新增grpc的客户端的硬编码

### 修改
1. 修改：指标的默认生效为默认不生效
2. 修改：链路的开关默认修改为true
3. 修改：应用名默认不可为空
4. 修改：优化natsJs的初始化下也支持nats的原生指标测量
5. 修改：修改grpc的服务启动监听机制问题
6. 修改：注册中心的配置值的优化
7. 修改：所有链路的引入默认处理

### 删除
1. 删除：删除链路开启下配置下面对导出的url的强依赖，修改为弱依赖

## v1.5.5
### 新增
1. 新增：新增多grpc服务端的功能

### 优化
1. 优化：统一化nats包里面的username名字为userName
2. 优化：http与grpc包的命名规范的统一

## v1.5.4
### 新增
1. 新增：支持多http
2. 新增：增加http的gin的各种配置化

### 优化
1. 调整：包http内部代码结构优化，简化代码
2. 升级：升级tdorm使用bug修复版本
3. 升级：升级tdorm版本，使用最新的v1.2.1

## v1.5.3
### 新增
1. 新增：分布式定时工具本地版
2. 新增：增加基于redis的分布式锁功能



## v1.5.2
### 新增
1. 新增：新增配置中心的使用功能
2. 新增：增加分布式定时任务的客户端使用功能
3. 新增：给http增加patch（部分更新）方法

### 修复

### 优化
1. 优化：修改http包关于返回结构还有parameter的参数类型
2. 优化：rabbitmq的consume的使用
3. 优化：http包在logger的level为debug时候自动将http的级别设置为debug
4. 优化：tdengine包的配置，用户名从原来的user-name修改为username，为了保持与其他配置的一致性
5. 升级：使用新版的cbb-base包v1.0.6版本，新增了更多的功能

## v1.5.1
### 新增

### 修复

### 优化
1. 优化：升级tdorm的版本：新版本支持websocket功能，且支持全部的tdengine的sql


## v1.5.0
整体代码重构，将基础工具包放到cbb-base项目中
### 新增
1. 新增：新增rabbitmq包的支持
2. 新增：新增配置中心的使用功能

### 修复
1. 修复：gorm的代理方式，用于解决gorm的sql的混乱问题

### 优化
1. 升级：升级tdorm的版本
2. 优化：修改http包关于返回结构还有parameter的参数类型
3. 优化：rabbitmq的consume的使用


## v1.4.2
### 新增
1. 新增：新增windows的服务
2. 新增：增加zip压缩和解压工具api
3. 新增：config包增加.env环境变量的读取
4. 新增：time包增加随机休眠的功能
5. 新增：增加异常处理包的判断api
6. 新增：file包增加目录拷贝功能

### 优化
1. 优化：优化跨域处理
2. 优化：优化logger对windows的支持（还有更多问题需要优化）
3. 优化：优化nats部分的api
4. 优化：优化config的载入用法
5. 优化：优化file文件的查询返回

### 更新
1. 修改：适配tdengine的日志包名
2. 修改：修改nats的一些api的可见性

## v1.4.1
### 删除
1. 删除：删除sqlite

### 新增
1. 新增：新增业务通用异常码
2. 新增：增加bitmap工具包

### 调整
1. 修改：调整对外接口的打印格式
2. 修改：修改异常统一封装结构

### 优化
1. 优化：优化异常处理

## v1.4.0
这个是新的工程迁移过来之后的初始版本
### 修改
1. 修改：将ISCList修改为List，将ISCMap修改为Map（不兼容旧版本）


## v1.3.1.1
### 优化
1. 优化：增加全局异常拦截的日志

## v1.3.1
### 新增
1. 增加：tcp服务的支持，支持otel
2. 增加：接入prometheus指标：gin、http、tcp、grpc、nats、orm、redis、tdengine
3. 增加：logger日志包增加租户id字段输出
4. 增加：util工具包中新增baseMap工具api
5. 增加：time工具包的毫秒转时间api
6. 增加：支持多Tdengine库实例
7. 增加：grpc的负载均衡功能
8. 增加：tcp的gnet参数配置增加支持

### 优化
1. 修改：修改服务退出的打印日志以及信号监控
2. 优化：server包的代码
3. 删除：系统多余无用的指标
4. 修改：grpc的客户端只能支持一个客户端的问题
5. 优化：tcp的一些api
6. 优化：gorm，nats的开关配置
7. 优化：退出时候，报的两个错误没有使用规范的logger包
8. 删除：两个时间名转化不合适的api
9. 修改：time解析时间使用默认时区为使用东八区


## v1.3.0
### 新增
1. 增加：丰富nats的各种功能
    - nats的账号连接功能
    - 非js模式下：发布订阅、请求回复、队列
    - js模式下：流配置、消费者配置、数据拉取、数据推送、数据批量拉取、并发拉取（不包括kv和对象存储）
2. 增加：nats对opentelemetry埋点的支持
    - nats支持otel
    - nats-js支持otel
3. 增加：对grpc的客户端和服务端支持
4. 增加：支持grpc的自动埋点
5. 增加：日志框架增加traceId日志输出

### 优化
1. 优化：修复日志模块的颜色异常问题
2. 优化：修复基本类型指针、结构体的属性指针的反射注入问题
3. 优化：升级使用tdengine-orm框架的版本v1.1.1

## v1.2.2
### 新增
1. 新增：新的业务返回方式，简化业务编写

### 优化
1. 优化：修改优化正常和异常返回结构


## v1.2.1
### 新增
1. 新增：增加BaseError用于全局异常拦截


## v1.2.0
### 新增
1. 新增：opentelemetry的配置
2. 新增：各组件对opentelemetry埋点的支持
   - gin支持otel
   - http支持otel
   - gorm支持otel
   - xorm支持otel
   - redis支持otel
   - tdengine支持otel
### 修改
1. 项目改名为gole-boot剥离seatak名字

## v1.1.0
### 新增
1. 增加tdengine的支持：使用自行开发的tdengine-orm


## v1.0.0 
### 新增：支持一些基础工具

---

当前全部功能

| 包名                                  |          简介          |
|-------------------------------------|:--------------------:|
| [util](/util)                       |      基础工具（更新中）       |
| [config](/config)                   |        配置文件管理        |
| [validate](/validate)               |         校验核查         |
| [logger](/logger)                   |          日志          |
| [database](/database)               |      数据库处理（待更新）      |
| [server](/server)                   |         服务处理         |
| [goid](/goid)                       | 局部id传递处理（theadlocal） |
| [json](/json)                       |     json字符串处理工具      |
| [cache](/cache)                     |         缓存工具         |
| [time](/time)                       |        时间管理工具        |
| [file](/file)                       |        文件管理工具        |
| [coder](/coder)                     |       编解码加解密工具       |
| [http](/client/http)                       |      http的辅助工具       |
| [listener](/event)               |        事件监听机制        |
| [bean](/bean)                       |        对象管理工具        |
| [bean](/grpc)                       |    grpc客户端和服务端封装     |
| [bean](/tcp)                        |     tcp服务端和客户端封装     |
| [debug](/debug)                     |     线上调试工具统一介绍文档     |
| [extend/orm](/extend/orm)           |     gorm、xorm的封装     |
| [extend/etcd](/extend/etcd)         |        etcd封装        |
| [extend/redis](/extend/redis)       |     go-redis的封装      |
| [extend/emqx](/extend/emqx)         |      emqx客户端的封装      |
| [extend/kafka](/extend/kafka)       |     kafka客户端的封装      |
| [extend/tdengine](/extend/tdengine) |     td-orm客户端封装      |
| [extend/nats](/extend/nats)         |      nats客户端封装       |
| [extend/otel](/extend/otel)         | opentelemetry先关功能封装  |
