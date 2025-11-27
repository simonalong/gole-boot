package event

import "github.com/simonalong/gole/listener"

var EventOfServerGrpcRunFinish = "event_of_server_grpc_run_finish"
var EventOfServerGrpcRunStart = "event_of_server_run_grpc_start"
var EventOfServerGrpcStop = "event_of_server_grpc_stop"
var EventOfServerAllGrpcRunFinish = "event_of_server_all_grpc_run_finish"

type ServerGrpcRunStartEvent struct {
	ServiceName string
}

// ServerGrpcRunFinishEvent tcp服务完成启动事件, 对应：event_of_server_grpc_run_finish
type ServerGrpcRunFinishEvent struct {
	ServiceName string
}

// ServerGrpcStopEvent grpc服务关闭事件, 对应：event_of_server_grpc_stop
type ServerGrpcStopEvent struct {
	ServiceName string
}

// ServerGrpcAllRunFinishEvent 所有grpc服务完成启动事件, 对应：event_of_server_all_grpc_run_finish（用于多http服务的情况）
type ServerGrpcAllRunFinishEvent struct {
}

func (e ServerGrpcRunStartEvent) Name() string {
	return EventOfServerGrpcRunStart
}

func (e ServerGrpcRunStartEvent) Group() string {
	return e.ServiceName
}

func (e ServerGrpcRunStartEvent) ToString() string {
	return ""
}

func (e ServerGrpcRunFinishEvent) Name() string {
	return EventOfServerGrpcRunFinish
}

func (e ServerGrpcRunFinishEvent) Group() string {
	return e.ServiceName
}

func (e ServerGrpcRunFinishEvent) ToString() string {
	return ""
}

func (e ServerGrpcStopEvent) Name() string {
	return EventOfServerGrpcStop
}

func (e ServerGrpcStopEvent) Group() string {
	return e.ServiceName
}

func (e ServerGrpcStopEvent) ToString() string {
	return ""
}

func (e ServerGrpcAllRunFinishEvent) Name() string {
	return EventOfServerAllGrpcRunFinish
}

func (e ServerGrpcAllRunFinishEvent) Group() string {
	return listener.DefaultGroup
}

func (e ServerGrpcAllRunFinishEvent) ToString() string {
	return ""
}
