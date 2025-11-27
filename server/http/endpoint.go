package http

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/simonalong/gole-boot/debug"
	"github.com/simonalong/gole-boot/errorx"
	"github.com/simonalong/gole/bean"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/maps"
	baseTime "github.com/simonalong/gole/time"
	"github.com/simonalong/gole/util"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func addSystemRoute(httpServer *ServerOfHttp) {
	// 注册 健康检查endpoint
	if config.GetValueBoolDefault("gole.endpoint.health.enable", false) {
		RegisterHealthCheckEndpoint(httpServer)
	}

	// 注册 配置查看和变更功能
	if config.GetValueBoolDefault("gole.endpoint.config.enable", false) {
		RegisterConfigWatchEndpoint(httpServer)
	}

	// 注册 bean管理的功能
	if config.GetValueBoolDefault("gole.endpoint.bean.enable", false) {
		RegisterBeanWatchEndpoint(httpServer)
	}

	// 注册 logger管理的功能
	if config.GetValueBoolDefault("gole.endpoint.logger.enable", false) {
		RegisterLoggerEndpoint(httpServer)
	}

	// 注册：调试相关的endpoint
	if config.GetValueBoolDefault("gole.debug.enable", true) {
		// 注册 debug的帮助命令
		RegisterHelpEndpoint(httpServer)
	}

	// 注册 swagger的功能
	if config.GetValueBoolDefault("gole.swagger.enable", false) {
		RegisterSwaggerEndpoint(httpServer)
	}
}

/// ----------------------------------------------------------------------------------

func RegisterHealthCheckEndpoint(httpServer *ServerOfHttp) {
	httpServer.AddRoute(HmGet, getPathAppendApiModel(httpServer.ServiceName, "/system/status"), systemStatus)
	httpServer.AddRoute(HmGet, getPathAppendApiModel(httpServer.ServiceName, "/system/init"), systemInit)
	httpServer.AddRoute(HmGet, getPathAppendApiModel(httpServer.ServiceName, "/system/destroy"), systemDestroy)
}

func RegisterConfigWatchEndpoint(httpServer *ServerOfHttp) {
	httpServer.AddRoute(HmGet, getPathAppendApiModel(httpServer.ServiceName, "/config/values/properties"), configGetConfigValues)
	httpServer.AddRoute(HmGet, getPathAppendApiModel(httpServer.ServiceName, "/config/values/yaml"), configGetConfigDeepValues)
	httpServer.AddRoute(HmGet, getPathAppendApiModel(httpServer.ServiceName, "/config/values/json"), configGetConfigJson)
	httpServer.AddRoute(HmGet, getPathAppendApiModel(httpServer.ServiceName, "/config/value/:key"), configGetConfigValue)
	httpServer.AddRoute(HmPut, getPathAppendApiModel(httpServer.ServiceName, "/config/update"), configUpdateConfig)
}

func RegisterBeanWatchEndpoint(httpServer *ServerOfHttp) {
	httpServer.AddRoute(HmGet, getPathAppendApiModel(httpServer.ServiceName, "/bean/name/all"), debugBeanAll)
	httpServer.AddRoute(HmGet, getPathAppendApiModel(httpServer.ServiceName, "/bean/name/list/:name"), debugBeanList)
	httpServer.AddRoute(HmPost, getPathAppendApiModel(httpServer.ServiceName, "/bean/field/get"), debugBeanGetField)
	httpServer.AddRoute(HmPut, getPathAppendApiModel(httpServer.ServiceName, "/bean/field/set"), debugBeanSetField)
	httpServer.AddRoute(HmPost, getPathAppendApiModel(httpServer.ServiceName, "/bean/fun/call"), debugBeanFunCall)
}

func RegisterLoggerEndpoint(httpServer *ServerOfHttp) {
	httpServer.AddRoute(HmGet, getPathAppendApiModel(httpServer.ServiceName, "/logger/group/list/:name"), loggerGroupList)
	httpServer.AddRoute(HmPut, getPathAppendApiModel(httpServer.ServiceName, "/logger/root"), loggerRootUpdate)
	httpServer.AddRoute(HmPut, getPathAppendApiModel(httpServer.ServiceName, "/logger/group"), loggerGroupUpdate)
}

func RegisterHelpEndpoint(httpServer *ServerOfHttp) {
	httpServer.AddRoute(HmGet, getPathAppendApiModel(httpServer.ServiceName, "/debug/help"), debugHelp)
}

