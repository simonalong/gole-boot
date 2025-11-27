/*
* 使用本例子之前请使用 ./deploy/xxx 文件夹中的部署进行启动即可
* 本文件全部使用：环境变量：gole.profiles.active=user
 */
package nats

import (
	"strconv"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/logger"
)

// gole.profiles.active=user
func TestNatsPush(t *testing.T) {
	nc, err := baseNats.New()
	if err != nil {
		logger.Fatal(err)
		return
	}
	defer nc.Close()

	for i := 0; i < 10000; i++ {
		pMsg := &nats.Msg{
			Subject: "test.sub.req",
			Data:    []byte("nats hello world " + strconv.Itoa(i)),
		}
		err = nc.PublishMsg(pMsg)

		time.Sleep(1 * time.Second)
	}
}
