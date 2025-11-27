package redis

var Cfg RedisConfig

// ---------------------------- gole.redis ----------------------------

// RedisConfig gole.redis前缀
type RedisConfig struct {
	Password string
	Username string

	// 单节点
	Standalone RedisStandaloneConfig
	// 哨兵
	Sentinel RedisSentinelConfig
	// 集群
	Cluster RedisClusterConfig

	// ----- 命令执行失败配置 -----
	// 命令执行失败时候，最大重试次数，默认3次，-1（不是0）则不重试
	MaxRetries int
	// （单位毫秒） 命令执行失败时候，每次重试的最小回退时间，默认8毫秒，-1则禁止回退
	MinRetryBackoff int
	// （单位毫秒）命令执行失败时候，每次重试的最大回退时间，默认512毫秒，-1则禁止回退
	MaxRetryBackoff int

	// ----- 超时配置 -----
	// （单位毫秒）超时：创建新链接的拨号超时时间，默认15秒
	DialTimeout int
	// （单位毫秒）超时：读超时，默认3秒，使用-1，使用-1则表示无超时，0的话是表示默认3秒
	ReadTimeout int
	// （单位毫秒）超时：写超时，默认是读超时3秒，使用-1，使用-1则表示无超时，0的话是表示默认3秒
	WriteTimeout int

	// ----- 连接池相关配置 -----
	// 连接池类型：fifo：true;lifo：false;和lifo相比，fifo开销更高
	PoolFIFO bool
	// 最大连接池大小：默认每个cpu核是10个连接，cpu核数可以根据函数runtime.GOMAXPROCS来配置，默认是runtime.NumCpu
	PoolSize int
	// 最小空闲连接数
	MinIdleConns int
	// （单位毫秒） 连接存活时长，默认不关闭
	MaxConnAge int
	// （单位毫秒）获取链接池中的链接都在忙，则等待对应的时间，默认读超时+1秒
	PoolTimeout int
	// （单位毫秒）空闲链接时间，超时则关闭，注意：该时间要小于服务端的超时时间，否则会出现拿到的链接失效问题，默认5分钟，-1表示禁用超时检查
	IdleTimeout int
	// （单位毫秒）空闲链接核查频率，默认1分钟。-1禁止空闲链接核查，即使配置了IdleTime也不行
	IdleCheckFrequency int
}

// RedisStandaloneConfig gole.redis.standalone
type RedisStandaloneConfig struct {
	Addr     string
	Database int
	// 网络类型，tcp或者unix，默认tcp
	Network  string `match:"value={tcp, unix}"  errMsg:"network值不合法，只可为两个值：tcp和unix"`
	ReadOnly bool
}

// RedisSentinelConfig gole.redis.sentinel
type RedisSentinelConfig struct {
	// 哨兵的集群名字
	Master string
	// 哨兵节点地址
	Addrs []string
	// 数据库节点
	Database int
	// 哨兵用户
	SentinelUser string
	// 哨兵密码
	SentinelPassword string
	// 将所有命令路由到从属只读节点。
	SlaveOnly bool
}

type RedisClusterConfig struct {
	// 节点地址
	Addrs []string
	// 最大重定向次数
	MaxRedirects int
	// 开启从节点的只读功能
	ReadOnly bool
	// 允许将只读命令路由到最近的主节点或从节点，它会自动启用 ReadOnly
	RouteByLatency bool
	// 允许将只读命令路由到随机的主节点或从节点，它会自动启用 ReadOnly
	RouteRandomly bool
}
