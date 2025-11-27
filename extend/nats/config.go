package nats

import (
	"github.com/nats-io/nats.go/jetstream"
	"time"
)

type ConfigOfNats struct {
	// # nats的url，可以为一个，集群情况下，可以填写多个（多个之间英文逗号区分）；也可以不填，默认：nats://127.0.0.1:4222
	Url string
	// 客户端名字
	Name string

	// --------------- 连接认证 ---------------
	// 认证方式：用户名和密码
	UserName string
	// 认证方式：密码，这个密码一定要用明文（服务端加密的话，代码中要用明文）
	Password string
	// 认证方式：token
	Token string
	// 认证方式：NKey；种子文件地址
	NkSeedFile string
	// credentials 证书验证文件
	CredentialsFile string

	// --------------- tls安全配置 ---------------
	// Secure启用默认跳过服务器验证的TLS安全连接。不推荐
	//Secure bool
	// TLSConfig 是用于安全传输的自定义TLS配置
	//TLSConfig *tls.Config
	// TLSHandshakeFirst 用于指示库在连接后和从服务器接收INFO协议之前执行TLS握手。如果启用了此选项，但服务器未配置为首先执行TLS握手，则连接将失败
	//TLSHandshakeFirst bool

	// --------------- 连接配置 ---------------
	// AllowReconnect 允许在遇到与当前服务器的断开连接时使用重新连接逻辑。
	AllowReconnect bool
	// MaxReconnect 最大重连次数。如果是否定的，那么它永远不会放弃尝试重新连接。默认为60。
	MaxReconnect int
	// ReconnectWait 重连后等待时间。默认为2s
	ReconnectWait time.Duration
	// ReconnectJitter 设置在未使用TLS的重新连接期间添加到ReconnectWait的随机延迟的上限。默认为100ms。
	ReconnectJitter time.Duration
	// ReconnectJitterTLS 设置使用TLS时在重新连接期间添加到ReconnectWait的随机延迟的上限。默认为1秒。
	ReconnectJitterTLS time.Duration
	// Timeout 设置连接上拨号操作的超时时间。默认为2s。
	Timeout time.Duration
	// DrainTimeout 设置数据流排完的超时时间。默认为30秒。
	DrainTimeout time.Duration
	// FlusherTimeout 是等待底层连接的写入操作完成（包括刷新循环）的最长时间。默认为1分钟
	FlusherTimeout time.Duration
	// PingInterval是客户端向服务器发送ping命令的时间段，如果为0或负数，则禁用。默认为2分钟
	PingInterval time.Duration
	// MaxPingsOut是在引发ErrStaleConnection错误之前可以等待响应的挂起ping命令的最大数量。默认为2。
	MaxPingsOut int
	// ReconnectBufSize 是重新连接期间支持bufio的大小。一旦完成此操作，发布操作将返回错误。默认值为 8388608 字节（8MB）。
	ReconnectBufSize int
	// SubChanLen是套接字Go例程和SyncSubscriptions的消息传递之间使用的缓冲通道的大小。注意：这不会影响由PendingLimits（）指定的AsyncSubscriptions。默认值为 65536。
	SubChanLen int
	// 如果无法连接到初始设置中的服务器，RetryOnFailedConnect会立即将连接设置为重新连接状态。
	// MaxReconnect和ReconnectWait选项用于此过程，类似于断开已建立的连接。如果设置了ReconnectHandler，
	//则它将在第一次成功的重新连接尝试时被调用（如果初始连接失败），如果设置了ClosedHandler，则在连接失败时（在用尽MaxReconnect尝试后）将被调用。
	RetryOnFailedConnect bool

	// --------------- websocket 配置 ---------------
	// 对于websocket连接，向服务器指示该连接支持压缩。如果服务器也这样做，那么数据将被压缩。
	Compression bool
	// 对于websocket连接，向连接url添加一个路径。当连接到代理后面的NATS时，这很有用。
	ProxyPath string

	// --------------- 其他配置 ---------------
	// NoRandomize 配置我们是否对服务器池进行随机化
	NoRandomize bool
	// NoEcho 配置如果我们也有匹配的订阅，服务器是否会回显在此连接上发送的消息。请注意，version >= 1.2
	NoEcho bool
	// 详细指示服务器为服务器成功处理的命令发送OK确认
	Verbose bool
	// Pedantic向服务器发出信号，表示是否应该对受试者进行进一步验证
	Pedantic bool
	// UseOldRequestStyle 强制使用旧的请求方法，即为每个请求使用新的收件箱和新的订阅。
	UseOldRequestStyle bool
	// NoCallbacksAfterClientClose 允许在调用Close（）后阻止调用回调。当用户代码调用Close时，客户端将不会收到通知。默认情况下是调用回调。
	NoCallbacksAfterClientClose bool
	// InboxPrefix允许自定义默认的_INBOX前缀
	InboxPrefix string
	// IgnoreAuthErrorAbort-如果设置为true，则客户端将退出默认连接行为，即如果服务器两次返回相同的身份验证错误，则中止后续的重新连接尝试（无论重新连接策略如何）。
	IgnoreAuthErrorAbort bool
	// SkipHostLookup跳过服务器主机名的DNS查找。
	SkipHostLookup bool

	// --------------- jetstream配置 ---------------
	// jetstream的配置
	Jetstream ConfigOfJetstream
}

