/*
* 使用本例子之前请使用 ./deploy/xxx 文件夹中的部署进行启动即可
* 本文件全部使用：环境变量：gole.profiles.active=user
 */
package nats

import (
	"github.com/magiconair/properties/assert"
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/logger"
	"log"
	"sync"
	"testing"
	"time"
)

func TestNatsQueueGroup1(t *testing.T) {
	nc, err := baseNats.New()
	if err != nil {
		logger.Fatal(err)
		return
	}
	defer nc.Close()

	count := sync.WaitGroup{}
	count.Add(3)
	_, err = nc.QueueSubscribe("test.queue1.req", "queue1", func(m *baseNats.MsgOfNats) {
		assert.Equal(t, string(m.Data), "ni hao")
		err := m.Respond([]byte("get it"))
		if err != nil {
			logger.Fatal(err)
			return
		}
		count.Done()
	})

	_, err = nc.QueueSubscribe("test.queue1.req", "queue1", func(m *baseNats.MsgOfNats) {
		assert.Equal(t, string(m.Data), "ni hao")
		err := m.Respond([]byte("get it"))
		if err != nil {
			logger.Fatal(err)
			return
		}
		count.Done()
	})

	_, err = nc.QueueSubscribe("test.queue1.req", "queue1", func(m *baseNats.MsgOfNats) {
		assert.Equal(t, string(m.Data), "ni hao")
		err := m.Respond([]byte("get it"))
		if err != nil {
			logger.Fatal(err)
			return
		}
		count.Done()
	})

	msg, err := nc.Request("test.queue1.req", []byte("ni hao"), time.Second)
	msg, err = nc.Request("test.queue1.req", []byte("ni hao"), time.Second)
	msg, err = nc.Request("test.queue1.req", []byte("ni hao"), time.Second)
	if err != nil {
		log.Fatal(err)
		return
	}
	count.Wait()
	assert.Equal(t, string(msg.Data), "get it")
	// 添加时间，等待数据上报完毕
	time.Sleep(10 * time.Second)
}

// 函数：QueueSubscribeSync 这个同步等待，是组内只有一个一直消费，这种其实负载均衡算是不支持了，这种相当于绑定了某个ip的负载均衡情况
//func TestNatsQueueGroup2(t *testing.T) {
//	nc, err := baseNats.New()
//	if err != nil {
//		logger.Fatal(err)
//		return
//	}
//	defer nc.Close()
//
//	count := sync.WaitGroup{}
//	count.Add(3)
//	sub, err := nc.QueueSubscribeSync("test.queue.req", "queue1")
//	go func() {
//		for msg, err := sub.NextMsg(time.Second); msg != nil && err == nil; {
//			fmt.Println("1")
//			fmt.Println(string(msg.Data))
//			count.Done()
//		}
//	}()
//
//	sub, err = nc.QueueSubscribeSync("test.queue.req", "queue1")
//	go func() {
//		for msg, err := sub.NextMsg(time.Second); msg != nil && err == nil; {
//			fmt.Println("2")
//			fmt.Println(string(msg.Data))
//			count.Done()
//		}
//	}()
//
//	sub, err = nc.QueueSubscribeSync("test.queue.req", "queue1")
//	go func() {
//		for msg, err := sub.NextMsg(time.Second); msg != nil && err == nil; {
//			fmt.Println("3")
//			fmt.Println(string(msg.Data))
//			count.Done()
//		}
//	}()
//
//	msg, err := nc.Request("test.queue.req", []byte("ni hao"), time.Second)
//	msg, err = nc.Request("test.queue.req", []byte("ni hao"), time.Second)
//	msg, err = nc.Request("test.queue.req", []byte("ni hao"), time.Second)
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//	count.Wait()
//	assert.Equal(t, string(msg.Data), "get it")
//}
