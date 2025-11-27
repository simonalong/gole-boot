package grpc

import (
	"context"
	"github.com/simonalong/gole-boot/errorx"
	"github.com/simonalong/gole/global"
	"google.golang.org/grpc"
)

func TraceUnaryInterceptorOfClient(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return invoker(global.GetGlobalContext(), method, req, reply, cc, opts...)
}

func ErrorUnaryInterceptorOfClient(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(global.GetGlobalContext(), method, req, reply, cc, opts...)
	if err == nil {
		return nil
	}
	newErr := errorx.StatusErrToBaseErr(err)
	if newErr == nil {
		newErr = errorx.SC_GRPC_ERR.WithError(err)
	}
	return newErr
}
