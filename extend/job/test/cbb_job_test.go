package test

import (
	"context"
	"github.com/simonalong/gole-boot/extend/job"
	"github.com/simonalong/gole-boot/server"
	"testing"
)

// gole.profiles.active=n1
func TestCbbJob1(t *testing.T) {
	// 添加任务处理器
	job.AddJobHandler("test.job1", testJob1)
	job.AddJobHandler("test.job2", testJob2)
	job.AddJobHandler("test.job3", testJob3)

	server.Run()
}

// gole.profiles.active=n2
func TestCbbJob2(t *testing.T) {
	// 添加任务处理器
	job.AddJobHandler("test.job1", testJob1)
	job.AddJobHandler("test.job2", testJob2)
	job.AddJobHandler("test.job3", testJob3)

	server.Run()
}

// gole.profiles.active=n3
func TestCbbJob3(t *testing.T) {
	// 添加任务处理器
	job.AddJobHandler("test.job1", testJob1)
	job.AddJobHandler("test.job2", testJob2)
	job.AddJobHandler("test.job3", testJob3)

	server.Run()
}

func testJob1(cxt context.Context, param *job.RunReq) string {
	//logger.Infof("执行任务1, %v", util.ToJsonString(param))
	//time.Sleep(5 * time.Second)
	return "success"
}
func testJob2(cxt context.Context, param *job.RunReq) string {
	//logger.Infof("执行任务2, %v", util.ToJsonString(param))
	return "success"
}
func testJob3(cxt context.Context, param *job.RunReq) string {
	//logger.Infof("执行任务3, %v", util.ToJsonString(param))
	return "success"
}
