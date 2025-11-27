package rsp

import (
	"errors"
	"github.com/gin-gonic/gin"
	errorx2 "github.com/simonalong/gole-boot/errorx"
	"net/http"
	"strings"
)

type ResponseBase struct {
	Code   string `json:"code"`
	Data   any    `json:"data,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func Success(ctx *gin.Context, object any) {
	ctx.JSON(http.StatusOK, object)
}

func Fail(ctx *gin.Context, code int, obj any) {
	ctx.JSON(code, obj)
}

func SuccessOfStandard(ctx *gin.Context, v any) {
	Done(ctx, v)
}

func FailOfStandard(ctx *gin.Context, code string, message string) {
	Done(ctx, nil, errorx2.New(code, message))
}

func FailOfStandardDetail(ctx *gin.Context, code, message, detail string) {
	Done(ctx, nil, errorx2.NewWithAll(code, message, detail))
}

func FailOfStandardErr(ctx *gin.Context, err *errorx2.BaseError) {
	Done(ctx, nil, err)
}

type ListWrap struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

// Done 用于返回服务端响应处理结果。
// data 作为响应数据的封装对象，返回给前端。
// message 作为业务逻辑层处理错误时返回的错误信息对象。
// example：
//
//	rsp, err := rpc_clients.UserService.GetUser(context.Background(), req)
//	if err != nil {
//		Done(ctx, nil, err)
//		return
//	}
//	Done(ctx, rsp.User)
//
// if you want to Done a page list result:
//
//	data := response.ListWrap{Total: total, BsList: rsp.Users}
//	Done(ctx, data)
func Done(ctx *gin.Context, data any, errs ...error) {
	var xe *errorx2.BaseError
	if len(errs) > 0 && errs[0] != nil {
		if !errors.As(errs[0], &xe) {
			xe = errorx2.SC_SERVER_ERROR.WithError(errs[0])
		}
	}

	if xe == nil {
		xe = errorx2.SC_OK
	}

	var body = ResponseBase{
		Code:   xe.Code,
		Msg:    xe.Msg,
		Detail: xe.Detail,
	}

	switch xe.Code {
	case "SC_OK", "SC_FOUND", "SC_MOVED_PERMANENTLY":
		body.Data = data
	}
	jsonResponse(ctx, &body, errorx2.GetHttpStatus(xe.Code))
	return
}

func jsonResponse(ctx *gin.Context, body *ResponseBase, status int) {
	if strings.Contains(ctx.GetHeader("User-Agent"), "curl") {
		ctx.IndentedJSON(status, body)
		ctx.Abort()
	} else {
		ctx.AbortWithStatusJSON(status, body)
	}
}
