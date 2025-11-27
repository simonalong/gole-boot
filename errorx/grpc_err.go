package errorx

import (
	"github.com/simonalong/gole/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func BaseErrToStatusErr(rspErr *BaseError) error {
	if rspErr == nil {
		return nil
	}
	st, err := status.New(codes.Code(GetGrpcStatus(rspErr.Code)), rspErr.Msg).WithDetails(&GrpcBaseError{
		Code:   rspErr.Code,
		Msg:    rspErr.Msg,
		Detail: rspErr.Detail,
	})
	if err != nil {
		logger.Error("rspErr转换grpc.Status异常")
		return err
	}
	return st.Err()
}

func StatusErrToBaseErr(originalErr error) *BaseError {
	if originalErr == nil {
		return nil
	}
	st, ok := status.FromError(originalErr)
	if !ok {
		return nil
	}
	// detail为空，说明没有进入到服务端，则查看status的code
	if st.Details() == nil || len(st.Details()) == 0 {
		return GetBaseErrFromGrpcStatus(uint32(st.Code())).WithDetail(st.Message())
	}

	// detail非空，说明进入到了服务端
	var grpcBaseErr *GrpcBaseError
	details := st.Details()
	for _, detail := range details {
		switch detail.(type) {
		case *GrpcBaseError:
			grpcBaseErr = detail.(*GrpcBaseError)
			break
		}
	}
	if grpcBaseErr == nil {
		return nil
	}
	return &BaseError{
		Code:   grpcBaseErr.Code,
		Msg:    grpcBaseErr.Msg,
		Detail: grpcBaseErr.Detail,
	}
}
