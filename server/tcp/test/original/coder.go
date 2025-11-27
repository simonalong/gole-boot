package original

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/panjf2000/gnet/v2"
)

var magicNumberBytes []byte

func init() {
	magicNumberBytes = make([]byte, tagSize)
	binary.BigEndian.PutUint16(magicNumberBytes, uint16(magicNumber))
}

// DemoXxCodec Protocol format:
//
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
type DemoXxCodec struct{}

// All current protocols.
const (
	magicNumber = 1314
	tagSize     = 2
	lengthSize  = 4
)

func (codec DemoXxCodec) Encode(buf []byte) (interface{}, error) {
	bodyOffset := tagSize + lengthSize
	msgLen := bodyOffset + len(buf)
	data := make([]byte, msgLen)

	// 写入magic
	binary.BigEndian.PutUint16(data[:tagSize], uint16(magicNumber))

	// 写入length
	binary.BigEndian.PutUint32(data[tagSize:bodyOffset], uint32(len(buf)))

	// 写入body
	copy(data[bodyOffset:msgLen], buf)
	return data, nil
}

func (codec *DemoXxCodec) Decode(c gnet.Conn) ([]byte, error) {
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
