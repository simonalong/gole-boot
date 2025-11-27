package nats

import (
	"sync"
	"testing"
	"time"

	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/logger"
	"github.com/stretchr/testify/assert"
)

// 使用本例子之前请使用 ./deploy/xxx 文件夹中的部署进行启动即可
// 使用环境变量：gole.profiles.active=user
func TestNatsConnectWithUser(t *testing.T) {
	nc, err := baseNats.New()
	if err != nil {
		logger.Fatal(err)
		return
	}
	defer nc.Close()

	count := sync.WaitGroup{}
	count.Add(1)
	_, err = nc.Subscribe("test.connect", func(msg *baseNats.MsgOfNats) {
		assert.Equal(t, string(msg.Data), "hello world")
		count.Done()
	})

	err = nc.Publish("test.connect", []byte("hello world"))
	count.Wait()

	// 添加时间，等待数据上报完毕
	time.Sleep(5 * time.Second)
}