type ConfigOfJetstream struct {
	Enable bool

	// --------------- stream 配置 ---------------
	Streams []ConfigOfJetstreamStream
	// --------------- consumer 配置 ---------------
	Consumers []ConfigOfJetstreamConsumer
	// --------------- kv 配置 ---------------
	KvStores []ConfigOfJsKvStores
}

type ConfigOfJetstreamStream struct {
	// --------------- 基本配置 ---------------
	Name        string
	Description string
	// Subjects 是流正在监听的主题列表。支持通配符。如果流是作为镜像创建的，则无法设置主题。
	Subjects []string

	// --------------- 流的消息限制和留存策略 ---------------
	// MaxConsumers指定流允许的最大消费者数量。
	MaxConsumers int
	// MaxMsgs是流将存储的最大消息数。
	// 达到限制后，流将遵守丢弃策略。
	// 如果未设置，服务器默认值为-1（无限制）。
	MaxMsgs int64
	// MaxBytes是流将存储的消息的最大总大小。
	// 达到限制后，流将遵守丢弃策略。
	// 如果未设置，服务器默认值为-1（无限制）。
	MaxBytes int64
	// MaxAge是流将保留的消息的最大年龄。
	MaxAge time.Duration
	// MaxMsgsPerSubject 是流将保留的每个主题的最大消息数。
	MaxMsgsPerSubject int64
	// MaxMsgSize 是流中任何单个消息的最大大小。
	MaxMsgSize int32
	// Retention 定义了流的消息保留策略。默认为LimitsPolicy。目前支持三种：0-Limits, 1-Interest, 2-WorkQueue
	Retention jetstream.RetentionPolicy
	// Discard 定义了当流在消息数量或总字节数方面达到限制时的消息处理策略；0-DiscardOld；1-DiscardNew；默认0
	Discard jetstream.DiscardPolicy
	// DiscardNewPerSubject 是一个标志，用于在达到限制时丢弃每个主题的新消息。要求DiscardPolicy设置为DiscardNew，并设置MaxMsgsPerSubject。
	DiscardNewPerSubject bool
	// Storage指定用于流的存储后端的类型。0-FileStorage；1-MemoryStorage
	Storage jetstream.StorageType

	// --------------- 集群配置 ---------------
	// Replicas是集群JetStream中的流副本数量。
	// 默认值为1，最大值为5。
	Replicas int
	// 放置用于通过标签和/或显式集群名称声明流应放置在何处
	Placement *jetstream.Placement

	// --------------- 来源配置 ---------------
	// Mirror定义了镜像另一个流的配置。
	Mirror *jetstream.StreamSource

	// Mirror定义了镜像另一个流的配置。
	// Mirror*StreamSource `json：“镜像，omitempty”`
	// Source是此流消息来源的其他流的列表。
	Sources []jetstream.StreamSource

	// --------------- 操作配置 ---------------
	// 密封流不允许通过限制或API发布或删除消息，
	// 密封的流不能通过配置更新解除密封。只能通过更新API对已创建的流进行设置。
	Sealed bool
	// DenyDelete限制通过API从流中删除消息的能力。默认为false。
	DenyDelete bool
	// DenyPurge限制了通过API从流中清除消息的能力。默认为false。
	DenyPurge bool
	// AllowRollup允许使用Nats Rollup标头用单个新消息替换流的所有内容或流中的主题。
	AllowRollup bool

	// --------------- 内容控制 ---------------
	// Duplicates是跟踪重复消息的窗口。
	// 如果未设置，服务器默认值为2分钟。
	Duplicates time.Duration
	// Compression指定消息存储压缩算法。
	// 默认设置为“无压缩”。
	Compression jetstream.StoreCompression
	// FirstSeq是流中第一条消息的初始序列号。
	FirstSeq uint64
	// SubjectTransform允许对匹配的消息主题应用转换。
	SubjectTransform *jetstream.SubjectTransformConfig
	// 重新发布允许在消息存储后立即将其重新发布到配置的主题。
	RePublish *jetstream.RePublish
	// AllowDirect允许使用直接获取API直接访问单个消息。默认为false。
	AllowDirect bool
	// MirrorDirect允许使用直接获取API直接访问原始流中的单个消息。默认为false。
	MirrorDirect bool

	// --------------- 消费配置 ---------------
	// NoAck是一个标志，用于禁用此流接收到的确认消息。
	// 如果设置为true，JetStream客户端的发布方法将无法按预期工作，因为它们依赖于确认。应改用核心NATS发布方法。请注意，这将降低消息传递的可靠性。
	NoAck bool
	// ConsumerLimits定义了消费者可以设置的某些值的限制，默认值适用于未设置这些设置的人
	ConsumerLimits jetstream.StreamConsumerLimits

	// --------------- 其他配置 ---------------
	// 元数据是一组应用程序定义的键值对，用于关联流上的元数据。此功能需要nats服务器v2.10.0或更高版本。
	Metadata map[string]string
}

