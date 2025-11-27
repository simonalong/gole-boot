package job

import "time"

type Config struct {
	Enable          bool          `json:"enable"`
	ServerAddress   string        `json:"serverAddress"`   // 分布式任务服务地址，示例：http://cbb-mid-srv-job:18080
	Timeout         time.Duration `json:"timeout"`         // 接口超时时间，默认5秒
	ExecutorName    string        `json:"executorName"`    // 执行器名称
	ExecutorAddress string        `json:"executorAddress"` // 执行器的服务地址和域名，示例：<host>:<port>
}
