package etcd

var Cfg EtcdConfig

type EtcdConfig struct {
	// etcd的服务ip:port列表
	Endpoints []string

	Username string
	Password string

	// 自动同步间隔：是用其最新成员更新端点的间隔；默认为0，即禁用自动同步；配置示例：1s、1000ms
	AutoSyncInterval string

	// 拨号超时：是指连接失败后的超时时间；配置示例：1s、1000ms
	DialTimeout string

	// 拨号保持连接时间：是客户端ping服务器以查看传输是否连接的时间；配置示例：1s、1000ms
	DialKeepAliveTime string

	// 拨号保持连接超时：是客户端等待响应保持连接探测的时间，如果在此时间内没有收到响应，则连接将被关闭；配置示例：1s、1000ms
	DialKeepAliveTimeout string

	// 拨号重试策略: 默认为空：表示默认不重试；1、2、3...表示重试多少次；always：表示一直重试
	DialRetry string

	// 最大呼叫：发送MSG大小是客户端请求发送的字节限制
	MaxCallSendMsgSize int

	// 最大调用recv MSG大小是客户端响应接收限制
	MaxCallRecvMsgSize int

	// 当设置拒绝旧集群时，将拒绝在过时的集群上创建客户端
	RejectOldCluster bool

	// 设置允许无流时将允许客户端发送keepalive ping到服务器没有任何活动流rp cs
	PermitWithoutStream bool
}
