package redis

import (
	"errors"
	redisprom "github.com/globocom/go-redis-prometheus"
	redisotel "github.com/go-redis/redis/extra/redisotel/v8"
	goredis "github.com/go-redis/redis/v8"
	"github.com/simonalong/gole-boot/constants"
	_ "github.com/simonalong/gole-boot/otel"
	"github.com/simonalong/gole/bean"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/logger"
	baseTime "github.com/simonalong/gole/time"
	"sync"

	//"github.com/simonalong/gole-boot/tracing"
	//"github.com/simonalong/gole-boot/tracing"
	"time"
)

var Hooks []goredis.Hook
var initLock sync.Mutex

type ConfigError struct {
	ErrMsg string
}

func (error *ConfigError) Error() string {
	return error.ErrMsg
}

func init() {
	config.Load()

	if config.Loaded && config.GetValueBoolDefault("gole.redis.enable", false) {
		err := config.GetValueObject("gole.redis", &Cfg)
		if err != nil {
			logger.Warn("读取redis配置异常")
			return
		}
	}
	Hooks = []goredis.Hook{}
}

func GetClient() (goredis.UniversalClient, error) {
	if bean.GetBean(constants.BeanNameRedis) != nil {
		return bean.GetBean(constants.BeanNameRedis).(goredis.UniversalClient), nil
	}
	initLock.Lock()
	defer initLock.Unlock()
	if bean.GetBean(constants.BeanNameRedis) != nil {
		return bean.GetBean(constants.BeanNameRedis).(goredis.UniversalClient), nil
	}
	rdbClient, err := NewClient()
	if err != nil {
		return nil, err
	}
	bean.AddBean(constants.BeanNameRedis, rdbClient)
	return rdbClient, nil
}

func NewClient() (goredis.UniversalClient, error) {
	if !config.GetValueBoolDefault("gole.redis.enable", false) {
		logger.Error("redis配置开关为关闭，请开启")
		return nil, errors.New("redis配置开关为关闭，请开启")
	}
	var rdbClient goredis.UniversalClient
	if Cfg.Sentinel.Master != "" {
		rdbClient = goredis.NewFailoverClient(getSentinelConfig())
	} else if len(Cfg.Cluster.Addrs) != 0 {
		rdbClient = goredis.NewClusterClient(getClusterConfig())
	} else {
		rdbClient = goredis.NewClient(getStandaloneConfig())
	}

	for _, hook := range Hooks {
		rdbClient.AddHook(hook)
	}

	// 支持opentelemetry埋点
	if config.GetValueBoolDefault("gole.opentelemetry.enable", true) {
		rdbClient.AddHook(redisotel.NewTracingHook())
	}

	if config.GetValueBoolDefault("gole.meter.redis.enable", false) {
		var applicationName string
		if val := config.GoleCfg.Application.Name; val != "" {
			applicationName = val
		} else {
			logger.Fatalf("gole.application.name 不可为空")
		}
		rdbClient.AddHook(redisprom.NewHook(
			redisprom.WithInstanceName(applicationName),
			redisprom.WithNamespace("base_boot"),
			redisprom.WithDurationBuckets([]float64{.001, .005, .01}),
		))
	}

	bean.AddBean(constants.BeanNameRedis, &rdbClient)
	return rdbClient, nil
}

func AddRedisHook(hook goredis.Hook) {
	Hooks = append(Hooks, hook)
}

