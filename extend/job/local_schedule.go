package job

import (
	"github.com/robfig/cron/v3"
	"github.com/simonalong/gole-boot/extend/redis"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/logger"
	"time"
)

var CronParser cron.Parser

func init() {
	CronParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	return
}

func ParseCron(cronStr string) (cron.Schedule, error) {
	return CronParser.Parse(cronStr)
}

// ScheduleCron 定时任务，支持cron表达式
// bizName 业务名称，用于业务唯一识别
// cron 表达式
// fun 业务函数
func ScheduleCron(bizName, cron string, bizFun func()) {
	checkRedis()
	schedule, err := ParseCron(cron)
	if err != nil {
		logger.Errorf("解析cron表达式失败, bizName=%v, cron=%v, err=%v", bizName, cron, err)
		return
	}
	// 执行定时任务
	go func() {
		for {
			currentTime := time.Now()
			nextTime := schedule.Next(currentTime)

			duration := nextTime.Sub(currentTime)
			time.Sleep(duration)

			// 使用分布式锁执行调度
			go redis.TryLock(bizName, duration, bizFun)
		}
	}()
}

// ScheduleFixRate 每隔一段时间执行一次，无论前次任务是否完成
// bizName 业务名称，用于业务唯一识别
// duration 固定频率，最小单位为秒
// fun 业务函数
func ScheduleFixRate(bizName string, duration time.Duration, bizFun func()) {
	checkRedis()
	if duration < time.Second {
		logger.Fatalf("固定频率不能小于1秒, bizName=%v, duration=%v", bizName, duration)
		return
	}
	first := true
	// 执行定时任务
	go func() {
		for {
			if first {
				// 首次调整整秒执行
				time.Sleep(time.Duration(1000-time.Now().UnixMilli()%1000) * time.Millisecond)
				first = false
			}

			// 使用整秒执行
			time.Sleep(duration - time.Second + time.Duration(1000-time.Now().UnixMilli()%1000)*time.Millisecond)

			// 使用分布式锁执行调度
			go redis.TryLock(bizName, duration, bizFun)
		}
	}()
}

func checkRedis() {
	if config.Loaded {
		// 核查redis是否开启，没有开启，则不启动
		if !config.GetValueBoolDefault("gole.redis.enable", false) {
			logger.Group("job").Fatalf("redis 服务未开启，无法启动cbb-job的客户端本地调度功能，请先开启：gole.redis.enable")
			return
		}
	}
}
