
# server
主要用于启动多服务，当前支持：grpc、http、tcp服务端同时启动

## 配置多类型服务启动
```yaml
gole:
  application:
    name: sample
  server:
    # http 服务
    http:
      enable: true
      port: 8180
    # tcp 服务
    tcp:
      enable: true
      port: 9000
    # grpc 服务
    grpc:
      enable: true
      port: 9090

```
代码
```go
import "github.com/simonalong/gole-boot/server"

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
    httpServer.RegisterRoute("/api/data", httpServer.HmGet, func(c *gin.Context) (any, error) {
        return "ok", nil
    })
    
    server.Run()
}
```
