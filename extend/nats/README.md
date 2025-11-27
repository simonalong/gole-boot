# nats

对业内的nats客户端进行配置化封装，用于简化获取。

nats包分为两类模式，一类是nats模式，一类是nats-jetstream模式（简称js模式），这两类相互间可以数据发送和订阅，只是js下是持久化的，可以保证数据至少一次发送到接收端可以保证消息的不丢失；业务使用时候请自行选择

## 快速使用：nats
配置文件（更多配置看最下面）
```yaml
gole:
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
```

简单示例：
```go
package nats

import (
    baseNats "github.com/simonalong/gole-boot/extend/nats"
    "github.com/simonalong/gole/logger"
    "testing"
)

func TestNatsConnectWithUser(t *testing.T) { 
    // 创建新的client，也支持获取单例客户端（整个项目只有一个客户端）
    nc, err := baseNats.GetClient()
    //nc, err := baseNats.New()
    if err != nil {
        logger.Fatal(err)
        return
    }
	
    _, err = nc.Subscribe("test.connect", func(msg *baseNats.MsgOfNats) {
        // xxx
    })

    err = nc.Publish("test.connect", []byte("hello world"))
}
```

相关API
```go
// ------------------ 订阅发布模式 ------------------
// 发布消息
func (client *Client) Publish(subj string, data []byte) error {}
func (client *Client) PublishMsg(m *nats.Msg) error {}
func (client *Client) PublishRequest(subj, reply string, data []byte) error {}

// 订阅消息
func (client *Client) Subscribe(subj string, cb MsgOfNatsHandler) (*nats.Subscription, error) {}
func (client *Client) QueueSubscribe(subj, queue string, cb MsgOfNatsHandler) (*nats.Subscription, error) {}
func (client *Client) ChanQueueSubscribe(subj, group string, ch chan *nats.Msg) (*nats.Subscription, error) {}
func (client *Client) ChanSubscribe(subj string, ch chan *nats.Msg) (*nats.Subscription, error) {}

// ------------------ 请求响应模式 ------------------
// 请求消息
func (client *Client) RequestMsg(msg *nats.Msg, timeout time.Duration) (*nats.Msg, error) {}
func (client *Client) RequestWithContext(ctx context.Context, subj string, data []byte) (*nats.Msg, error) {}
func (client *Client) RequestMsgWithContext(ctx context.Context, msg *nats.Msg) (*nats.Msg, error) {}

// 响应消息：（注意：这个是msg直接返回）
func (m *MsgOfNats) Respond(data []byte) error {}
```

## 快速使用：nats-jetstream
jetstream里面最重要的是两个部分流和消费者，这里配置也是两部分
1. 流的配置
2. 消费者配置

配置文件
```yaml
gole:
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
      # 流的配置
      streams:
        - name: stream-name1
          # Subjects 是流正在监听的主题列表。支持通配符。如果流是作为镜像创建的，则无法设置主题。
          subjects:
            - test.*.req
      # 消费者配置
      consumers:
        # Name是消费者的可选名称。如果未设置，则会自动生成一个。
        - name: consumer1
          # Durable是消费者可选的持久名称。如果同时设置了Durable和Name，则它们必须相等
          durable: consumer1
```