func getStandaloneConfig() *goredis.Options {
	addr := "127.0.0.1:6379"
	if Cfg.Standalone.Addr != "" {
		addr = Cfg.Standalone.Addr
	}

	redisConfig := &goredis.Options{
		Addr: addr,

		DB:       Cfg.Standalone.Database,
		Network:  Cfg.Standalone.Network,
		Username: Cfg.Username,
		Password: Cfg.Password,

		MaxRetries:      Cfg.MaxRetries,
		MinRetryBackoff: baseTime.NumToTimeDuration(Cfg.MinRetryBackoff, time.Millisecond),
		MaxRetryBackoff: baseTime.NumToTimeDuration(Cfg.MaxRetryBackoff, time.Millisecond),

		DialTimeout:  baseTime.NumToTimeDuration(Cfg.DialTimeout, time.Millisecond),
		ReadTimeout:  baseTime.NumToTimeDuration(Cfg.ReadTimeout, time.Millisecond),
		WriteTimeout: baseTime.NumToTimeDuration(Cfg.WriteTimeout, time.Millisecond),

		PoolFIFO:           Cfg.PoolFIFO,
		PoolSize:           Cfg.PoolSize,
		MinIdleConns:       Cfg.MinIdleConns,
		MaxConnAge:         baseTime.NumToTimeDuration(Cfg.MaxConnAge, time.Millisecond),
		PoolTimeout:        baseTime.NumToTimeDuration(Cfg.PoolTimeout, time.Millisecond),
		IdleTimeout:        baseTime.NumToTimeDuration(Cfg.IdleTimeout, time.Millisecond),
		IdleCheckFrequency: baseTime.NumToTimeDuration(Cfg.IdleCheckFrequency, time.Millisecond),
	}

	// -------- 命令执行失败配置 --------
	if config.GetValueString("gole.redis.max-retries") == "" {
		// # 命令执行失败时候，最大重试次数，默认3次，-1（不是0）则不重试
		redisConfig.MaxRetries = 3
	}

	if config.GetValueString("gole.redis.min-retry-backoff") == "" {
		// #（单位毫秒） 命令执行失败时候，每次重试的最小回退时间，默认8毫秒，-1则禁止回退
		redisConfig.MinRetryBackoff = 8 * time.Millisecond
	}

	if config.GetValueString("gole.redis.max-retry-backoff") == "" {
		// #（单位毫秒） 命令执行失败时候，每次重试的最小回退时间，默认8毫秒，-1则禁止回退
		redisConfig.MinRetryBackoff = 512 * time.Millisecond
	}

	// -------- 超时配置 --------
	if config.GetValueString("gole.redis.dial-timeout") == "" {
		// # （单位毫秒）超时：创建新链接的拨号超时时间，默认15秒
		redisConfig.DialTimeout = 15 * time.Second
	}

	if config.GetValueString("gole.redis.read-timeout") == "" {
		// # （单位毫秒）超时：读超时，默认3秒，使用-1，使用-1则表示无超时，0的话是表示默认3秒
		redisConfig.ReadTimeout = 3 * time.Second
	}

	if config.GetValueString("gole.redis.write-timeout") == "" {
		// # （单位毫秒）超时：写超时，默认是读超时3秒，使用-1，使用-1则表示无超时，0的话是表示默认3秒
		redisConfig.WriteTimeout = 3 * time.Second
	}

	// -------- 连接池相关配置 --------
	if config.GetValueString("gole.redis.pool-fifo") == "" {
		// # 连接池类型：fifo：true;lifo：false;和lifo相比，fifo开销更高
		redisConfig.PoolFIFO = false
	}

	if config.GetValueString("gole.redis.pool-size") == "" {
		// # 最大连接池大小：默认每个cpu核是10个连接，cpu核数可以根据函数runtime.GOMAXPROCS来配置，默认是runtime.NumCpu
		redisConfig.PoolSize = 10
	}

	if config.GetValueString("gole.redis.min-idle-conns") == "" {
		// # 最小空闲连接数
		redisConfig.MinIdleConns = 10
	}

	if config.GetValueString("gole.redis.max-conn-age") == "" {
		// #（单位毫秒） 连接存活时长，默认不关闭
		redisConfig.MaxConnAge = 12 * 30 * 24 * time.Hour
	}

	if config.GetValueString("gole.redis.pool-timeout") == "" {
		// #（单位毫秒）获取链接池中的链接都在忙，则等待对应的时间，默认读超时+1秒
		redisConfig.PoolTimeout = time.Second
	}

	if config.GetValueString("gole.redis.idle-timeout") == "" {
		// #（单位毫秒）空闲链接时间，超时则关闭，注意：该时间要小于服务端的超时时间，否则会出现拿到的链接失效问题，默认5分钟，-1表示禁用超时检查
		redisConfig.IdleTimeout = 5 * time.Minute
	}

	if config.GetValueString("gole.redis.idle-check-frequency") == "" {
		// #（单位毫秒）空闲链接核查频率，默认1分钟。-1禁止空闲链接核查，即使配置了IdleTime也不行
		redisConfig.IdleCheckFrequency = time.Minute
	}
	return redisConfig
}

