package tcp

import "github.com/panjf2000/gnet/v2"

type Decoder interface {
	Decode(c gnet.Conn) (interface{}, error)
}

type CodecGenerate func() Decoder
