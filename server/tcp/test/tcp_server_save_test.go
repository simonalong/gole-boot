package test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/panjf2000/gnet/v2"
	"github.com/simonalong/gole-boot/server/tcp"
	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/time"
	"github.com/simonalong/gole/util"
	"testing"
)

var consList util.BsList[*gnet.Conn]
var ConnectMap map[string]*gnet.Conn

func init() {
	magicNumberBytes = make([]byte, tagSize)
	binary.BigEndian.PutUint16(magicNumberBytes, uint16(magicNumber))
	util.NewList[gnet.Conn]()
	ConnectMap = map[string]*gnet.Conn{}
}

type ServerMsgDemoSaveConCodec struct {
}

func (codec *ServerMsgDemoSaveConCodec) Decode(c gnet.Conn) (interface{}, error) {
	consList.Add(&c)
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

//func (codec *ServerMsgDemoSaveConCodec) Encode(bufData interface{}) ([]byte, error) {
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

func TestTcpSaveConServer(t *testing.T) {
	// 设置编码解码器
	tcp.SetDecoder(func() tcp.Decoder { return &ServerMsgDemoSaveConCodec{} })

	// 设置连接开启情况下
	tcp.ConnectOpen(func(c gnet.Conn) {
		ConnectMap[c.RemoteAddr().Network()+":"+c.RemoteAddr().String()] = &c
	})

	// 设置连接断开情况下
	tcp.ConnectClose(func(c gnet.Conn) error {
		delete(ConnectMap, c.RemoteAddr().Network()+":"+c.RemoteAddr().String())
		return nil
	})

	// 设置数据接收处理器
	tcp.Receive(func(dataReq interface{}) ([]byte, error) {
		// 收到消息
		logger.Infof("收到客户端消息：%v", string(dataReq.([]byte)))

		// 自己的业务代码...

		logger.Info("给客户端返回消息：", "收到了")

		// 返回响应
		return []byte("收到了"), nil
	})

	timer := time.NewTimerWithFire(5, func(t *time.Timer) {
		for index, con := range consList {
			logger.Info("给客户端发送数据")
			err := tcp.Send(*con, []byte("hello"))
			if err != nil {
				if errors.Is(err, tcp.ConnectIsUnAvailableErr) {
					consList.Delete(index)
				}
				logger.Error(err)
				continue
			}
		}
	})
	timer.Start()

	// 启动服务
	tcp.RunServer()
}