type ConfigOfJetstreamConsumer struct {
	// --------------- 基本配置 ---------------
	// Name是消费者的可选名称。如果未设置，则会自动生成一个。
	// 名称不能包含空格、.、*、>、，路径分隔符（正斜杠或反斜杠）和不可打印字符。
	Name string
	// Durable是消费者可选的持久名称。如果同时设置了Durable和Name，则它们必须相等。除非设置了InactiveThreshold，否则不会自动清理耐用消费品。
	// Durable不能包含空格、.、*、>、，路径分隔符（正斜杠或反斜杠）和不可打印字符。
	Durable string
	// 描述
	Description string
	// 是否为有序消费
	Order bool

	// --------------- 投递配置 ---------------
	// 消息投递策略：DeliverPolicy定义了从哪个点开始从流中传递消息。默认为0：DeliverAllPolicy。【OrderConsumer可以用】
	//
	//	0：DeliverAllPolicy：从流的一开始就开始传递消息。默认设置
	//	1：DeliverLastPolicy：将使用收到的最后一个序列启动消费者。
	//	2：DeliverNewPolicy：只会传递在创建消费者后发送的新消息。
	//	3：DeliverByStartSequencePolicy：将从ConsumerConfig中配置了OptStartSeq的给定序列开始传递消息。
	//	4：DeliverByStartTimePolicy：将从ConsumerConfig中配置了OptStartTime的给定时间开始传递消息
	//	5：DeliverLastPerSubjectPolicy：将向消费者发送收到的所有主题的最后一条消息。
	DeliverPolicy jetstream.DeliverPolicy
	// 可选的序列号，用于开始消息传递。仅当DeliverPolicy设置为DeliverByStartSequencePolicy时适用。【OrderConsumer可以用】
	OptStartSeq uint64
	// 开始消息传递的可选时间。仅当DeliverPolicy设置为DeliverByStartTimePolicy时适用。【OrderConsumer可以用】
	OptStartTime *time.Time
	// BackOff指定了在确认失败后重试消息传递的可选回退间隔。它覆盖了AckWait。
	// BackOff仅适用于未在指定时间内确认的消息，不适用于未确认的消息。
	// 指定的间隔数必须小于或等于MaxDeliver。如果间隔数较低，则最后一个间隔将用于所有剩余的尝试。
	BackOff []time.Duration
	// 过滤主题：FilterSubject可用于过滤从流中传递的消息。FilterSubject是FilterSubjects独有的。
	FilterSubject string
	// ReplayPolicy定义了向消费者发送消息的速率。【OrderConsumer可以用】
	// 如果设置了ReplayOriginalPolicy，则消息将以与存储在流中的时间间隔相同的时间间隔发送。
	// 例如，这可用于模拟开发环境中的生产流量。如果设置了ReplayInstantPolicy，则会尽快发送消息。
	// 默认为，0：ReplayInstantPolicy
	//	0：ReplayInstantPolicy：将尽可能快地重放消息。
	//	1：ReplayOriginalPolicy：将保持与收到消息相同的时间。
	ReplayPolicy jetstream.ReplayPolicy

	// --------------- 确认配置 ---------------
	// 确认策略：定义了消费者的确认策略。默认为0：AckExplicitPolicy
	//	0：AckExplicitPolicy：要求所有消息都使用ack或nack
	//	1：AckAllPolicy：在标记序列号时，也会隐式地标记低于此序列号的所有序列
	//	2：AckNonePolicy：不要求对已传递的消息进行确认
	AckPolicy jetstream.AckPolicy
	// AckWait定义了服务器在重新发送消息之前等待确认的时间。如果未设置，服务器默认值为30秒。
	AckWait time.Duration

	// --------------- 投递限制 ---------------
	// RateLimit指定了可选的最大消息传递速率，单位为比特每秒。
	RateLimit uint64
	// SampleFrequency是一个可选频率，用于对确认的可观察性采样频率进行采样。
	SampleFrequency string
	// MaxDeliver定义了邮件的最大投递尝试次数。适用于因ack策略而重新发送的任何邮件。如果未设置，服务器默认值为-1（无限制）
	MaxDeliver int
	// MaxWaiting是等待完成的拉取请求的最大数量。如果未设置，这将继承流的ConsumerLimits的设置，或者（如果未设置）继承帐户设置的设置。如果两者都没有设置，则服务器默认值为512。
	MaxWaiting int
	// MaxAckPending是未确认邮件的最大数量。一旦达到此限制，服务器将暂停向消费者发送消息。如果未设置，服务器默认值为1000。
	// 设置为-1表示无限制
	MaxAckPending int
	// MaxRequestBatch是单个拉取请求可以进行的可选最大批处理大小。当使用MaxRequestMaxBytes设置时，批处理大小将受到首先达到的限制的约束。
	MaxRequestBatch int
	// MaxRequestExpires是单个拉取请求等待消息可用于拉取的最长持续时间。
	MaxRequestExpires time.Duration
	// MaxRequestMaxBytes是给定批中可以请求的可选最大总字节数。当使用MaxRequestBatch设置时，批大小将受到首先达到的限制的约束。
	MaxRequestMaxBytes int
	// HeadersOnly指示是否只应发送消息的标头（而不发送有效载荷）。默认为false。【OrderConsumer可以用】
	HeadersOnly bool

	// --------------- 其他配置 ---------------
	// InactiveThreshold是一个持续时间，如果消费者在指定的持续时间内处于非活动状态，则指示服务器清理消费者。【OrderConsumer可以用】
	// 默认情况下，持久消费者不会被清理，但如果设置了InactiveThreshold，则会被清理。如果没有设置，则会继承流的ConsumerLimits的设置。如果两者都没有设置，服务器默认值为5秒。
	//
	// 如果服务器没有收到拉取请求（对于拉取消费者），或者没有检测到对传递主题的兴趣（对于推送消费者），则消费者被视为不活跃，如果没有要传递的消息则不是这样。
	InactiveThreshold time.Duration
	// Replicas消费者状态的副本数量。默认情况下，消费者从流中继承副本的数量。
	Replicas int
	// MemoryStorage是一个标志，用于强制消费者使用内存存储，而不是从流中继承存储类型。
	MemoryStorage bool
	// FilterSubjects允许按主题过滤流中的消息。此字段仅适用于FilterSubject。需要nats服务器v2.10.0或更高版本。【OrderConsumer可以用】
	FilterSubjects []string
	// 元数据是一组应用程序定义的键值对，用于关联消费者的元数据。此功能需要nats服务器v2.10.0或更高版本。
	Metadata map[string]string

	// 在单个重建循环中重新创建消费者的最大尝试次数。默认为无限制。【OrderConsumer可以用】
	MaxResetAttempts int
}

