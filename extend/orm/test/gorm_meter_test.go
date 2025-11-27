package test

import (
	"github.com/gin-gonic/gin"
	orm2 "github.com/simonalong/gole-boot/extend/orm"
	httpServer "github.com/simonalong/gole-boot/server/http"
	"github.com/simonalong/gole/logger"
	baseTime "github.com/simonalong/gole/time"
	"testing"
	"time"
)

// gole.profiles.active=demo-http
func TestGormMeter(t *testing.T) {
	db, err := orm2.GetGormClient()
	if err != nil {
		logger.Errorf("初始化数据库失败：%v", err)
		return
	}

	httpServer.Get("/orm/init", func(context *gin.Context) (any, error) {
		// 删除表
		//db.Exec("drop table base_demo")

		//新增表
		db.Exec("CREATE TABLE base_demo(\n" +
			"  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',\n" +
			"  `name` char(20) NOT NULL COMMENT '名字',\n" +
			"  `age` INT NOT NULL COMMENT '年龄',\n" +
			"  `address` char(20) NOT NULL COMMENT '名字',\n" +
			"  \n" +
			"  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',\n" +
			"  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',\n" +
			"\n" +
			"  PRIMARY KEY (`id`)\n" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='测试表'")

		// 新增
		//db.Create(&BaseDemo{Name: "zhou", Age: 18, Address: "杭州"})
		//db.Create(&BaseDemo{Name: "zhou", Age: 11, Address: "杭州2"})

		// 查询：一行
		//var demo BaseDemo
		//db.First(&demo).Where("name=?", "zhou")

		//dd, _ := db.DB()
		//dd.Query("select * from base_demo")

		// 查询：多行
		//fmt.Println(demo)
		return "ok", nil
	})

	httpServer.Get("/orm/insert", func(context *gin.Context) (any, error) {
		db.Create(&BaseDemo{Name: baseTime.TimeToStringYmd(time.Now()), Age: 11, Address: "杭州2"})
		return "ok", nil
	})

	httpServer.Get("/orm/query", func(context *gin.Context) (any, error) {
		var demo BaseDemo
		db.First(&demo).Where("name=?", "zhou")
		return demo, nil
		//dd, _ := db.DB()
		//rows, _ := dd.Query("select * from base_demo")
		//rowMaps := maps.FromSqlRows(rows)
		//var dataMaps []map[string]interface{}
		//for _, rowMap := range rowMaps {
		//	dataMaps = append(dataMaps, rowMap.ToMap())
		//}
		//rsp.Done(context, dataMaps)
	})
	httpServer.RunServer()
}
