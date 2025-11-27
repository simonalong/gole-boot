package original

import (
	"context"
	"fmt"
	service2 "github.com/simonalong/gole-boot/server/grpc/test/original/service"
	"github.com/simonalong/gole/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"net"
	"strconv"
	"testing"
	"time"
)

const (
	exampleScheme      = "base"
	exampleServiceName = "demo-service"
)

// var addrs = []string{"localhost:9091", "localhost:9092", "localhost:9093"}
var addrs = []string{"192.168.0.56:9091", "192.168.0.56:9092", "192.168.0.56:9093"}

type serverImpl1 struct {
	service2.UnimplementedGreeterServer
}
type serverImpl2 struct {
	service2.UnimplementedGreeterServer
}
type serverImpl3 struct {
	service2.UnimplementedGreeterServer
}

func (s *serverImpl1) SayHello(ctx context.Context, in *service2.HelloRequest) (*service2.HelloReply, error) {
	fmt.Println("Received:", in.GetName())
	return &service2.HelloReply{Message: "Hello " + in.GetName() + "from server1"}, nil
}
func (s *serverImpl2) SayHello(ctx context.Context, in *service2.HelloRequest) (*service2.HelloReply, error) {
	fmt.Println("Received:", in.GetName())
	return &service2.HelloReply{Message: "Hello " + in.GetName() + "from server2"}, nil
}
func (s *serverImpl3) SayHello(ctx context.Context, in *service2.HelloRequest) (*service2.HelloReply, error) {
	fmt.Println("Received:", in.GetName())
	return &service2.HelloReply{Message: "Hello " + in.GetName() + "from server3"}, nil
}

func TestServer1(t *testing.T) {
	// 建立 TCP 连接
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 9091))
	if err != nil {
		fmt.Println("failed to listen: ", err)
	}
	// 创建 gRPC 服务
	s := grpc.NewServer()
	// 注册服务
	s.RegisterService(&service2.Greeter_ServiceDesc, &serverImpl1{})

	// 运行服务
	if err := s.Serve(lis); err != nil {
		fmt.Println("failed to serve: ", err)
	}
}

func TestServer2(t *testing.T) {
	// 建立 TCP 连接
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 9092))
	if err != nil {
		fmt.Println("failed to listen: ", err)
	}
	// 创建 gRPC 服务
	s := grpc.NewServer()
	// 注册服务
	s.RegisterService(&service2.Greeter_ServiceDesc, &serverImpl2{})

	// 运行服务
	if err := s.Serve(lis); err != nil {
		fmt.Println("failed to serve: ", err)
	}
}

func TestServer3(t *testing.T) {
	// 建立 TCP 连接
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 9093))
	if err != nil {
		fmt.Println("failed to listen: ", err)
	}
	// 创建 gRPC 服务
	s := grpc.NewServer()
	// 注册服务
	s.RegisterService(&service2.Greeter_ServiceDesc, &serverImpl3{})

	// 运行服务
	if err := s.Serve(lis); err != nil {
		fmt.Println("failed to serve: ", err)
	}
}

func TestClientLb(t *testing.T) {
	//flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(
		"gole:///demo-service",
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := service2.NewGreeterClient(conn)

	for i := 0; i < 10000; i++ {
		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		rspMsg, err := c.SayHello(ctx, &service2.HelloRequest{Name: "name " + strconv.Itoa(i)})
		if err != nil {
			logger.Fatalf("could not greet: %v", err)
		}
		logger.Info("Greeting: ", rspMsg.GetMessage())
		time.Sleep(2 * time.Second)
	}
}

// Following is an example name resolver implementation. Read the name
// resolution example to learn more about it.

type exampleResolverBuilder struct{}

func (*exampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &exampleResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			exampleServiceName: addrs,
		},
	}
	r.start()
	return r, nil
}
func (*exampleResolverBuilder) Scheme() string { return exampleScheme }

type exampleResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *exampleResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*exampleResolver) Close()                                  {}

func init() {
	resolver.Register(&exampleResolverBuilder{})
}
