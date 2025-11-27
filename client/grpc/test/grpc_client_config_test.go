package test

import (
	"github.com/simonalong/gole-boot/client/grpc"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClientConfig(t *testing.T) {
	config.LoadFile("./application-client-all.yaml")
	cfgOfGrpcClient := grpc.ConfigOfGrpcClient{}
	err := config.GetValueObject("gole.grpc.demo1", &cfgOfGrpcClient)
	if err != nil {
		logger.Errorf("异常：%v", err)
		return
	}

	assert.Equal(t, "localhost", cfgOfGrpcClient.Host)
	assert.Equal(t, 9090, cfgOfGrpcClient.Port)
	assert.Equal(t, "test", cfgOfGrpcClient.Authority)
	assert.Equal(t, true, cfgOfGrpcClient.DisableServiceConfig)
	assert.Equal(t, true, cfgOfGrpcClient.DisableRetry)
	assert.Equal(t, true, cfgOfGrpcClient.DisableHealthCheck)
	assert.Equal(t, "json", *cfgOfGrpcClient.DefaultServiceConfigRawJSON)
	assert.Equal(t, "30m0s", cfgOfGrpcClient.IdleTimeout.String())
	assert.Equal(t, "10s", cfgOfGrpcClient.ConnectParams.MinConnectTimeout.String())
	assert.Equal(t, "10s", cfgOfGrpcClient.ConnectParams.Backoff.BaseDelay.String())
	assert.Equal(t, 10.0, cfgOfGrpcClient.ConnectParams.Backoff.Multiplier)
	assert.Equal(t, 10.0, cfgOfGrpcClient.ConnectParams.Backoff.Jitter)
	assert.Equal(t, "10s", cfgOfGrpcClient.ConnectParams.Backoff.MaxDelay.String())
}
