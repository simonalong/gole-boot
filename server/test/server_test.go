package test

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/panjf2000/gnet/v2"
	"github.com/simonalong/gole-boot/server"
	baseGrpc "github.com/simonalong/gole-boot/server/grpc"
	service2 "github.com/simonalong/gole-boot/server/grpc/test/original/service"
	httpServer "github.com/simonalong/gole-boot/server/http"
	"github.com/simonalong/gole-boot/server/http/rsp"
	"github.com/simonalong/gole-boot/server/tcp"
	"github.com/simonalong/gole/logger"
	"testing"
)

type ServerImpl struct {
	service2.UnimplementedGreeterServer
}

func (s *ServerImpl) SayHello(ctx context.Context, in *service2.HelloRequest) (*service2.HelloReply, error) {
	fmt.Println("Received: ", in.GetName())
	return &service2.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// gole.profiles.active=server-grpc-http
// grpc和http同时配置的话，会怎么启动
func TestServerGrpcAndHttp(t *testing.T) {

	// 注册服务
	baseGrpc.Register(&service2.Greeter_ServiceDesc, &ServerImpl{})
	// 启动http服务
	httpServer.AddGinRoute("/api/data", httpServer.HmGet, func(c *gin.Context) {
		rsp.Done(c, "ok")
	})

	server.Run()
}

type ServerMsgDemoSaveConCodec struct {
}

func (codec *ServerMsgDemoSaveConCodec) Decode(c gnet.Conn) (interface{}, error) {
	return []byte("阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿是老地方阿斯蒂芬爱仕达发斯蒂芬sdf阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬阿斯蒂芬"), nil
}

//func (codec *ServerMsgDemoSaveConCodec) Encode(buf interface{}) ([]byte, error) {
//	return nil, nil
//}

// gole.profiles.active=server-grpc-http-tcp
// grpc和http同时配置的话，会怎么启动
func TestServer(t *testing.T) {

	// tcp：设置编码解码器
	tcp.SetDecoder(func() tcp.Decoder { return &ServerMsgDemoSaveConCodec{} })
	tcp.Receive(func(msg interface{}) ([]byte, error) {
		// 自己的业务代码...
		logger.Info("给客户端返回消息：", "收到了")
		return nil, nil
	})

	// grpc：注册服务
	baseGrpc.Register(&service2.Greeter_ServiceDesc, &ServerImpl{})

	// http：启动http服务
	httpServer.AddGinRoute("/api/data", httpServer.HmGet, func(c *gin.Context) {
		rsp.Done(c, "ok")
	})

	server.Run()
}