type ConfigOfJsKvStores struct {
	// Bucket是KeyValue存储的名称。存储桶名称必须是唯一的，并且只能包含字母数字字符、破折号和下划线。
	Bucket string
	// 描述
	Description string
	// MaxValueSize 是值的最大大小，单位为字节。如果未指定，默认值为-1（无限制）
	MaxValueSize int32
	// History 是每个键要保留的历史值的数量。如果未指定，则默认值为1。Max是64。
	History uint8
	// TTL是密钥的到期时间。默认情况下，密钥不会过期。
	TTL time.Duration
	// MaxBytes是KeyValue存储的最大字节数。如果未指定，默认值为-1（无限制）。
	MaxBytes int64
	// 存储是用于KeyValue存储的存储类型。如果未指定，则默认值为FileStorage。
	Storage jetstream.StorageType
	// Replicas是集群jetstream中KeyValue存储要保留的副本数量。默认值为1，最大值为5。
	Replicas int
	//	放置用于通过标签和/或显式集群名称声明流应放置在何处。
	Placement *jetstream.Placement
	// 重新发布允许在消息存储后立即将其重新发布到配置的主题。
	RePublish *jetstream.RePublish
	// Mirror定义了镜像另一个KeyValue存储的配置。
	Mirror *jetstream.StreamSource
	// Sources定义KeyValue存储源的配置。
	Sources []jetstream.StreamSource
	// 压缩设置底层流压缩。注意：nats服务器2.10.0支持压缩+
	Compression bool
}

