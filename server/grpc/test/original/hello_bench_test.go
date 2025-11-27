package original

import (
	"context"
	service2 "github.com/simonalong/gole-boot/server/grpc/test/original/service"
	"github.com/simonalong/gole/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

// serverImpl is used to implement helloworld.GreeterServer.

var client service2.GreeterClient

func init() {
	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()
	c := service2.NewGreeterClient(conn)
	// Contact the server and print out its response.
	client = c
}

func TestOnce(t *testing.T) {

}

func BenchmarkClient(b *testing.B) {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	for i := 0; i < b.N; i++ {
		//for i := 0; i < 1; i++ {
		_, err := client.SayHello(context.Background(), &service2.HelloRequest{Name: "name1"})
		if err != nil {
			logger.Fatalf("could not greet: %v", err)
		}
	}
}
