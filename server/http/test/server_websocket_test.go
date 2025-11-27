package test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	httpServer "github.com/simonalong/gole-boot/server/http"
)

// gole.profiles.active=debug
func TestWebsocket(t *testing.T) {
	httpServer.Get("ws", func(c *gin.Context) (any, error) {
		handleWebSocket(c)
		return nil, nil
	})
	httpServer.RunServer()
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域请求，生产环境可限制特定域名
	},
}

func handleWebSocket(c *gin.Context) {
	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "升级连接失败"})
		return
	}
	defer conn.Close()

	// 处理 WebSocket 通信逻辑
	for {
		// 接收消息
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			break
		}
		message := string(p)
		// 输出接收到的消息
		println("接收到的消息:", message)
		// 发送消息回客户端
		if err := conn.WriteMessage(messageType, []byte("服务器已收到消息: "+message)); err != nil {
			break
		}
	}
}