func jetStreamStreamConvert(streamConfig ConfigOfJetstreamStream) jetstream.StreamConfig {
	config := jetstream.StreamConfig{
		Name:                 streamConfig.Name,
		Description:          streamConfig.Description,
		Subjects:             streamConfig.Subjects,
		Retention:            streamConfig.Retention,
		MaxConsumers:         streamConfig.MaxConsumers,
		Discard:              streamConfig.Discard,
		DiscardNewPerSubject: streamConfig.DiscardNewPerSubject,
		MaxAge:               streamConfig.MaxAge,
		MaxMsgsPerSubject:    streamConfig.MaxMsgsPerSubject,
		MaxMsgSize:           streamConfig.MaxMsgSize,
		Storage:              streamConfig.Storage,
		NoAck:                streamConfig.NoAck,
		Placement:            streamConfig.Placement,
		Mirror:               streamConfig.Mirror,
		Sources:              generateStreamSource(streamConfig.Sources),
		Sealed:               streamConfig.Sealed,
		DenyDelete:           streamConfig.DenyDelete,
		FirstSeq:             streamConfig.FirstSeq,
		SubjectTransform:     streamConfig.SubjectTransform,
		RePublish:            streamConfig.RePublish,
		AllowDirect:          streamConfig.AllowDirect,
		MirrorDirect:         streamConfig.MirrorDirect,
		ConsumerLimits:       streamConfig.ConsumerLimits,
		Metadata:             streamConfig.Metadata,
	}

	if streamConfig.MaxMsgs != 0 {
		config.MaxMsgs = streamConfig.MaxMsgs
	}

	if streamConfig.MaxBytes != 0 {
		config.MaxBytes = streamConfig.MaxBytes
	}

	if streamConfig.Replicas != 0 {
		config.Replicas = streamConfig.Replicas
	}

	if streamConfig.Duplicates != time.Duration(0) {
		config.Duplicates = streamConfig.Duplicates
	}

	return config
}

func jetStreamConsumerConvert(consumerConfig ConfigOfJetstreamConsumer) jetstream.ConsumerConfig {
	config := jetstream.ConsumerConfig{
		Name:               consumerConfig.Name,
		Durable:            consumerConfig.Durable,
		Description:        consumerConfig.Description,
		DeliverPolicy:      consumerConfig.DeliverPolicy,
		OptStartSeq:        consumerConfig.OptStartSeq,
		OptStartTime:       consumerConfig.OptStartTime,
		AckPolicy:          consumerConfig.AckPolicy,
		BackOff:            consumerConfig.BackOff,
		FilterSubject:      consumerConfig.FilterSubject,
		ReplayPolicy:       consumerConfig.ReplayPolicy,
		RateLimit:          consumerConfig.RateLimit,
		SampleFrequency:    consumerConfig.SampleFrequency,
		HeadersOnly:        consumerConfig.HeadersOnly,
		MaxRequestBatch:    consumerConfig.MaxRequestBatch,
		MaxRequestExpires:  consumerConfig.MaxRequestExpires,
		MaxRequestMaxBytes: consumerConfig.MaxRequestMaxBytes,
		Replicas:           consumerConfig.Replicas,
		MemoryStorage:      consumerConfig.MemoryStorage,
		FilterSubjects:     consumerConfig.FilterSubjects,
		Metadata:           consumerConfig.Metadata,
	}

	if consumerConfig.AckWait != time.Duration(0) {
		config.AckWait = consumerConfig.AckWait
	}

	if consumerConfig.MaxDeliver != 0 {
		config.MaxDeliver = consumerConfig.MaxDeliver
	}

	if consumerConfig.MaxWaiting != 0 {
		config.MaxWaiting = consumerConfig.MaxWaiting
	}

	if consumerConfig.MaxAckPending != 0 {
		config.MaxAckPending = consumerConfig.MaxAckPending
	}

	if consumerConfig.InactiveThreshold != time.Duration(0) {
		config.InactiveThreshold = consumerConfig.InactiveThreshold
	}

	return config
}

func jetStreamOrderConsumerConvert(consumerConfig ConfigOfJetstreamConsumer) jetstream.OrderedConsumerConfig {
	return jetstream.OrderedConsumerConfig{}
}

func generateStreamSource(streamSources []jetstream.StreamSource) []*jetstream.StreamSource {
	var sources []*jetstream.StreamSource
	for _, source := range streamSources {
		sources = append(sources, &source)
	}
	return sources
}
