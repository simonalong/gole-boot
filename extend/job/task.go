package job

import (
	"context"
	"fmt"
	"github.com/simonalong/gole/logger"
	"runtime/debug"
)

// TaskFunc 任务执行函数，返回值：是成功情况下返回的字段
type TaskFunc func(cxt context.Context, param *RunReq) string

// Task 任务
type Task struct {
	Id        int64
	Name      string
	Ext       context.Context
	Param     *RunReq
	fn        TaskFunc
	Cancel    context.CancelFunc
	StartTime int64
	EndTime   int64
}

// Run 运行任务
func (t *Task) Run(callback func(code int, msg string)) {
	defer func(cancel func()) {
		if err := recover(); err != nil {
			logger.Errorf("任务运行异常：%v", err)
			debug.PrintStack() //堆栈跟踪
			callback(FailureCode, fmt.Sprintf("任务运行异常:%v", err))
			cancel()
		}
	}(t.Cancel)

	resultChan := make(chan *string)
	go t.call(resultChan)
	select {
	case <-t.Ext.Done():
		logger.Errorf("任务【%v】超时：执行时间超过设定的超时时间 %v秒", t.Id, t.Param.ExecutorTimeout)
		callback(FailureCode, "task timeout")
		return
	case data := <-resultChan:
		callback(SuccessCode, *data)
		return
	}
}

func (t *Task) call(resultChan chan *string) {
	msg := t.fn(t.Ext, t.Param)
	resultChan <- &msg
	logger.Debugf("任务【%v】执行完毕", t.Id)
	return
}

// Info 任务信息
func (t *Task) Info() string {
	return fmt.Sprintf("任务ID[%d]任务名称[%s]参数:%s", t.Id, t.Name, t.Param.ExecutorParams)
}
