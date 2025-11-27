package test

import (
	"testing"

	httpServer "github.com/simonalong/gole-boot/server/http"
)

// gole.profiles.active=debug
func TestDebug(t *testing.T) {
	httpServer.RunServer()
}
