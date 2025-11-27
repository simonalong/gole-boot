package test

import (
	"github.com/simonalong/gole-boot/extend/redis"
	"github.com/simonalong/gole/logger"
	"testing"
	"time"
)

// 分布式锁：阻塞式获取分布式锁
func TestLock(t *testing.T) {
	redis.Lock("test", 15*time.Second, func() {
		logger.Info("加锁成功，开始执行业务")
		// 业务逻辑
		time.Sleep(10 * time.Second)
		logger.Info("业务执行完毕")
	})
}

func TestLock_other(t *testing.T) {
	redis.Lock("test", 15*time.Second, func() {
		logger.Info("加锁成功，开始执行业务")
		// 业务逻辑
		time.Sleep(10 * time.Second)
		logger.Info("业务执行完毕")
	})
}

// 分布式锁：非阻塞式获取分布式锁
func TestTryLock(t *testing.T) {
	redis.TryLock("test", 15*time.Second, func() {
		logger.Info("加锁成功，开始执行业务")
		// 业务逻辑
		time.Sleep(10 * time.Second)
		logger.Info("业务执行完毕")
	})
}

func TestTryLock_other(t *testing.T) {
	redis.TryLock("test", 15*time.Second, func() {
		logger.Info("加锁成功，开始执行业务")
		time.Sleep(10 * time.Second)
		logger.Info("业务执行完毕")
	})
}
