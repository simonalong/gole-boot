package sys

import (
	"fmt"
	"github.com/simonalong/gole/config"
	"sync"
)

var BootVersion = "v1.6.2"

var loadLock sync.Mutex
var Loaded = false

func init() {
	config.Load()
	PrintBanner()
}

func PrintBanner() {
	loadLock.Lock()
	defer loadLock.Unlock()
	if Loaded {
		return
	}
	Loaded = true
	if config.GetValueBoolDefault("gole.application.banner", true) {
		fmt.Printf("%s", Banner)
	}
	fmt.Printf("---------------------------------------- cbb-gole-boot: %s ----------------------------------------\n", BootVersion)
	fmt.Printf("profile：%s\n", config.CurrentProfile)
	appName := config.BaseCfg.Application.Name
	if appName != "" {
		fmt.Printf("服务名：%v\n", config.BaseCfg.Application.Name)
	}
	appVersion := config.BaseCfg.Application.Version
	if appVersion != "" {
		fmt.Printf("版本号：%v\n", appVersion)
	}
	fmt.Printf("-------------------------------------------------------------------------------------------------------\n")
}
