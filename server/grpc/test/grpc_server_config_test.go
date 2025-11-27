package test

import (
	"github.com/simonalong/gole-boot/server/grpc"
	"github.com/simonalong/gole/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServerConfig(t *testing.T) {
	config.LoadFile("./application-server-all.yaml")
	var cfgOfGrpcServer grpc.ConfigOfGrpcServer
	config.GetValueObject("gole.server-grpc", &cfgOfGrpcServer)

	assert.Equal(t, 9090, cfgOfGrpcServer.Port)
	assert.Equal(t, uint32(1000), cfgOfGrpcServer.MaxConcurrentStreams)
	assert.Equal(t, 4194304, cfgOfGrpcServer.MaxReceiveMessageSize)
	assert.Equal(t, 4194304, cfgOfGrpcServer.MaxSendMessageSize)
	assert.Equal(t, "12s", cfgOfGrpcServer.KeepaliveParams.MaxConnectionIdle.String())
	assert.Equal(t, "10h0m0s", cfgOfGrpcServer.KeepaliveParams.MaxConnectionAge.String())
	assert.Equal(t, "12h0m0s", cfgOfGrpcServer.KeepaliveParams.MaxConnectionAgeGrace.String())
	assert.Equal(t, "2h0m0s", cfgOfGrpcServer.KeepaliveParams.Time.String())
	assert.Equal(t, "20s", cfgOfGrpcServer.KeepaliveParams.Timeout.String())
	assert.Equal(t, int32(65536), cfgOfGrpcServer.InitialWindowSize)
	assert.Equal(t, int32(65536), cfgOfGrpcServer.InitialConnWindowSize)
	assert.Equal(t, 32768, cfgOfGrpcServer.WriteBufferSize)
	assert.Equal(t, 32768, cfgOfGrpcServer.ReadBufferSize)
	assert.Equal(t, true, cfgOfGrpcServer.SharedWriteBuffer)
	assert.Equal(t, "2m0s", cfgOfGrpcServer.ConnectionTimeout.String())
	assert.Equal(t, uint32(512), *cfgOfGrpcServer.MaxHeaderListSize)
	assert.Equal(t, uint32(512), *cfgOfGrpcServer.HeaderTableSize)
	assert.Equal(t, uint32(100), cfgOfGrpcServer.NumServerWorkers)
	assert.Equal(t, true, cfgOfGrpcServer.WaitForHandlers)
}