func RegisterSwaggerEndpoint(httpServer *ServerOfHttp) {
	httpServer.AddRouteGinHandler(HmGet, "/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

/// ----------------------------------------------------------------------------------

func loggerGroupList(c *gin.Context) (any, error) {
	loggerGroupNames := logger.GetLoggerGroupList(c.Param("name"))
	loggerGroupLevelMap := maps.New()
	for _, loggerGroupName := range loggerGroupNames {
		loggerGroup := logger.Group(loggerGroupName)
		loggerGroupLevelMap.Set(loggerGroupName, loggerGroup.Level.String())
	}
	return loggerGroupLevelMap.ToMap(), nil
}

// root级别变更会影响所有的分组级别
func loggerRootUpdate(c *gin.Context) (any, error) {
	rootLevel := c.Query("level")
	logger.SetGlobalLevel(rootLevel)

	// 修改所有的分组级别
	allGroupNames := logger.GetLoggerGroupList("")
	for _, groupName := range allGroupNames {
		logger.SetGroupLevel(groupName, rootLevel)
	}
	return "ok", nil
}

func loggerGroupUpdate(c *gin.Context) (any, error) {
	loggerGroupName := c.Query("group")
	loggerGroupLevel := c.Query("level")
	logger.SetGroupLevel(loggerGroupName, loggerGroupLevel)
	return "ok", nil
}

func debugBeanAll(c *gin.Context) (any, error) {
	return bean.GetBeanNames(""), nil
}

func debugBeanList(c *gin.Context) (any, error) {
	return bean.GetBeanNames(c.Param("name")), nil
}

func debugBeanGetField(c *gin.Context) (any, error) {
	fieldGetReq := bean.FieldGetReq{}
	_, err := util.DataToEntity(c.Request.Body, &fieldGetReq)
	if err != nil {
		return "", errorx.SC_BAD_REQUEST.WithError(err)
	}
	return bean.GetField(fieldGetReq.Bean, fieldGetReq.Field), nil
}

func debugBeanSetField(c *gin.Context) (any, error) {
	fieldSetReq := bean.FieldSetReq{}
	_, err := util.DataToEntity(c.Request.Body, &fieldSetReq)
	if err != nil {
		return "", errorx.SC_BAD_REQUEST.WithError(err)
	}
	bean.SetField(fieldSetReq.Bean, fieldSetReq.Field, fieldSetReq.Value)
	return fieldSetReq.Value, nil
}

func debugBeanFunCall(c *gin.Context) (any, error) {
	funCallReq := bean.FunCallReq{}
	_, err := util.DataToEntity(c.Request.Body, &funCallReq)
	if err != nil {
		return "", errorx.SC_BAD_REQUEST.WithError(err)
	}
	return bean.CallFun(funCallReq.Bean, funCallReq.Fun, funCallReq.Parameter), nil
}

func debugHelp(c *gin.Context) (any, error) {
	return debug.GetHelpPrintMap(), nil
}

func configGetConfigValues(c *gin.Context) (any, error) {
	return config.GetConfigValues(), nil
}

func configGetConfigDeepValues(c *gin.Context) (any, error) {
	data, err := util.ObjectToYaml(config.GetConfigDeepValues())
	if err != nil {
		return nil, errorx.SC_BAD_REQUEST.WithError(err)
	}
	return data, nil
}

func configGetConfigJson(c *gin.Context) (any, error) {
	dataMap, err := util.JsonToMap(util.ObjectToJson(config.GetConfigDeepValues()))
	if err != nil {
		return nil, errorx.SC_BAD_REQUEST.WithError(err)
	}
	return dataMap, nil
}

func configGetConfigValue(c *gin.Context) (any, error) {
	return config.GetConfigValue(c.Param("key")), nil
}

func configUpdateConfig(c *gin.Context) (any, error) {
	valueMap := map[string]any{}
	_, err := util.DataToEntity(c.Request.Body, &valueMap)
	if err != nil {
		logger.Errorf("解析失败，%v", err.Error())
		return nil, errorx.SC_BAD_REQUEST.WithError(err)
	}

	key, _ := valueMap["key"]
	value, _ := valueMap["value"]
	config.UpdateConfig(util.ToString(key), value)
	return "ok", nil
}

var procId = os.Getpid()
var startTime = baseTime.TimeToStringYmdHms(time.Now())

const defaultVersion = "unknown"

func systemStatus(c *gin.Context) (any, error) {
	return maps.Of("status", "ok", "running", true, "pid", procId, "startupAt", startTime, "version", getVersion()).ToMap(), nil
}

func systemInit(c *gin.Context) (any, error) {
	return maps.Of("status", "ok"), nil
}

func systemDestroy(c *gin.Context) (any, error) {
	return maps.Of("status", "ok"), nil
}

func getVersion() string {
	return config.GetValueStringDefault("gole.application.version", defaultVersion)
}
