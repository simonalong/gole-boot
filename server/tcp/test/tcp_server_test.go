package test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/panjf2000/gnet/v2"
	"github.com/simonalong/gole-boot/server/tcp"
	"github.com/simonalong/gole/logger"
	"testing"
)

// * 0           2                       6
// * +-----------+-----------------------+
// * |   magic   |       body len        |
// * +-----------+-----------+-----------+
// * |                                   |
// * +                                   +
// * |           body bytes              |
// * +                                   +
// * |            ... ...                |
// * +-----------------------------------+.
const (
	magicNumber = 1314
	tagSize     = 2
	lengthSize  = 4
)

var magicNumberBytes []byte

func init() {
	magicNumberBytes = make([]byte, tagSize)
	binary.BigEndian.PutUint16(magicNumberBytes, uint16(magicNumber))
}

type ServerMsgDemoCodec struct {
}

func (codec *ServerMsgDemoCodec) Decode(c gnet.Conn) (interface{}, error) {
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

//func (codec *ServerMsgDemoCodec) Encode(bufData interface{}) ([]byte, error) {
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

func TestTcpServer(t *testing.T) {
	// 设置编码解码器
	tcp.SetDecoder(func() tcp.Decoder { return &ServerMsgDemoCodec{} })

	// 设置处理器
	tcp.Receive(func(dataReq interface{}) ([]byte, error) {
		data := dataReq.([]byte)
		// 收到消息
		logger.Infof("收到客户端消息：%v", string(data))

		// 自己的业务代码...

		logger.Info("给客户端返回消息：", "收到了")

		// 返回响应
		return []byte("收到了"), nil
	})

	// 启动服务
	tcp.RunServer()
}
