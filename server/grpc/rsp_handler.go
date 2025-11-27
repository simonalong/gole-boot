package grpc

import (
	"context"
	"errors"
	"github.com/simonalong/gole-boot/errorx"
	"github.com/simonalong/gole/logger"
	baseTime "github.com/simonalong/gole/time"
	"github.com/simonalong/gole/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"time"
)

type Request struct {
	Context context.Context `json:"-"`
	Method  string          `json:"method"`
	Ip      string          `json:"ip"`
	Request interface{}     `json:"request"`
}

type Response struct {
	Request    *Request          `json:"request"`
	Error      *errorx.BaseError `json:"error,omitempty"`
	Cost       int64             `json:"cost"`
	CostStr    string            `json:"costStr"`
	StatusCode uint32            `json:"statusCode"`
	Data       interface{}       `json:"data,omitempty"`
}

func generateReq(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo) *Request {
	var clientIP string
	p, ok := peer.FromContext(ctx)
	if ok {
		clientIP = p.Addr.String()
	}
	return &Request{
		Context: ctx,
		Method:  info.FullMethod,
		Ip:      clientIP,
		Request: req,
	}
}

func generateRsp(startTime time.Time, req *Request, res interface{}, baseErr *errorx.BaseError) *Response {
	if baseErr == nil {
		return &Response{
			Request: req,
			CostStr: baseTime.ParseDurationOfTimeForViewEn(startTime, time.Now()),
			Cost:    time.Now().Sub(startTime).Milliseconds(),
			Data:    res,
		}
	}

	return &Response{
		Request:    req,
		Error:      baseErr,
		CostStr:    baseTime.ParseDurationOfTimeForViewEn(startTime, time.Now()),
		Cost:       time.Now().Sub(startTime).Milliseconds(),
		StatusCode: errorx.GetGrpcStatus(baseErr.Code),
		Data:       res,
	}
}

func RspHandleUnaryInterceptorOfServer() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		request := generateReq(ctx, req, info)

		resp, err := handler(ctx, req)
		var baseErr *errorx.BaseError
		if err != nil && !errors.As(err, &baseErr) {
			// 不是标准的异常，则框架层统一封装掉
			baseErr = errorx.SC_SERVER_ERROR.WithError(err)
		}

		response := generateRsp(startTime, request, resp, baseErr)

		// 异常：记录错误日志
		if !errorx.IsOk(baseErr) {
			logger.Errorf("调用异常：%v", util.ToJsonString(response))
			return resp, baseErr
		}

		// 1s以内按照分组debug
		if response.Cost < time.Second.Milliseconds() {
			logger.Group("grpc-server").Debugf("调用信息：%v", util.ToJsonString(response))
			return resp, baseErr
		}

		// 超过1s但是在10s以内则发告警日志
		if response.Cost < 10*time.Second.Milliseconds() {
			logger.Warnf("调用时间超过1s：%v", util.ToJsonString(response))
			return resp, baseErr
		}

		// 超过10s则发异常日志
		logger.Errorf("调用时间超过10s：%v", util.ToJsonString(response))
		return resp, baseErr
	}
}