func getSentinelConfig() *goredis.FailoverOptions {
	redisConfig := &goredis.FailoverOptions{
		SentinelAddrs: Cfg.Sentinel.Addrs,
		MasterName:    Cfg.Sentinel.Master,

		DB:               Cfg.Sentinel.Database,
		Username:         Cfg.Username,
		Password:         Cfg.Password,
		SentinelUsername: Cfg.Sentinel.SentinelUser,
		SentinelPassword: Cfg.Sentinel.SentinelPassword,

		MaxRetries:      Cfg.MaxRetries,
		MinRetryBackoff: baseTime.NumToTimeDuration(Cfg.MinRetryBackoff, time.Millisecond),
		MaxRetryBackoff: baseTime.NumToTimeDuration(Cfg.MaxRetryBackoff, time.Millisecond),

		DialTimeout:  baseTime.NumToTimeDuration(Cfg.DialTimeout, time.Millisecond),
		ReadTimeout:  baseTime.NumToTimeDuration(Cfg.ReadTimeout, time.Millisecond),
		WriteTimeout: baseTime.NumToTimeDuration(Cfg.WriteTimeout, time.Millisecond),

		PoolFIFO:           Cfg.PoolFIFO,
		PoolSize:           Cfg.PoolSize,
		MinIdleConns:       Cfg.MinIdleConns,
		MaxConnAge:         baseTime.NumToTimeDuration(Cfg.MaxConnAge, time.Millisecond),
		PoolTimeout:        baseTime.NumToTimeDuration(Cfg.PoolTimeout, time.Millisecond),
		IdleTimeout:        baseTime.NumToTimeDuration(Cfg.IdleTimeout, time.Millisecond),
		IdleCheckFrequency: baseTime.NumToTimeDuration(Cfg.IdleCheckFrequency, time.Millisecond),
	}

	// -------- 命令执行失败配置 --------
	if config.GetValueString("gole.redis.max-retries") == "" {
		// # 命令执行失败时候，最大重试次数，默认3次，-1（不是0）则不重试
		redisConfig.MaxRetries = 3
	}

	if config.GetValueString("gole.redis.min-retry-backoff") == "" {
		// #（单位毫秒） 命令执行失败时候，每次重试的最小回退时间，默认8毫秒，-1则禁止回退
		redisConfig.MinRetryBackoff = 8 * time.Millisecond
	}

	if config.GetValueString("gole.redis.max-retry-backoff") == "" {
		// #（单位毫秒） 命令执行失败时候，每次重试的最小回退时间，默认8毫秒，-1则禁止回退
		redisConfig.MinRetryBackoff = 512 * time.Millisecond
	}

	// -------- 超时配置 --------
	if config.GetValueString("gole.redis.dial-timeout") == "" {
		// # （单位毫秒）超时：创建新链接的拨号超时时间，默认15秒
		redisConfig.DialTimeout = 15 * time.Second
	}

	if config.GetValueString("gole.redis.read-timeout") == "" {
		// # （单位毫秒）超时：读超时，默认3秒，使用-1，使用-1则表示无超时，0的话是表示默认3秒
		redisConfig.ReadTimeout = 3 * time.Second
	}

	if config.GetValueString("gole.redis.write-timeout") == "" {
		// # （单位毫秒）超时：写超时，默认是读超时3秒，使用-1，使用-1则表示无超时，0的话是表示默认3秒
		redisConfig.WriteTimeout = 3 * time.Second
	}

	// -------- 连接池相关配置 --------
	if config.GetValueString("gole.redis.pool-fifo") == "" {
		// # 连接池类型：fifo：true;lifo：false;和lifo相比，fifo开销更高
		redisConfig.PoolFIFO = false
	}

	if config.GetValueString("gole.redis.pool-size") == "" {
		// # 最大连接池大小：默认每个cpu核是10个连接，cpu核数可以根据函数runtime.GOMAXPROCS来配置，默认是runtime.NumCpu
		redisConfig.PoolSize = 10
	}

	if config.GetValueString("gole.redis.min-idle-conns") == "" {
		// # 最小空闲连接数
		redisConfig.MinIdleConns = 10
	}

	if config.GetValueString("gole.redis.max-conn-age") == "" {
		// #（单位毫秒） 连接存活时长，默认不关闭
		redisConfig.MaxConnAge = 12 * 30 * 24 * time.Hour
	}

	if config.GetValueString("gole.redis.pool-timeout") == "" {
		// #（单位毫秒）获取链接池中的链接都在忙，则等待对应的时间，默认读超时+1秒
		redisConfig.PoolTimeout = time.Second
	}

	if config.GetValueString("gole.redis.idle-timeout") == "" {
		// #（单位毫秒）空闲链接时间，超时则关闭，注意：该时间要小于服务端的超时时间，否则会出现拿到的链接失效问题，默认5分钟，-1表示禁用超时检查
		redisConfig.IdleTimeout = 5 * time.Minute
	}

	if config.GetValueString("gole.redis.idle-check-frequency") == "" {
		// #（单位毫秒）空闲链接核查频率，默认1分钟。-1禁止空闲链接核查，即使配置了IdleTime也不行
		redisConfig.IdleCheckFrequency = time.Minute
	}
	return redisConfig
}

