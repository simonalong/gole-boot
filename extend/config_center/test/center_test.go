package test

import (
	"fmt"
	_ "github.com/simonalong/gole-boot/extend/config_center"
	"github.com/simonalong/gole/config"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		fmt.Println(config.GetValueString("biz.key"))
		time.Sleep(time.Second)
	}
}
