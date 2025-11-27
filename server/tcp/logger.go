package tcp

import (
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/simonalong/gole/logger"
)

func initGnetLogger() {
	logging.SetDefaultLoggerAndFlusher(&GnetLogger{}, logging.GetDefaultFlusher())
}

type GnetLogger struct {
}

func (gnetLog *GnetLogger) Debugf(format string, args ...interface{}) {
	logger.Group("gnet").Debugf(format, args...)
}

func (gnetLog *GnetLogger) Infof(format string, args ...interface{}) {
	logger.Group("gnet").Infof(format, args...)
}

func (gnetLog *GnetLogger) Warnf(format string, args ...interface{}) {
	logger.Group("gnet").Warnf(format, args...)
}

func (gnetLog *GnetLogger) Errorf(format string, args ...interface{}) {
	logger.Group("gnet").Errorf(format, args...)
}

func (gnetLog *GnetLogger) Fatalf(format string, args ...interface{}) {
	logger.Group("gnet").Fatalf(format, args...)
}
