package rsp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
	"unsafe"

	"github.com/simonalong/gole/config"
	baseTime "github.com/simonalong/gole/time"

	"github.com/gin-gonic/gin"
	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/util"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ResponseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqPrint := config.GetValueBoolDefault("gole.server.http.request.print.enable", true)
		rspPrint := config.GetValueBoolDefault("gole.server.http.response.print.enable", false)

		if !reqPrint && !rspPrint {
			// 处理请求
			c.Next()
			return
		}

		// 开始时间
		startTime := time.Now()
		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Errorf("read request body failed,err = %s.", err)
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		var body any
		bodyStr := string(data)
		if "" != bodyStr && unsafe.Sizeof(bodyStr) < 10240 {
			if strings.HasPrefix(bodyStr, "{") && strings.HasSuffix(bodyStr, "}") {
				bodys := map[string]any{}
				_, _ = util.StrToObject(bodyStr, &bodys)
				body = bodys
			} else if strings.HasPrefix(bodyStr, "[") && strings.HasSuffix(bodyStr, "]") {
				var bodys []any
				_, _ = util.StrToObject(bodyStr, &bodys)
				body = bodys
			}
		}
		request := Request{
			Method:     c.Request.Method,
			Uri:        c.Request.RequestURI,
			Ip:         c.ClientIP(),
			Parameters: c.Params,
			Headers:    c.Request.Header,
			Body:       body,
		}
		if reqPrint && !rspPrint {
			logger.Group("http.server.req").Debugf("请求：%v", util.ToJsonString(request))
		}

		// 处理请求
		c.Next()

		responseMessage := Response{
			Request:    request,
			StatusCode: c.Writer.Status(),
			CostStr:    baseTime.ParseDurationOfTimeForViewEn(startTime, time.Now()),
			Cost:       time.Now().Sub(startTime).Milliseconds(),
		}

		statusCode := c.Writer.Status()
		// 1xx和2xx都是成功
		if (statusCode >= 300) && statusCode != 0 {
			if statusCode == 401 {
				// token失败重新登录的这个就不打印了，太多了
				return
			}
			// 失败必打印
			var response ResponseBase
			if err := json.Unmarshal([]byte(blw.body.String()), &response); err != nil {
				return
			}
			responseMessage.Response = response
			logger.Errorf("调用异常：%v", util.ToJsonString(responseMessage))
		} else {
			// 成功只有在开启开关的情况下才打印
			if rspPrint {
				var response ResponseBase
				if err := json.Unmarshal([]byte(blw.body.String()), &response); err != nil {
					return
				} else {
					responseMessage.Response = response
					logger.Group("http.server.rsp").Debugf("响应：%v", util.ToJsonString(responseMessage))
				}
			}
		}
	}
}

//
//func printReq(requestUri string, requestData Request) {
//	logger.Group("gin.req").Debugf("请求：%v", util.ToJsonString(requestData))
//}
//
//func printRsq(requestUri string, responseMessage Response) {
//	includeUri := config.GetValueArray("gole.server.http.response.print.include-uri")
//	printFlag := false
//	if len(includeUri) != 0 {
//		for _, uri := range includeUri {
//			if uri == "*" || strings.HasPrefix(requestUri, util.ToString(uri)) {
//				printFlag = true
//				break
//			}
//		}
//	}
//
//	excludeUri := config.GetValueArray("gole.server.http.response.print.exclude-uri")
//	if len(excludeUri) != 0 {
//		for _, uri := range excludeUri {
//			if strings.HasPrefix(requestUri, util.ToString(uri)) {
//				printFlag = false
//				break
//			}
//		}
//	}
//
//	rspLogLevel := config.GetValueStringDefault("gole.server.http.response.print.level", "info")
//	if printFlag {
//		logger.Record(rspLogLevel, "响应：%v", util.ToJsonString(responseMessage))
//	}
//}

type Request struct {
	Method     string      `form:"method"`
	Uri        string      `form:"uri"`
	Ip         string      `form:"ip"`
	Headers    http.Header `form:"headers"`
	Parameters gin.Params  `form:"parameters,omitempty"`
	Body       any         `form:"body,omitempty"`
}

type Response struct {
	Request    Request      `json:"request"`
	Response   ResponseBase `json:"response"`
	Cost       int64        `json:"cost"`
	CostStr    string       `json:"costStr"`
	StatusCode int          `json:"statusCode"`
}
