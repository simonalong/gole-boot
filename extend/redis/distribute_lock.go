package redis

import (
	"context"
	"errors"
	"fmt"
	goredisV8 "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/simonalong/gole-boot/constants"
	"github.com/simonalong/gole/bean"
	"github.com/simonalong/gole/logger"
	"time"
)

type DistributedLock struct {
	mutex     *redsync.Mutex
	stopRenew chan struct{} // 停止续期信号
	expiry    time.Duration // 锁过期时间
}

// Lock 分布式锁：获取锁失败则阻塞业务，等待锁释放，等待时长为expiry；如果还是失败，则返回失败
func Lock(lockKey string, expiry time.Duration, bizFun func()) error {
	redisClient := getRedisClient()
	dsLock := newDistributedLock(redisClient, lockKey, expiry, 2, expiry/2)
	if dsLock == nil {
		logger.Errorf("创建分布式锁失败, key=%v", lockKey)
		return errors.New(fmt.Sprintf("创建分布式锁失败, key=%v", lockKey))
	}
	ctx := context.Background()
	err := dsLock.lock(ctx)
	if err != nil {
		return err
	}
	defer func(dsLock *DistributedLock) {
		_, err := dsLock.Unlock()
		if err != nil {
			logger.Errorf("解锁失败, key=%v, err=%v", lockKey, err)
		}
	}(dsLock)

	// 执行业务
	bizFun()
	return nil
}

// TryLock 分布式锁：获取锁失败则不做任何处理
func TryLock(lockKey string, expiry time.Duration, bizFun func()) bool {
	redisClient := getRedisClient()
	dsLock := newDistributedLock(redisClient, lockKey, expiry, 1, 0)
	if dsLock == nil {
		logger.Errorf("创建分布式锁失败, key=%v", lockKey)
		return false
	}
	ctx := context.Background()
	lockSuccess, err := dsLock.TryLock(ctx)
	if !lockSuccess || err != nil {
		return false
	}
	defer func(dsLock *DistributedLock) {
		_, err := dsLock.Unlock()
		if err != nil {
			logger.Errorf("解锁失败, key=%v, err=%v", lockKey, err)
		}
	}(dsLock)

	// 执行业务
	bizFun()
	return true
}

func getRedisClient() goredisV8.UniversalClient {
	obj := bean.GetBean(constants.BeanNameRedis)
	if obj != nil {
		return *obj.(*goredisV8.UniversalClient)
	}
	_redisClient, err := NewClient()
	if err != nil {
		logger.Fatalf("初始化redis客户端失败, err=%v", err)
		return nil
	}
	return _redisClient
}

// newDistributedLock 创建分布式锁实例
func newDistributedLock(redisClient goredisV8.UniversalClient, lockKey string, expiry time.Duration, tries int, delay time.Duration) *DistributedLock {
	if redisClient == nil {
		logger.Fatalf("redis示例为空，请先配置")
		return nil
	}
	pool := goredis.NewPool(redisClient)
	rs := redsync.New(pool)

	mutex := rs.NewMutex(
		lockKey,
		redsync.WithExpiry(expiry),
		redsync.WithTries(tries),
		redsync.WithRetryDelay(delay),
	)

	return &DistributedLock{
		mutex:     mutex,
		stopRenew: make(chan struct{}),
		expiry:    expiry,
	}
}

// Lock 阻塞式获取锁（自动续期）
func (dl *DistributedLock) lock(ctx context.Context) error {
	if err := dl.mutex.LockContext(ctx); err != nil {
		return err
	}
	// 启动自动续期
	dl.startRenewal()
	return nil
}

// TryLock 非阻塞尝试获取锁
func (dl *DistributedLock) TryLock(ctx context.Context) (bool, error) {
	err := dl.mutex.TryLockContext(ctx)
	if err != nil {
		return false, err
	}
	// 启动自动续期
	dl.startRenewal()
	return true, nil
}

// Unlock 释放锁并停止续期
func (dl *DistributedLock) Unlock() (bool, error) {
	close(dl.stopRenew) // 停止续期协程
	return dl.mutex.Unlock()
}

// 后台协程自动续期
func (dl *DistributedLock) startRenewal() {
	go func() {
		ticker := time.NewTicker(dl.expiry / 3) // 续期间隔 = TTL/3
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if _, err := dl.mutex.Extend(); err != nil {
					return // 续期失败则退出
				}
			case <-dl.stopRenew:
				return // 收到停止信号
			}
		}
	}()
}
