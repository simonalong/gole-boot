/*
* 使用本例子之前请使用 ./deploy/xxx 文件夹中的部署进行启动即可
* 本文件全部使用：环境变量：gole.profiles.active=user
 */
package nats

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"github.com/nats-io/nats.go"
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/logger"
)

func TestNatsReqRsp1(t *testing.T) {
	nc, _, err := baseNats.GetJetStreamClient()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer nc.Close()

	count := sync.WaitGroup{}
	count.Add(1)
	_, err = nc.Subscribe("test.req_rsp1.req", func(msg *baseNats.MsgOfNats) {
		assert.Equal(t, string(msg.Data), "ni hao")
		time.Sleep(2 * time.Second)
		err := msg.Respond([]byte("hello world"))
		if err != nil {
			log.Fatal(err)
			return
		}
		count.Done()
	})

	msg, err := nc.Request("test.req_rsp1.req", []byte("ni hao"), time.Second)
	if err != nil {
		logger.Errorf("Request failed: %v", err)
		//log.Fatal(err)
		return
	}
	count.Wait()
	assert.Equal(t, string(msg.Data), "hello world")
}

func TestNatsReqRsp2(t *testing.T) {
	nc, err := baseNats.New()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer nc.Close()

	count := sync.WaitGroup{}
	count.Add(1)
	_, err = nc.Subscribe("test.req_rsp2.req", func(msg *baseNats.MsgOfNats) {
		assert.Equal(t, string(msg.Data), "ni hao")
		err := msg.Respond([]byte("hello world"))
		if err != nil {
			log.Fatal(err)
			return
		}
		count.Done()
	})

	pMsg := &nats.Msg{
		Subject: "test.req_rsp2.req",
		Data:    []byte("ni hao"),
	}
	msg, err := nc.RequestMsg(pMsg, time.Second)
	if err != nil {
		log.Fatal(err)
		return
	}
	count.Wait()
	assert.Equal(t, string(msg.Data), "hello world")
}

func TestNatsReqRsp3(t *testing.T) {
	nc, err := baseNats.New()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer nc.Close()

	count := sync.WaitGroup{}
	count.Add(2)
	_, err = nc.Subscribe("test.req_rsp3.req", func(msg *baseNats.MsgOfNats) {
		assert.Equal(t, string(msg.Data), "ni hao")
		err := msg.Respond([]byte("hello world"))
		if err != nil {
			log.Fatal(err)
			return
		}
		count.Done()
	})

	_, err = nc.Subscribe("test.req_rsp3.rsp", func(msg *baseNats.MsgOfNats) {
		assert.Equal(t, string(msg.Data), "hello world")
		count.Done()
	})

	// 原来这个时候，这个reply才排上用场了
	err = nc.PublishRequest("test.req_rsp3.req", "test.req_rsp3.rsp", []byte("ni hao"))
	if err != nil {
		log.Fatal(err)
		return
	}
	count.Wait()
}
