package test

import (
	orm2 "github.com/simonalong/gole-boot/extend/orm"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/logger"
	"testing"
)

func TestConnect(t *testing.T) {
	config.LoadYamlFile("./application-connect.yaml")
	_, err := orm2.GetGormClient()
	if err != nil {
		logger.Fatalf("连接数据库异常：%v", err)
	}
	logger.Info("连接数据库成功")
}
