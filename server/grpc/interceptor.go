package grpc

import (
	"context"
	"errors"
	"github.com/simonalong/gole-boot/errorx"
	"github.com/simonalong/gole/global"
	"google.golang.org/grpc"
)

func TraceUnaryInterceptorOfServer() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		global.SetGlobalContext(ctx)
		resp, err := handler(ctx, req)
		return resp, err
	}
}

func ErrorUnaryInterceptorOfServer() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}
		var rspErr *errorx.BaseError
		if !errors.As(err, &rspErr) {
			// 不是标准的异常，则框架层统一封装掉
			rspErr = errorx.SC_SERVER_ERROR.WithError(err)
		}
		return resp, errorx.BaseErrToStatusErr(rspErr)
	}
}

func ErrorStreamInterceptorOfServer() func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, ss)
		if err == nil {
			return nil
		}
		var rspErr *errorx.BaseError
		if !errors.As(err, &rspErr) {
			// 不是标准的异常，则框架层统一封装掉
			rspErr = errorx.SC_SERVER_ERROR.WithError(err)
		}
		return errorx.BaseErrToStatusErr(rspErr)
	}
}
