package event

import "github.com/simonalong/gole/listener"

var EventOfServerHttpRunStart = "event_of_server_http_run_start"
var EventOfServerHttpRunFinish = "event_of_server_http_run_finish"
var EventOfServerHttpStop = "event_of_server_http_stop"

var EventOfServerAllHttpRunFinish = "event_of_server_all_http_run_finish"

// ServerHttpRunStartEvent http服务开始启动事件, 对应：event_of_server_http_run_start
type ServerHttpRunStartEvent struct {
	ServiceName string
}

// ServerHttpRunFinishEvent http服务完成启动事件, 对应：event_of_server_http_run_finish
//type ServerHttpRunFinishEvent struct {
//	ServiceName string
//}

// ServerHttpStopEvent http服务关闭事件, 对应：event_of_server_http_stop
type ServerHttpStopEvent struct {
	ServiceName string
}

// ServerHttpAllRunFinishEvent 所有http服务完成启动事件, 对应：event_of_server_all_http_run_finish（用于多http服务的情况）
type ServerHttpAllRunFinishEvent struct {
}

func (e ServerHttpRunStartEvent) Name() string {
	return EventOfServerHttpRunStart
}

func (e ServerHttpRunStartEvent) Group() string {
	return e.ServiceName
}

func (e ServerHttpRunStartEvent) ToString() string {
	return ""
}

//func (e ServerHttpRunFinishEvent) Name() string {
//	return EventOfServerHttpRunFinish
//}
//
//func (e ServerHttpRunFinishEvent) Group() string {
//	return e.ServiceName
//}
//
//func (e ServerHttpRunFinishEvent) ToString() string {
//	return ""
//}

func (e ServerHttpStopEvent) Name() string {
	return EventOfServerHttpStop
}

func (e ServerHttpStopEvent) Group() string {
	return e.ServiceName
}

func (e ServerHttpStopEvent) ToString() string {
	return ""
}

func (e ServerHttpAllRunFinishEvent) Name() string {
	return EventOfServerAllHttpRunFinish
}

func (e ServerHttpAllRunFinishEvent) Group() string {
	return listener.DefaultGroup
}

func (e ServerHttpAllRunFinishEvent) ToString() string {
	return ""
}