func getClusterConfig() *goredis.ClusterOptions {
	if len(Cfg.Cluster.Addrs) == 0 {
		Cfg.Cluster.Addrs = []string{"127.0.0.1:6379"}
	}

	redisConfig := &goredis.ClusterOptions{
		Addrs: Cfg.Cluster.Addrs,

		Username: Cfg.Username,
		Password: Cfg.Password,

		MaxRedirects:   Cfg.Cluster.MaxRedirects,
		ReadOnly:       Cfg.Cluster.ReadOnly,
		RouteByLatency: Cfg.Cluster.RouteByLatency,
		RouteRandomly:  Cfg.Cluster.RouteRandomly,

		MaxRetries:      Cfg.MaxRetries,
		MinRetryBackoff: baseTime.NumToTimeDuration(Cfg.MinRetryBackoff, time.Millisecond),
		MaxRetryBackoff: baseTime.NumToTimeDuration(Cfg.MaxRetryBackoff, time.Millisecond),

		DialTimeout:  baseTime.NumToTimeDuration(Cfg.DialTimeout, time.Millisecond),
		ReadTimeout:  baseTime.NumToTimeDuration(Cfg.ReadTimeout, time.Millisecond),
		WriteTimeout: baseTime.NumToTimeDuration(Cfg.WriteTimeout, time.Millisecond),
		PoolFIFO:     Cfg.PoolFIFO,
		PoolSize:     Cfg.PoolSize,
		MinIdleConns: Cfg.MinIdleConns,

		MaxConnAge:         baseTime.NumToTimeDuration(Cfg.MaxConnAge, time.Millisecond),
		PoolTimeout:        baseTime.NumToTimeDuration(Cfg.PoolTimeout, time.Millisecond),
		IdleTimeout:        baseTime.NumToTimeDuration(Cfg.IdleTimeout, time.Millisecond),
		IdleCheckFrequency: baseTime.NumToTimeDuration(Cfg.IdleCheckFrequency, time.Millisecond),
	}

	// -------- 命令执行失败配置 --------
	if config.GetValueString("gole.redis.max-retries") == "" {
		// # 命令执行失败时候，最大重试次数，默认3次，-1（不是0）则不重试
		redisConfig.MaxRetries = 3
	}

	if config.GetValueString("gole.redis.min-retry-backoff") == "" {
		// #（单位毫秒） 命令执行失败时候，每次重试的最小回退时间，默认8毫秒，-1则禁止回退
		redisConfig.MinRetryBackoff = 8 * time.Millisecond
	}

	if config.GetValueString("gole.redis.max-retry-backoff") == "" {
		// #（单位毫秒） 命令执行失败时候，每次重试的最小回退时间，默认8毫秒，-1则禁止回退
		redisConfig.MinRetryBackoff = 512 * time.Millisecond
	}

	// -------- 超时配置 --------
	if config.GetValueString("gole.redis.dial-timeout") == "" {
		// # （单位毫秒）超时：创建新链接的拨号超时时间，默认15秒
		redisConfig.DialTimeout = 15 * time.Second
	}

	if config.GetValueString("gole.redis.read-timeout") == "" {
		// # （单位毫秒）超时：读超时，默认3秒，使用-1，使用-1则表示无超时，0的话是表示默认3秒
		redisConfig.ReadTimeout = 3 * time.Second
	}

	if config.GetValueString("gole.redis.write-timeout") == "" {
		// # （单位毫秒）超时：写超时，默认是读超时3秒，使用-1，使用-1则表示无超时，0的话是表示默认3秒
		redisConfig.WriteTimeout = 3 * time.Second
	}

	// -------- 连接池相关配置 --------
	if config.GetValueString("gole.redis.pool-fifo") == "" {
		// # 连接池类型：fifo：true;lifo：false;和lifo相比，fifo开销更高
		redisConfig.PoolFIFO = false
	}

	if config.GetValueString("gole.redis.pool-size") == "" {
		// # 最大连接池大小：默认每个cpu核是10个连接，cpu核数可以根据函数runtime.GOMAXPROCS来配置，默认是runtime.NumCpu
		redisConfig.PoolSize = 10
	}

	if config.GetValueString("gole.redis.min-idle-conns") == "" {
		// # 最小空闲连接数
		redisConfig.MinIdleConns = 10
	}

	if config.GetValueString("gole.redis.max-conn-age") == "" {
		// #（单位毫秒） 连接存活时长，默认不关闭
		redisConfig.MaxConnAge = 12 * 30 * 24 * time.Hour
	}

	if config.GetValueString("gole.redis.pool-timeout") == "" {
		// #（单位毫秒）获取链接池中的链接都在忙，则等待对应的时间，默认读超时+1秒
		redisConfig.PoolTimeout = time.Second
	}

	if config.GetValueString("gole.redis.idle-timeout") == "" {
		// #（单位毫秒）空闲链接时间，超时则关闭，注意：该时间要小于服务端的超时时间，否则会出现拿到的链接失效问题，默认5分钟，-1表示禁用超时检查
		redisConfig.IdleTimeout = 5 * time.Minute
	}

	if config.GetValueString("gole.redis.idle-check-frequency") == "" {
		// #（单位毫秒）空闲链接核查频率，默认1分钟。-1禁止空闲链接核查，即使配置了IdleTime也不行
		redisConfig.IdleCheckFrequency = time.Minute
	}
	return redisConfig
}
