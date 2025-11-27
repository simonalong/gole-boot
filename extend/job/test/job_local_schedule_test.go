package test

import (
	"github.com/simonalong/gole-boot/extend/job"
	"github.com/simonalong/gole/logger"
	"testing"
	"time"
)

// gole.profiles.active=local-schedule
func TestScheduleCron1(t *testing.T) {
	job.ScheduleCron("cbb-job-demo:TestScheduleCron", "* * * * * ?", job1)

	time.Sleep(12 * time.Hour)
}

// gole.profiles.active=local-schedule
func TestScheduleCron2(t *testing.T) {
	job.ScheduleCron("cbb-job-demo:TestScheduleCron", "* * * * * ?", job1)

	time.Sleep(12 * time.Hour)
}

// gole.profiles.active=local-schedule
func TestScheduleFixRate1(t *testing.T) {
	job.ScheduleFixRate("cbb-job-demo:TestScheduleFixRate", time.Second, job2)

	time.Sleep(12 * time.Hour)
}

// gole.profiles.active=local-schedule
func TestScheduleFixRate2(t *testing.T) {
	job.ScheduleFixRate("cbb-job-demo:TestScheduleFixRate", time.Second, job2)

	time.Sleep(12 * time.Hour)
}

func job1() {
	logger.Info("job1")
	time.Sleep(500 * time.Millisecond)
}

func job2() {
	logger.Info("job2")
	time.Sleep(500 * time.Millisecond)
}

func job3() {
	logger.Info("job3")
}