简单示例：
```go
package nats

import (
    "context"
    "github.com/nats-io/nats.go"
    "github.com/nats-io/nats.go/jetstream"
    "github.com/stretchr/testify/assert"
    "github.com/simonalong/gole/config"
    baseNats "github.com/simonalong/gole-boot/extend/nats"
    "github.com/simonalong/gole/logger"
    baseTime "github.com/simonalong/gole/time"
    "testing"
    "time"
)

func TestNatsConnectWithUser(t *testing.T) {
    // 获取 nats 和 natsJetstream 对象；这两个对象建议保存到业务代码中
	// 也是支持单例获取和创建
    nc, js, err := baseNats.GetJetStreamClient()
    //nc, js, err := baseNats.NewJetStream()
    if err != nil {
        logger.Fatal(err)
        return
    }

    // 获取消费者，请使用配置中的流名称和消费者名称
    consumer, err := baseNats.GetStreamConsumer(js, "stream-name1", "consumer1")
    if err != nil {
        logger.Fatal(err)
        return
    }

	
    // 订阅消息
    _, err = consumer.Consume(func(msg jetstream.Msg) {
        // xxxx
		
        msg.Ack()
    })

    // 发布消息
    _, err = js.Publish(context.Background(), "test.pub.req", []byte("hello world"))
    time.Sleep(12 * time.Hour)
}
````

相关api
```go
// ------------------ 订阅发布模式 ------------------
// 发布消息
func (jsClient *JetStreamClient) Publish(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {}
func (jsClient *JetStreamClient) PublishMsg(ctx context.Context, msg *nats.Msg, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {}
func (jsClient *JetStreamClient) PublishAsync(subject string, payload []byte, opts ...jetstream.PublishOpt) (jetstream.PubAckFuture, error) {}
func (jsClient *JetStreamClient) PublishMsgAsync(msg *nats.Msg, opts ...jetstream.PublishOpt) (jetstream.PubAckFuture, error) {}

// nats非jetstream的api也可以发布消息：相关api同上

// 订阅消息：consume （这个支持埋点）
func (jsConsumer *JetStreamConsumer) Consume(handler jetstream.MessageHandler, opts ...jetstream.PullConsumeOpt) (jetstream.ConsumeContext, error) {}

// 订阅消息：fetch（以下不支持埋点）
Fetch(batch int, opts ...FetchOpt) (MessageBatch, error) {}
Fetch(batch int, opts ...FetchOpt) (MessageBatch, error) {}
FetchBytes(maxBytes int, opts ...FetchOpt) (MessageBatch, error){}
FetchNoWait(batch int) (MessageBatch, error){}
Messages(opts ...PullMessagesOpt) (MessagesContext, error){}
Next(opts ...FetchOpt) (Msg, error){}

// ------------------ kv存储模式（todo）------------------

// ------------------ 对象存储模式（todo）------------------

```

提示：
非jetstream和jetstream的nats相互间发送和接收都是可以的，只要subject可以覆盖

## 全部配置
nats的配置太多了，我这里分了三部分，用到的话，后续可以参考
1. nats 非jetstream配置
2. nats-jetstream 流配置
3. nats-jetstream 消费配置

### nats 全部配置
```yaml
gole:
  nats:
    enable: true
    # --------------- 基本信息 ---------------
    # nats的url，可以为一个，集群情况下，可以填写多个（多个之间英文逗号区分）；也可以不填，默认：nats://127.0.0.1:4222
    url: nats://127.0.0.1:4222
    # 客户端名字
    name: xx-demo-service
    # --------------- 连接认证 ---------------
    # 认证方式：用户名和密码
    user-name: admin
    # 认证方式：密码，这个密码一定要用明文（服务端加密的话，代码中要用明文）
    password: admin-demo123@xxxx.com
    # 认证方式：token
    token: xxxxxxxxxxxxxxx
    # 认证方式：NKey；种子文件地址
    nk-seed-file: ./nkeys/seed.txt
    # 证书验证文件
    credentials-file: ~/.local/share/nats/nsc/keys/creds/xxxOperate/xxxAccount/xxxUser.creds

    #--------------- 连接配置 ---------------
    # 允许在遇到与当前服务器的断开连接时使用重新连接逻辑。
    allow-reconnect: true
    # 最大重连次数。如果是否定的，那么它永远不会放弃尝试重新连接。默认为60。
    max-reconnect: 100
    # ReconnectWait 重连后等待时间。默认为2s
    reconnect-wait: 12s
    # ReconnectJitter 设置在未使用TLS的重新连接期间添加到ReconnectWait的随机延迟的上限。默认为100ms。
    reconnect-jitter: 100ms
    # ReconnectJitterTLS 设置使用TLS时在重新连接期间添加到ReconnectWait的随机延迟的上限。默认为1秒。如下两个配置都可以
    reconnect-jitter-tLS: 1s
    #  Timeout 设置连接上拨号操作的超时时间。默认为2s。
    timeout: 2s
    # DrainTimeout 设置数据流排完的超时时间。默认为30秒。
    drain-timeout: 30s
    #  FlusherTimeout 是等待底层连接的写入操作完成（包括刷新循环）的最长时间。默认为1分钟
    flusher-timeout: 1m
    # PingInterval是客户端向服务器发送ping命令的时间段，如果为0或负数，则禁用。默认为2分钟
    ping-interval: 2m
    # MaxPingsOut是在引发ErrStaleConnection错误之前可以等待响应的挂起ping命令的最大数量。默认为2。
    max-pings-out: 2
    #  ReconnectBufSize 是重新连接期间支持bufio的大小。一旦完成此操作，发布操作将返回错误。默认值为8388608字节（8MB）。
    reconnect-buf-size: 8388608
    # SubChanLen是套接字Go例程和SyncSubscriptions的消息传递之间使用的缓冲通道的大小。注意：这不会影响由PendingLimits（）指定的AsyncSubscriptions。默认值为65536。
    sub-chan-len: 65536
    #  如果无法连接到初始设置中的服务器，RetryOnFailedConnect会立即将连接设置为重新连接状态。
    #  MaxReconnect和ReconnectWait选项用于此过程，类似于断开已建立的连接。如果设置了ReconnectHandler，
    #  则它将在第一次成功的重新连接尝试时被调用（如果初始连接失败），如果设置了ClosedHandler，则在连接失败时（在用尽MaxReconnect尝试后）将被调用。
    retry-on-failed-connect: false

    #--------------- websocket 配置 ---------------
    # 对于websocket连接，向服务器指示该连接支持压缩。如果服务器也这样做，那么数据将被压缩。
    compression: false
    # 对于websocket连接，向连接url添加一个路径。当连接到代理后面的NATS时，这很有用。
    proxy-path: ws://xxx/xxx/xx

    #--------------- 其他配置 ---------------
    # NoRandomize 配置我们是否对服务器池进行随机化
    no-randomize: false
    # NoEcho 配置如果我们也有匹配的订阅，服务器是否会回显在此连接上发送的消息。请注意，version >= 1.2
    no-echo: false
    #  详细指示服务器为服务器成功处理的命令发送OK确认
    verbose: false
    #  Pedantic向服务器发出信号，表示是否应该对受试者进行进一步验证
    pedantic: false
    #  UseOldRequestStyle 强制使用旧的请求方法，即为每个请求使用新的收件箱和新的订阅。
    use-old-request-style: false
    #  NoCallbacksAfterClientClose 允许在调用Close（）后阻止调用回调。当用户代码调用Close时，客户端将不会收到通知。默认情况下是调用回调。
    no-callbacks-after-client-close: false
    #  InboxPrefix允许自定义默认的_INBOX前缀
    inbox-prefix: _api_demo
    #  IgnoreAuthErrorAbort-如果设置为true，则客户端将退出默认连接行为，即如果服务器两次返回相同的身份验证错误，则中止后续的重新连接尝试（无论重新连接策略如何）。
    ignore-auth-error-abort: false
    #  SkipHostLookup跳过服务器主机名的DNS查找。
    skip-host-lookup: false
```

### nats-jetstream 流全部配置
说明：有些结构字段太多了，我这里没有贴全，使用的时候，直接看nats相关的结构即可
```yaml
gole:
  nats:
    enable: true
    
    # ...... 这里省略nats配置 ......
    
    # nats-js的配置
    jetstream:
      # 是否使用jetstream
      enable: true
      # 流配置
      streams:
        - name: stream-name1
          description: 描述1
          # Subjects 是流正在监听的主题列表。支持通配符。如果流是作为镜像创建的，则无法设置主题。
          subjects:
            - test1.*.req
            - test2.*.req
          #--------------- 流的消息限制和留存策略 ---------------
          # MaxConsumers指定流允许的最大消费者数量。
          max-consumers: 1000
          #	MaxMsgs是流将存储的最大消息数。
          #	达到限制后，流将遵守丢弃策略。
          #	如果未设置，服务器默认值为-1（无限制）。
          max-msgs: 1000
          #	MaxBytes是流将存储的消息的最大总大小。
          #	达到限制后，流将遵守丢弃策略。
          #	如果未设置，服务器默认值为-1（无限制）。
          max-bytes: 1000
          #	MaxAge是流将保留的消息的最大年龄。
          max-age: 3h
          #	MaxMsgsPerSubject 是流将保留的每个主题的最大消息数。
          max-msgs-per-subject: 1000
          # MaxMsgSize 是流中任何单个消息的最大大小。
          max-msg-size: 1000
          # Retention 定义了流的消息保留策略。默认为LimitsPolicy。目前支持三种：0-Limits, 1-Interest, 2-WorkQueue
          retention: 2
          #	Discard 定义了当流在消息数量或总字节数方面达到限制时的消息处理策略；0-DiscardOld；1-DiscardNew；默认0
          discard: 1
          #	DiscardNewPerSubject 是一个标志，用于在达到限制时丢弃每个主题的新消息。要求DiscardPolicy设置为DiscardNew，并设置MaxMsgsPerSubject。
          discard-new-per-subject: true
          # Storage指定用于流的存储后端的类型。0-FileStorage；1-MemoryStorage
          storage: 1

          #--------------- 集群配置 ---------------
          # Replicas是集群JetStream中的流副本数量。
          # 默认值为1，最大值为5。
          replicas: 3
          #	放置用于通过标签和/或显式集群名称声明流应放置在何处
          placement:
            cluster: cluster-name
            tags:
              - "tag1"

          #--------------- 来源配置 ---------------
          # Mirror定义了镜像另一个流的配置。
          mirror:
            name: mirror-name

          # Source是此流消息来源的其他流的列表。
          sources:
            - name: source0-name

          #--------------- 操作配置 ---------------
          # 密封流不允许通过限制或API发布或删除消息，
          # 密封的流不能通过配置更新解除密封。只能通过更新API对已创建的流进行设置。
          sealed: true
          # DenyDelete限制通过API从流中删除消息的能力。默认为false。
          deny-delete: true
          # DenyPurge限制了通过API从流中清除消息的能力。默认为false。
          deny-purge: true
          #	AllowRollup允许使用Nats Rollup标头用单个新消息替换流的所有内容或流中的主题。
          allow-rollup: true

          # --------------- 内容控制 ---------------
          # Duplicates是跟踪重复消息的窗口。如果未设置，服务器默认值为2分钟。
          duplicates: 3m
          # Compression指定消息存储压缩算法。默认设置为“无压缩”。
          compression: 1
          #	FirstSeq是流中第一条消息的初始序列号。
          first-seq: 1
          #	SubjectTransform允许对匹配的消息主题应用转换。
          subject-transform:
            source: source
            destination: destination
          #	重新发布允许在消息存储后立即将其重新发布到配置的主题。
          re-publish:
            source: source
            destination: destination
            headers-only: true
          # AllowDirect允许使用直接获取API直接访问单个消息。默认为false。
          allow-direct: true
          # MirrorDirect允许使用直接获取API直接访问原始流中的单个消息。默认为false。
          mirror-direct: true
          # --------------- 消费配置 ---------------
          # NoAck是一个标志，用于禁用此流接收到的确认消息。如果设置为true，JetStream客户端的发布方法将无法按预期工作，因为它们依赖于确认。应改用核心NATS发布方法。请注意，这将降低消息传递的可靠性。
          no-ack: true
          #	ConsumerLimits定义了消费者可以设置的某些值的限制，默认值适用于未设置这些设置的人
          consumer-limits:
            inactive-threshold: 3s
            max-ack-pending: 12
          # --------------- 其他配置 ---------------
          # 元数据是一组应用程序定义的键值对，用于关联流上的元数据。此功能需要nats服务器v2.10.0或更高版本。
          metadata:
            key1: value1
            key2: value2
```

### nats-jetstream 消费配置
```yaml
gole:
  nats:
    enable: true
    
    # ...... 这里省略nats配置 ......
    
    # nats-js的配置
    jetstream:
      # 是否使用jetstream
      enable: true
      # 流配置
      streams:
        - name: stream-name1
          subjects:
            - test1.*.req
      # 消费者配置
      consumers:
          # --------------- 基本配置 ---------------
          # Name是消费者的可选名称。如果未设置，则会自动生成一个。
        - name: consumer1
          # Durable是消费者可选的持久名称。如果同时设置了Durable和Name，则它们必须相等。除非设置了InactiveThreshold，否则不会自动清理耐用消费品。
          durable: consumer1
          # 描述
          description: 消费者描述
          # 创建有序消费者
          order: true
          # --------------- 投递配置 ---------------
          # 消息投递策略：DeliverPolicy定义了从哪个点开始从流中传递消息。默认为0：DeliverAllPolicy
          #	0：DeliverAllPolicy：从流的一开始就开始传递消息。默认设置
          #	1：DeliverLastPolicy：将使用收到的最后一个序列启动消费者。
          #	2：DeliverNewPolicy：只会传递在创建消费者后发送的新消息。
          #	3：DeliverByStartSequencePolicy：将从ConsumerConfig中配置了OptStartSeq的给定序列开始传递消息。
          #	4：DeliverByStartTimePolicy：将从ConsumerConfig中配置了OptStartTime的给定时间开始传递消息
          #	5：DeliverLastPerSubjectPolicy：将向消费者发送收到的所有主题的最后一条消息。
          deliver-policy: 3
          # 可选的序列号，用于开始消息传递。仅当DeliverPolicy设置为DeliverByStartSequencePolicy时适用
          opt-start-seq: 12
          # 开始消息传递的可选时间。仅当DeliverPolicy设置为DeliverByStartTimePolicy时适用
          opt-start-time: 2024-07-30
          #	BackOff指定了在确认失败后重试消息传递的可选回退间隔。它覆盖了AckWait。
          #	BackOff仅适用于未在指定时间内确认的消息，不适用于未确认的消息。
          #	指定的间隔数必须小于或等于MaxDeliver。如果间隔数较低，则最后一个间隔将用于所有剩余的尝试。
          back-off:
            - 3s
          #	过滤主题：FilterSubject可用于过滤从流中传递的消息。FilterSubject是FilterSubjects独有的。
          filter-subject: test.req
          #	ReplayPolicy定义了向消费者发送消息的速率。
          #	如果设置了ReplayOriginalPolicy，则消息将以与存储在流中的时间间隔相同的时间间隔发送。
          #	例如，这可用于模拟开发环境中的生产流量。如果设置了ReplayInstantPolicy，则会尽快发送消息。默认为，0：ReplayInstantPolicy
          #	0：ReplayInstantPolicy：将尽可能快地重放消息。
          #	1：ReplayOriginalPolicy：将保持与收到消息相同的时间。
          replay-policy: 1

          #	--------------- 确认配置 ---------------
          # 确认策略：定义了消费者的确认策略。默认为0：AckExplicitPolicy
          #	0：AckExplicitPolicy：要求所有消息都使用ack或nack
          #	1：AckAllPolicy：在标记序列号时，也会隐式地标记低于此序列号的所有序列
          #	2：AckNonePolicy：不要求对已传递的消息进行确认
          ack-policy: 2
          #	// AckWait定义了服务器在重新发送消息之前等待确认的时间。如果未设置，服务器默认值为30秒。
          ack-wait: 5m

          #--------------- 投递限制 ---------------
          # RateLimit指定了可选的最大消息传递速率，单位为比特每秒。
          rate-limit: 1000
          #	SampleFrequency是一个可选频率，用于对确认的可观察性采样频率进行采样。
          sample-frequency: 10s
          # MaxDeliver定义了邮件的最大投递尝试次数。适用于因ack策略而重新发送的任何邮件。如果未设置，服务器默认值为-1（无限制）
          max-deliver: 1000
          # MaxWaiting是等待完成的拉取请求的最大数量。如果未设置，这将继承流的ConsumerLimits的设置，或者（如果未设置）继承帐户设置的设置。如果两者都没有设置，则服务器默认值为512。
          max-waiting: 1000
          # MaxAckPending是未确认邮件的最大数量。一旦达到此限制，服务器将暂停向消费者发送消息。如果未设置，服务器默认值为1000。设置为-1表示无限制
          max-ack-pending: 1000
          # MaxRequestBatch是单个拉取请求可以进行的可选最大批处理大小。当使用MaxRequestMaxBytes设置时，批处理大小将受到首先达到的限制的约束。
          max-request-batch: 1000
          # MaxRequestExpires是单个拉取请求等待消息可用于拉取的最长持续时间。
          max-request-expires: 10s
          #	MaxRequestMaxBytes是给定批中可以请求的可选最大总字节数。当使用MaxRequestBatch设置时，批大小将受到首先达到的限制的约束。
          max-request-max-bytes: 1000
          # HeadersOnly指示是否只应发送消息的标头（而不发送有效载荷）。默认为false。
          headers-only: true

          # --------------- 其他配置 ---------------
          # InactiveThreshold是一个持续时间，如果消费者在指定的持续时间内处于非活动状态，则指示服务器清理消费者。
          # 默认情况下，持久消费者不会被清理，但如果设置了InactiveThreshold，则会被清理。如果没有设置，则会继承流的ConsumerLimits的设置。如果两者都没有设置，服务器默认值为5秒。
          # 如果服务器没有收到拉取请求（对于拉取消费者），或者没有检测到对传递主题的兴趣（对于推送消费者），则消费者被视为不活跃，如果没有要传递的消息则不是这样。
          inactive-threshold: 10s
          # Replicas消费者状态的副本数量。默认情况下，消费者从流中继承副本的数量。
          replicas: 3
          # MemoryStorage是一个标志，用于强制消费者使用内存存储，而不是从流中继承存储类型。
          memory-storage: true
          # FilterSubjects允许按主题过滤流中的消息。此字段仅适用于FilterSubject。需要nats服务器v2.10.0或更高版本。
          filter-subjects:
            - test.req
          # 元数据是一组应用程序定义的键值对，用于关联消费者的元数据。此功能需要nats服务器v2.10.0或更高版本。
          metadata:
            key1: value1
            key2: value2
```

更多详细用法，请见测试代码和源码
