package original

import (
	"context"
	"fmt"
	service2 "github.com/simonalong/gole-boot/server/grpc/test/original/service"
	"github.com/simonalong/gole/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net"
	"testing"
	"time"
)

func tracerProviderServer(serviceName string) error {
	ctx := context.Background()

	var exporter *otlptrace.Exporter
	// 创建 OTLP 导出器；导出到otel-collector
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint("localhost:4317"),
		// 如果你使用的是自签名证书，这里可以设置为 WithInsecure()
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))),
	)
	if err != nil {
		logger.Fatalf("创建exporter失败: %v", err)
		return err
	}

	resources, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		logger.Error("不能设置资源: ", err)
		return err
	}

	// 创建TracerProvider
	otel.SetTracerProvider(sdkTrace.NewTracerProvider(
		sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
		sdkTrace.WithBatcher(exporter),
		sdkTrace.WithResource(resources),
	))

	// 配置多进程传播
	b3Propagator := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}, b3Propagator)
	otel.SetTextMapPropagator(propagator)
	return nil
}

func tracerProviderClient(serviceName string) error {
	ctx := context.Background()

	var exporter *otlptrace.Exporter
	// 创建 OTLP 导出器；导出到otel-collector
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint("localhost:4317"),
		// 如果你使用的是自签名证书，这里可以设置为 WithInsecure()
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))),
	)
	if err != nil {
		logger.Fatalf("创建exporter失败: %v", err)
		return err
	}

	resources, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		logger.Error("不能设置资源: ", err)
		return err
	}

	// 创建TracerProvider
	otel.SetTracerProvider(sdkTrace.NewTracerProvider(
		sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
		sdkTrace.WithBatcher(exporter),
		sdkTrace.WithResource(resources),
	))

	// 配置多进程传播
	b3Propagator := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}, b3Propagator)
	otel.SetTextMapPropagator(propagator)
	return nil
}

func TestServerOtel(t *testing.T) {
	_ = tracerProviderServer("grpc-server")

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	s.RegisterService(&service2.Greeter_ServiceDesc, &serverImpl{})
	//pb.RegisterGreeterServer(s, &server{})
	logger.Infof("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}

func TestClientOtel(t *testing.T) {
	_ = tracerProviderClient("grpc-client")

	conn, err := grpc.NewClient("localhost:9090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)

	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := service2.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &service2.HelloRequest{Name: "name1"})
	if err != nil {
		s, ok := status.FromError(err)
		if !ok {
			logger.Fatalf("err is not standard grpc error: %v", err)
		}
		fmt.Println(s.Code())
		//log.Fatalf("could not greet: %v", err)
	}
	logger.Infof("Greeting: %s", r.GetMessage())
	time.Sleep(5 * time.Second)
}
