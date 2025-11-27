package event

import "github.com/simonalong/gole/listener"

var EventOfServerRunFinish = "event_of_server_run_finish"

// ServerRunFinishEvent 服务（http、tcp、udp）运行完成事件
type ServerRunFinishEvent struct{}

func (e ServerRunFinishEvent) Name() string {
	return EventOfServerRunFinish
}

func (e ServerRunFinishEvent) Group() string {
	return listener.DefaultGroup
}

func (e ServerRunFinishEvent) ToString() string {
	return ""
}
