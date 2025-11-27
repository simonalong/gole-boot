package test

import (
	"fmt"
	"github.com/simonalong/gole-boot/extend/tdengine"
	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/maps"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	tdClient, err := tdengine.NewClient()
	if err != nil {
		panic(err)
		return
	}
	baseMap := maps.Of("ts", time.Now(), "name", "大牛市-boot", "age", 28, "address", "浙江杭州市")
	num, err := tdClient.Insert("td_china", baseMap)
	assert.Equal(t, 1, num)
}

// 使用多库配置：gole.profiles.active=multi
// 需要提前创建库
func TestInsertOfMulti(t *testing.T) {
	tdClient, err := tdengine.NewClientWithName("demo1")
	if err != nil {
		panic(err)
		return
	}

	//创建超级表
	_, err = tdClient.Exec("create stable if not exists td_orm1.td_demo1(ts timestamp, name nchar(32), age int, address nchar(128)) tags (station nchar(128))")

	// 添加数据顺便创建子表（不存在的话）
	baseMap := maps.Of("ts", time.Now(), "name", "大牛市-boot", "age", 28, "address", "浙江杭州市")
	tagMap := maps.Of("station", "china")
	num, err := tdClient.InsertWithTag("td_china", "td_demo1", tagMap, baseMap)
	assert.Equal(t, 1, num)

	tdOrm2, err := tdengine.NewClientWithName("demo2")
	if err != nil {
		panic(err)
		return
	}

	baseMap = maps.Of("ts", time.Now(), "name", "大牛市-boot", "age", 28, "address", "浙江杭州市")
	num, err = tdOrm2.Insert("td_china", baseMap)
	assert.Equal(t, 1, num)
}

const StableNameOfTaskLog = "task_log"

func TestDemo(t *testing.T) {
	tdClient, err := tdengine.NewClient()
	if err != nil {
		panic(err)
		return
	}

	productNo := "demo1ForTest"
	devNo := "demo1ForTest"
	tableName := getTableNameOfTaskLog(productNo, devNo)
	taskLogDo := &TaskLogDo{
		Ts:                time.Now(),
		TraceId:           "xxxxxxxxxx",
		PropertyOrService: 1,
		DriverCode:        "drv-stk-aep-water",
		Module:            "default",
		RunStatus:         2,
		Enable:            1,
		RetryNum:          1,
		MaxRetryNum:       3,
		ErrMsg:            "",
		ServiceId:         "",
		Content:           "{\"module\":\"default\",\"propertyName\":\"powerType\",\"propertyValue\":3}",
	}
	//taskLogDo := maps.Of("ts", time.Now().UnixMilli(), "trace_id", "xxxxxxxxxx")
	tagsMap := maps.OfSort("product_no", productNo, "dev_no", devNo)

	//_, err = tdClient.InsertWithTag(tableName, StableNameOfTaskLog, tagsMap, taskLogDo)
	_, err = tdClient.InsertEntityWithTag(tableName, StableNameOfTaskLog, tagsMap, taskLogDo)
	if err != nil {
		logger.Errorf("InsertTaskLog err:%s", err.Error())
	}
}

func getTableNameOfTaskLog(productNo, devNo string) string {
	return fmt.Sprintf("task_log_%s_%s", productNo, devNo)
}

type TaskLogDo struct {
	Ts                time.Time `json:"ts"`
	TraceId           string    `json:"trace_id"`            // 跟踪id
	PropertyOrService int       `json:"property_or_service"` // 属性设置还是函数调用：1-属性设置；2-函数调用
	DriverCode        string    `json:"driver_code"`         // 驱动code
	Module            string    `json:"module"`              // 模块
	ExecuteTime       time.Time `json:"execute_time"`        // 执行时间
	RunStatus         int       `json:"run_status"`          // 任务执行结果：0-未处理；1-处理中；2-成功；3-异常
	Enable            int       `json:"enable"`              // 是否启用：0-禁用，1-启用
	RetryNum          int       `json:"retry_num"`           // 任务重试次数
	MaxRetryNum       int       `json:"max_retry_num"`       // 任务最大重试次数
	ErrMsg            string    `json:"err_msg"`             // 异常信息
	ServiceId         string    `json:"service_id"`          // 函数名（只有为函数调用时候才需要）
	Content           string    `json:"content"`             // 属性设置：则为属性的参数 kv；函数调用：则为函数的参数 json 字符串
	CreateTime        time.Time `json:"create_time"`         // 创建时间
	UpdateTime        time.Time `json:"update_time"`         // 更新时间
}
