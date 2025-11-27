package nats

import (
	"sync"
	"testing"

	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/logger"
	"github.com/stretchr/testify/assert"
)

// 使用本例子之前请使用 ./deploy/xxx 文件夹中的部署进行启动即可

// 使用环境变量：gole.profiles.active=all
func TestNatsConfigWithAll(t *testing.T) {
	nc, err := baseNats.New()
	if err != nil {
		logger.Fatal(err)
		return
	}
	defer nc.Close()

	assert.Equal(t, "nats://127.0.0.1:4222", baseNats.CfgOfNats.Url)
	assert.Equal(t, "xx-demo-service", baseNats.CfgOfNats.Name)

	// --------------- 连接认证 ---------------
	assert.Equal(t, "admin", baseNats.CfgOfNats.UserName)
	assert.Equal(t, "admin-demo123@xxxx.com", baseNats.CfgOfNats.Password)
	assert.Equal(t, "xxxxxxxxxxxxxxx", baseNats.CfgOfNats.Token)
	assert.Equal(t, "./nkeys/seed.txt", baseNats.CfgOfNats.NkSeedFile)
	//assert.Equal(t, true, baseNats.CfgOfNats.Secure)
	//assert.Equal(t, true, baseNats.CfgOfNats.TLSHandshakeFirst)

	// --------------- 连接配置 ---------------
	assert.Equal(t, true, baseNats.CfgOfNats.AllowReconnect)
	assert.Equal(t, 100, baseNats.CfgOfNats.MaxReconnect)
	assert.Equal(t, "12s", baseNats.CfgOfNats.ReconnectWait.String())
	assert.Equal(t, "100ms", baseNats.CfgOfNats.ReconnectJitter.String())
	assert.Equal(t, "1s", baseNats.CfgOfNats.ReconnectJitterTLS.String())
	assert.Equal(t, "2s", baseNats.CfgOfNats.Timeout.String())
	assert.Equal(t, "30s", baseNats.CfgOfNats.DrainTimeout.String())
	assert.Equal(t, "1m0s", baseNats.CfgOfNats.FlusherTimeout.String())
	assert.Equal(t, "2m0s", baseNats.CfgOfNats.PingInterval.String())
	assert.Equal(t, 2, baseNats.CfgOfNats.MaxPingsOut)
	assert.Equal(t, 8388608, baseNats.CfgOfNats.ReconnectBufSize)
	assert.Equal(t, 65536, baseNats.CfgOfNats.SubChanLen)
	assert.Equal(t, false, baseNats.CfgOfNats.RetryOnFailedConnect)

	// --------------- websocket 配置 ---------------
	assert.Equal(t, false, baseNats.CfgOfNats.Compression)
	assert.Equal(t, "ws://xxx/xxx/xx", baseNats.CfgOfNats.ProxyPath)

	// --------------- 其他配置 ---------------
	assert.Equal(t, false, baseNats.CfgOfNats.NoRandomize)
	assert.Equal(t, false, baseNats.CfgOfNats.NoEcho)
	assert.Equal(t, false, baseNats.CfgOfNats.Verbose)
	assert.Equal(t, false, baseNats.CfgOfNats.Pedantic)
	assert.Equal(t, false, baseNats.CfgOfNats.UseOldRequestStyle)
	assert.Equal(t, false, baseNats.CfgOfNats.NoCallbacksAfterClientClose)
	assert.Equal(t, "_api_demo", baseNats.CfgOfNats.InboxPrefix)
	assert.Equal(t, false, baseNats.CfgOfNats.IgnoreAuthErrorAbort)
	assert.Equal(t, false, baseNats.CfgOfNats.SkipHostLookup)
}

// 使用如下的环境变量进而测试不同的认证情况
// 使用环境变量：gole.profiles.active=user
// 使用环境变量：gole.profiles.active=token
// 使用环境变量：gole.profiles.active=creds
// 使用环境变量：gole.profiles.active=nkey

// 使用环境变量：gole.profiles.active=cluster
func TestNatsConnectWithUser(t *testing.T) {
	nc, err := baseNats.New()
	if err != nil {
		logger.Fatal(err)
		return
	}
	defer nc.Close()

	logger.Infof("ok")

	count := sync.WaitGroup{}
	count.Add(1)
	_, err = nc.Subscribe("test.connect", func(msg *baseNats.MsgOfNats) {
		assert.Equal(t, string(msg.Data), "hello world")
		count.Done()
	})

	err = nc.Publish("test.connect", []byte("hello world"))
	count.Wait()
}
