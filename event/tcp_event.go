package event

import "github.com/simonalong/gole/listener"

var EventOfServerTcpRunFinish = "event_of_server_tcp_run_finish"
var EventOfServerTcpRunStart = "event_of_server_run_tcp_start"
var EventOfServerTcpStop = "event_of_server_tcp_stop"

type ServerTcpRunStartEvent struct{}

// ServerTcpRunFinishEvent tcp服务完成启动事件, 对应：event_of_server_tcp_run_finish
type ServerTcpRunFinishEvent struct{}

type ServerTcpStopEvent struct{}

func (e ServerTcpRunStartEvent) Name() string {
	return EventOfServerTcpRunStart
}

func (e ServerTcpRunStartEvent) Group() string {
	return listener.DefaultGroup
}
func (e ServerTcpRunStartEvent) ToString() string {
	return EventOfServerTcpRunStart
}

func (e ServerTcpStopEvent) Name() string {
	return EventOfServerTcpStop
}

func (e ServerTcpStopEvent) Group() string {
	return listener.DefaultGroup
}
func (e ServerTcpStopEvent) ToString() string {
	return EventOfServerTcpStop
}

func (e ServerTcpRunFinishEvent) Name() string {
	return EventOfServerTcpRunFinish
}

func (e ServerTcpRunFinishEvent) Group() string {
	return listener.DefaultGroup
}
func (e ServerTcpRunFinishEvent) ToString() string {
	return EventOfServerTcpRunFinish
}
