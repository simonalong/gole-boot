package original

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/simonalong/gole/logger"
	"testing"
)

func TestGnetServer(t *testing.T) {
	err := gnet.Run(&DemoXxxHandler{}, "tcp://localhost:9992")
	if err != nil {
		logger.Errorf("异常%v", err)
	}
}
