package test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/panjf2000/gnet/v2"
	"github.com/simonalong/gole-boot/server/tcp"
	"github.com/simonalong/gole/logger"
	"net"
	"testing"
	"time"
)

type ClientMsgCodec struct {
}

func (codec *ClientMsgCodec) Decode(c gnet.Conn) (interface{}, error) {
	tagAndLengthSize := tagSize + lengthSize
	buf, _ := c.Peek(tagAndLengthSize)
	// 校验head：长度
	if len(buf) < tagAndLengthSize {
		return nil, errors.New("incomplete packet")
	}
	// 校验head：数据
	if !bytes.Equal(magicNumberBytes, buf[:tagSize]) {
		return nil, errors.New("invalid magic number")
	}

	// 读取数据内容
	bodyLen := binary.BigEndian.Uint32(buf[tagSize:tagAndLengthSize])
	msgLen := tagAndLengthSize + int(bodyLen)
	if c.InboundBuffered() < msgLen {
		return nil, errors.New("incomplete packet")
	}
	buf, _ = c.Peek(msgLen)
	_, _ = c.Discard(msgLen)

	return buf[tagAndLengthSize:msgLen], nil
}

//func (codec *ClientMsgCodec) Encode(bufData interface{}) ([]byte, error) {
//	buf := bufData.([]byte)
//	bodyOffset := tagSize + lengthSize
//	msgLen := bodyOffset + len(buf)
//	data := make([]byte, msgLen)
//
//	// 写入magic
//	binary.BigEndian.PutUint16(data[:tagSize], uint16(magicNumber))
//
//	// 写入length
//	binary.BigEndian.PutUint32(data[tagSize:bodyOffset], uint32(len(buf)))
//
//	// 写入body
//	copy(data[bodyOffset:msgLen], buf)
//	return data, nil
//}

type ClientDemoReceiver struct {
	gnet.BuiltinEventEngine
}

func (handler *ClientDemoReceiver) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	c.SetContext(&ClientMsgCodec{})
	logger.Infof("收到连接：%s:%s", c.RemoteAddr().Network(), c.RemoteAddr().String())
	return
}

func (handler *ClientDemoReceiver) OnTraffic(c gnet.Conn) (action gnet.Action) {
	codec := c.Context().(tcp.Decoder)
	if codec == nil {
		return gnet.Close
	}
	req, err := codec.Decode(c)
	if err != nil {
		logger.Errorf("解码异常：%v", err)
		return gnet.None
	}
	logger.Infof("tcp客户端收到解析完后的消息：%v", string(req.([]byte)))
	return gnet.None
}

// 使用gnet进行客户端数据发送
func TestTcpClient(t *testing.T) {
	client, err := gnet.NewClient(&ClientDemoReceiver{})
	if err != nil {
		logger.Fatalf("启动client失败：%v", err)
		return
	}
	client.Start()
	defer client.Stop()

	con, err := client.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		logger.Fatalf("启动client失败：%v", err)
		return
	}
	_, err = con.Write([]byte("服务端你好"))
	if err != nil {
		logger.Fatalf("写入失败：%v", err)
		return
	}

	time.Sleep(1000 * time.Hour)
}

// 使用go/net进行客户端数据发送
func TestGoNetClient(t *testing.T) {
	// 连接到服务端，假设服务端运行在本地的8080端口
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer conn.Close()

	// 发送数据
	conn.Write([]byte("Hello world!"))

	time.Sleep(13 * time.Second)
}
