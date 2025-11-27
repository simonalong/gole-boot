package nats

import (
	baseNats "github.com/simonalong/gole-boot/extend/nats"
	"github.com/simonalong/gole/logger"
	"testing"
	"time"
)

// 使用本例子之前请使用 ./deploy/xxx 文件夹中的部署进行启动即可
func TestNatsConnectWithUser(t *testing.T) {
	_, _, err := baseNats.GetJetStreamClient()
	if err != nil {
		logger.Fatal(err)
		return
	}
	time.Sleep(12 * time.Hour)
}
