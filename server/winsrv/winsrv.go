package winsrv

import (
	kardSrv "github.com/kardianos/service"
	"github.com/simonalong/gole/logger"
	"os"
)

type Runnable func()

type WinService struct {
	DisplayName string
	Runnable    Runnable
}

func (winSvc *WinService) Start(s kardSrv.Service) error {
	logger.Infof("【%v】Start.......", winSvc.DisplayName)
	go winSvc.Runnable()
	return nil
}

func (winSvc *WinService) Stop(s kardSrv.Service) error {
	logger.Infof("【%v】Stop.......", winSvc.DisplayName)
	return nil
}

func Run(svcConfig *kardSrv.Config, runnable Runnable) {
	winS := &WinService{
		DisplayName: svcConfig.DisplayName,
		Runnable:    runnable,
	}
	winServer, err := kardSrv.New(winS, svcConfig)
	if err != nil {
		logger.Fatalf("【%v】创建失败：%v", svcConfig.DisplayName, err)
		return
	}

	if len(os.Args) > 1 {
		err := kardSrv.Control(winServer, os.Args[1])
		if err != nil {
			logger.Errorf("%v: %v", getErrStr(os.Args[1]), err)
		} else {
			logger.Infof("%v", getSuccessStr(os.Args[1]))
		}
		return
	}
	err = winServer.Run()
	if err != nil {
		logger.Fatalf("【%v】运行失败: %v", svcConfig.DisplayName, err)
		return
	}
}
func getErrStr(cmd string) string {
	switch cmd {
	case "install":
		return "安装服务失败"
	case "uninstall":
		return "卸载服务失败"
	case "start":
		return "服务启动失败"
	case "stop":
		return "服务启动失败"
	}
	return ""
}
func getSuccessStr(cmd string) string {
	switch cmd {
	case "install":
		return "安装服务成功"
	case "uninstall":
		return "卸载服务成功"
	case "start":
		return "服务启动成功"
	case "stop":
		return "服务停止成功"
	}
	return ""
}
