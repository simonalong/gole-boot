package http

var EventOfServerHttpInit = "event_of_server_http_init"

// ServerHttpInitEvent
// @Description: 服务器启动后，添加中间件事件；一般用于给业务这边添加中间件模块使用
type ServerHttpInitEvent struct {
	ServiceName  string
	ServerOfHttp *ServerOfHttp
}

func (e ServerHttpInitEvent) Name() string {
	return EventOfServerHttpInit
}

func (e ServerHttpInitEvent) Group() string {
	return e.ServiceName
}

func (e ServerHttpInitEvent) ToString() string {
	return ""
}
