package otel

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/simonalong/gole-boot/event"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/global"
	"github.com/simonalong/gole/goid"
	"github.com/simonalong/gole/listener"
	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/validate"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var CfgOfOtel ConfigOfOpentelemetry

type ConfigOfOpentelemetry struct {
	// opentelemetry-collector服务的地址，可以为空，为空，则表示不搜集
	ExporterUrl string
	// 服务名，默认为 gole.application.name
	ServiceName string `match:"isUnBlank"`
}

func init() {
	config.Load()

	if config.Loaded && config.GetValueBoolDefault("gole.opentelemetry.enable", true) {
		err := config.GetValueObject("gole.opentelemetry", &CfgOfOtel)
		if err != nil {
			logger.Warn("读取 opentelemetry 配置异常")
			return
		}
	} else {
		return
	}

	if CfgOfOtel.ServiceName == "" {
		if val := config.BaseCfg.Application.Name; val != "" {
			CfgOfOtel.ServiceName = val
		} else {
			logger.Fatalf("gole.application.name 不可为空")
		}
	}

	success, _, errMsg := validate.Check(CfgOfOtel)
	if !success {
		logger.Errorf("opentelemetry 配置异常：%v", errMsg)
		return
	}

	initTracer()

	global.ContextLocalStorage = goid.NewLocalStorage()
}

func initTracer() {
	ctx := context.Background()

	var exporter *otlptrace.Exporter
	// 创建 OTLP 导出器；导出到otel-collector
	if strings.HasPrefix(CfgOfOtel.ExporterUrl, "http://") {
		traceClientHttp := otlptracehttp.NewClient(otlptracehttp.WithEndpoint(CfgOfOtel.ExporterUrl[len("http://"):]), otlptracehttp.WithInsecure())
		otlptracehttp.WithCompression(1)
		_exporter, err := otlptrace.New(ctx, traceClientHttp)
		if err != nil {
			logger.Fatal("创建exporter失败: %v", err)
			return
		}
		exporter = _exporter
	} else {
		if CfgOfOtel.ExporterUrl != "" {
			_exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(CfgOfOtel.ExporterUrl),
				// 如果你使用的是自签名证书，这里可以设置为 WithInsecure()
				otlptracegrpc.WithInsecure(),
				otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))),
			)
			if err != nil {
				logger.Fatal("创建exporter失败: %v", err)
				return
			}
			exporter = _exporter
		}
	}

	resources, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.name", CfgOfOtel.ServiceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		logger.Error("不能设置资源: ", err)
		return
	}

	// 创建TracerProvider
	var traceProvider *sdkTrace.TracerProvider
	if exporter != nil {
		traceProvider = sdkTrace.NewTracerProvider(
			sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
			sdkTrace.WithBatcher(exporter),
			sdkTrace.WithResource(resources),
		)
	} else {
		traceProvider = sdkTrace.NewTracerProvider(
			sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
			sdkTrace.WithResource(resources),
		)
	}
	otel.SetTracerProvider(traceProvider)

	// 配置多进程传播
	b3Propagator := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}, b3Propagator)
	otel.SetTextMapPropagator(propagator)

	// 配置一个Tracer
	global.Tracer = otel.Tracer(CfgOfOtel.ServiceName)

	logger.Info("埋点配置完成")

	// 监听(http)服务退出事件
	listener.AddListener(event.EventOfServerHttpStop, func(event listener.BaseEvent) {
		logger.Warn("http应用关闭完成")
		if exporter != nil {
			err := exporter.Shutdown(context.Background())
			if err != nil {
				logger.Error("关闭异常：", err.Error())
			}
		}
	})

	// 监听(tcp)服务退出事件
	listener.AddListener(event.EventOfServerTcpStop, func(event listener.BaseEvent) {
		logger.Warn("tcp应用关闭完成")
		if exporter != nil {
			err := exporter.Shutdown(context.Background())
			if err != nil {
				logger.Error("关闭异常：", err.Error())
			}
		}
	})

	// 监听(grpc)服务退出事件
	listener.AddListener(event.EventOfServerGrpcStop, func(event listener.BaseEvent) {
		logger.Warn("grpc应用关闭完成")
		if exporter != nil {
			err := exporter.Shutdown(context.Background())
			if err != nil {
				logger.Error("关闭异常：", err.Error())
			}
		}
	})
}

func GlobalContextLoad() gin.HandlerFunc {
	return func(c *gin.Context) {
		if global.GetGlobalContext() != nil {
			global.SetGlobalContext(c.Request.Context())
		}
	}
}
