package debug

import (
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/util"
)

func GetHelpPrintMap() map[string]interface{} {
	port := config.GetValueIntDefault("gole.server.http.port", 8080)
	cmdMap := map[string]interface{}{}
	cmdMap["帮助"] = "curl http://localhost:" + pre(port) + "/debug/help"
	cmdMap["日志"] = logHelp(port)
	cmdMap["http接口出入参"] = httpHelp(port)
	cmdMap["bean管理"] = beanHelp(port)
	cmdMap["pprof"] = pprofHelp(port)
	cmdMap["配置处理"] = configHelp(port)
	return cmdMap
}

func logHelp(port int) map[string]interface{} {
	return map[string]interface{}{
		"日志分组列表": "curl http://localhost:" + pre(port) + "/logger/list/{name}",
		"动态修改日志": "curl -X PUT http://localhost:" + pre(port) + "/config/update -d '{\"key\":\"gole.logger.level\", \"value\":\"debug\"}'",
	}
}
func beanHelp(port int) map[string]interface{} {
	return map[string]interface{}{
		"获取注册的所有bean": "curl http://localhost:" + pre(port) + "/bean/name/all",
		"查询注册的bean":   "curl http://localhost:" + pre(port) + "/bean/name/list/{name}",
		"查询bean的属性值":  "curl -X POST http://localhost:" + pre(port) + "/bean/field/get' -d '{\"bean\": \"xx\", \"field\":\"xxx\"}'",
		"修改bean的属性值":  "curl -X POST http://localhost:" + pre(port) + "/bean/field/set' -d '{\"bean\": \"xx\", \"field\": \"xxx\", \"value\": \"xxx\"}'",
		"调用bean的函数":   "curl -X POST http://localhost:" + pre(port) + "/bean/fun/call' -d '{\"bean\": \"xx\", \"fun\": \"xxx\", \"parameter\": {\"p1\":\"xx\", \"p2\": \"xxx\"}}'",
	}
}
func httpHelp(port int) map[string]interface{} {
	return map[string]interface{}{
		"指定url打印请求":     "curl -X PUT http://localhost:" + pre(port) + "/config/update -d '{\"key\":\"gole.server.http.request.print.include-uri[0]\", \"value\":\"/api/xx/xxx\"}'",
		"指定url不打印请求":    "curl -X PUT http://localhost:" + pre(port) + "/config/update -d '{\"key\":\"gole.server.http.request.print.exclude-uri[0]\", \"value\":\"/api/xx/xxx\"}'",
		"指定url打印请求和响应":  "curl -X PUT http://localhost:" + pre(port) + "/config/update -d '{\"key\":\"gole.server.http.response.print.include-uri[0]\", \"value\":\"/api/xx/xxx\"}'",
		"指定url不打印请求和响应": "curl -X PUT http://localhost:" + pre(port) + "/config/update -d '{\"key\":\"gole.server.http.response.print.exclude-uri[0]\", \"value\":\"/api/xx/xxx\"}'",
	}
}
func pprofHelp(port int) map[string]interface{} {
	return map[string]interface{}{
		"动态启用pprof": "curl -X PUT http://localhost:" + pre(port) + "/config/update -d '{\"key\":\"gole.server.http.pprof.enable\", \"value\":\"true\"}'",
	}
}
func configHelp(port int) map[string]interface{} {
	return map[string]interface{}{
		"服务所有配置":         "curl http://localhost:" + pre(port) + "/config/values",
		"服务所有配置(yaml结构)": "curl http://localhost:" + pre(port) + "/config/values/yaml",
		"服务某个配置":         "curl http://localhost:" + pre(port) + "/config/value/{key}",
		"修改服务的配置":        "curl -X PUT http://localhost:" + pre(port) + "/config/update -d '{\"key\":\"xxx\", \"value\":\"yyy\"}'",
	}
}

func pre(port int) string {
	return util.ToString(port) + "/" + apiPreAndModule()
}

func apiPreAndModule() string {
	apiPrefix := util.ISCString(config.GetValueStringDefault("gole.server.http.api.prefix", "")).Trim("/")
	if apiPrefix == "" {
		apiPrefix = "api"
	}
	return apiPrefix.ToString()
}
