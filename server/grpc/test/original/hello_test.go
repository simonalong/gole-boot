package original

import (
	"context"
	"fmt"
	service2 "github.com/simonalong/gole-boot/server/grpc/test/original/service"
	"github.com/simonalong/gole/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
	"time"
)

// serverImpl is used to implement helloworld.GreeterServer.
type serverImpl struct {
	service2.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *serverImpl) SayHello(ctx context.Context, in *service2.HelloRequest) (*service2.HelloReply, error) {
	fmt.Println("Received:", in.GetName())
	return &service2.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func TestServer(t *testing.T) {
	// 建立 TCP 连接
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 9090))
	if err != nil {
		fmt.Println("failed to listen: ", err)
	}
	// 创建 gRPC 服务
	s := grpc.NewServer()
	// 注册服务
	s.RegisterService(&service2.Greeter_ServiceDesc, &serverImpl{})

	// 运行服务
	if err := s.Serve(lis); err != nil {
		fmt.Println("failed to serve: ", err)
	}
}

func TestClient(t *testing.T) {
	//flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := service2.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.SayHello(ctx, &service2.HelloRequest{Name: "name1"})
	if err != nil {
		logger.Fatalf("could not greet: %v", err)
	}
	//logger.Info("Greeting: ", r.GetMessage())
}
