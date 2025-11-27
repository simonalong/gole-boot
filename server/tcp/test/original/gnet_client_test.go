package original

import (
	"github.com/simonalong/gole/logger"
	"net"
	"testing"
	"time"
)

// 这里使用普通的go/net包进行发送数据
func TestGoNetClient(t *testing.T) {
	// 连接到服务端，假设服务端运行在本地的8080端口
	conn, err := net.Dial("tcp", "127.0.0.1:9992")
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer conn.Close()

	// 对要发送的数据编码
	dst, err := DemoXxCodec{}.Encode([]byte("Hello world!"))

	// 发送数据
	conn.Write(dst.([]byte))

	time.Sleep(13 * time.Second)
}

// 这里使用普通的go/net包进行发送数据
func TestGoNetClient2(t *testing.T) {
	// 连接到服务端，假设服务端运行在本地的8080端口
	conn, err := net.Dial("tcp", "127.0.0.1:9992")
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer conn.Close()

	// 对要发送的数据编码
	dst, err := DemoXxCodec{}.Encode([]byte("Hello world!"))

	// 发送数据
	conn.Write(dst.([]byte))

	time.Sleep(13 * time.Second)
}

// 这里使用普通的go/net包进行发送数据
func TestGoNetClient3(t *testing.T) {
	// 连接到服务端，假设服务端运行在本地的8080端口
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer conn.Close()

	// 对要发送的数据编码
	dst, err := DemoXxCodec{}.Encode([]byte("Hello world!"))

	// 发送数据
	conn.Write(dst.([]byte))

	time.Sleep(13 * time.Second)
}
