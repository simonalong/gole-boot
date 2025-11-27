package http

import (
	"github.com/gin-gonic/gin"
	"github.com/simonalong/gole-boot/server/http/rsp"
	"net/http"
)

type RouteHandler func(c *gin.Context) (any, error)

type ServerOfHttp struct {
	ServiceName string
	IRouter     gin.IRouter
}

func (server *ServerOfHttp) GetEngine() *gin.Engine {
	return server.IRouter.(*gin.Engine)
}

func (server *ServerOfHttp) GetRouterGroup() *gin.RouterGroup {
	return server.IRouter.(*gin.RouterGroup)
}

// AddMiddleware 添加中间件
func (server *ServerOfHttp) AddMiddleware(handler ...gin.HandlerFunc) *ServerOfHttp {
	if server.IRouter == nil {
		return nil
	}
	server.IRouter.Use(handler...)
	return server
}

func (server *ServerOfHttp) Post(path string, serverHandler RouteHandler) *ServerOfHttp {
	return server.AddRoute(HmPost, getPathAppendApiModel(server.ServiceName, path), serverHandler)
}

func (server *ServerOfHttp) Delete(path string, serverHandler RouteHandler) *ServerOfHttp {
	return server.AddRoute(HmDelete, getPathAppendApiModel(server.ServiceName, path), serverHandler)
}

func (server *ServerOfHttp) Put(path string, serverHandler RouteHandler) *ServerOfHttp {
	return server.AddRoute(HmPut, getPathAppendApiModel(server.ServiceName, path), serverHandler)
}

func (server *ServerOfHttp) Patch(path string, serverHandler RouteHandler) *ServerOfHttp {
	return server.AddRoute(HmPatch, getPathAppendApiModel(server.ServiceName, path), serverHandler)
}

func (server *ServerOfHttp) Head(path string, serverHandler RouteHandler) *ServerOfHttp {
	return server.AddRoute(HmHead, getPathAppendApiModel(server.ServiceName, path), serverHandler)
}

func (server *ServerOfHttp) Get(path string, serverHandler RouteHandler) *ServerOfHttp {
	return server.AddRoute(HmGet, getPathAppendApiModel(server.ServiceName, path), serverHandler)
}

func (server *ServerOfHttp) Options(path string, serverHandler RouteHandler) *ServerOfHttp {
	return server.AddRoute(HmOptions, getPathAppendApiModel(server.ServiceName, path), serverHandler)
}

func (server *ServerOfHttp) All(path string, serverHandler RouteHandler) *ServerOfHttp {
	return server.AddRoute(HmAll, getPathAppendApiModel(server.ServiceName, path), serverHandler)
}

func (server *ServerOfHttp) AddRoute(method Method, path string, routeHandler RouteHandler) *ServerOfHttp {
	return server.AddRouteGinHandler(method, path, ginHandler(routeHandler))
}

func ginHandler(serverHandler RouteHandler) func(context *gin.Context) {
	return func(context *gin.Context) {
		result, err := serverHandler(context)
		rsp.Done(context, result, err)
	}
}

func (server *ServerOfHttp) AddRouteGinHandler(method Method, path string, ginHandler gin.HandlerFunc) *ServerOfHttp {
	if server == nil {
		return nil
	}
	iRouter := server.IRouter
	handler := ginHandler
	switch method {
	case HmAll:
		iRouter.GET(path, handler)
		iRouter.POST(path, handler)
		iRouter.PUT(path, handler)
		iRouter.DELETE(path, handler)
		iRouter.OPTIONS(path, handler)
		iRouter.HEAD(path, handler)
	case HmGet:
		iRouter.GET(path, handler)
	case HmPost:
		iRouter.POST(path, handler)
	case HmPut:
		iRouter.PUT(path, handler)
	case HmPatch:
		iRouter.PATCH(path, handler)
	case HmDelete:
		iRouter.DELETE(path, handler)
	case HmOptions:
		iRouter.OPTIONS(path, handler)
	case HmHead:
		iRouter.HEAD(path, handler)
	}
	return server
}

func (server *ServerOfHttp) RegisterHttpHandler(method Method, path string, handler http.Handler) *ServerOfHttp {
	return server.AddRouteGinHandler(method, path, gin.WrapH(handler))
}

func (server *ServerOfHttp) Group(relativePath string, handlers ...gin.HandlerFunc) *ServerOfHttp {
	var newRouterGroup *gin.RouterGroup
	if engine := server.IRouter.(*gin.Engine); engine != nil {
		newRouterGroup = engine.Group(relativePath, handlers...)
	} else if routerGroup := server.IRouter.(*gin.RouterGroup); routerGroup != nil {
		newRouterGroup = routerGroup.Group(relativePath, handlers...)
	}
	return &ServerOfHttp{
		ServiceName: server.ServiceName,
		IRouter:     newRouterGroup,
	}
}
