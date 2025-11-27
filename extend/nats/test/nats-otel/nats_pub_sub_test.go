/*
* 使用本例子之前请使用 ./deploy/xxx 文件夹中的部署进行启动即可
* 本文件全部使用：环境变量：gole.profiles.active=user
 */
package nats

import (
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/logger"
	"github.com/stretchr/testify/assert"
)

// 使用环境变量：gole.profiles.active=user
func TestNatsPubMsg(t *testing.T) {
	nc, err := baseNats.New()
	if err != nil {
		logger.Fatal(err)
		return
	}
	defer nc.Close()

	count := sync.WaitGroup{}
	count.Add(1)
	_, err = nc.Subscribe("test.pub.sub1", func(msg *baseNats.MsgOfNats) {
		assert.Equal(t, string(msg.Data), "hello world")
		count.Done()
	})

	pMsg := &nats.Msg{
		Subject: "test.pub.sub1",
		Data:    []byte("hello world"),
	}
	err = nc.PublishMsg(pMsg)
	count.Wait()

	// 添加时间，等待数据上报完毕
	time.Sleep(5 * time.Second)
}

func TestNatsPubMsg2(t *testing.T) {
	nc, err := baseNats.New()
	if err != nil {
		logger.Fatal(err)
		return
	}
	defer nc.Close()

	count := sync.WaitGroup{}
	count.Add(1)
	_, err = nc.Subscribe("test.pub.sub2", func(msg *baseNats.MsgOfNats) {
		assert.Equal(t, msg.Header.Get("head1"), "xxx")
		assert.Equal(t, string(msg.Data), "hello world")
		count.Done()
	})

	pMsg := &nats.Msg{
		Subject: "test.pub.sub2",
		Data:    []byte("hello world"),
		Header: map[string][]string{
			"head1": {"xxx"},
		},
	}
	err = nc.PublishMsg(pMsg)
	count.Wait()

	// 添加时间，等待数据上报完毕
	time.Sleep(10 * time.Second)
}
