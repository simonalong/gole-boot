package emqx

import (
	"errors"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/simonalong/gole-boot/otel"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/logger"
	t0 "time"
)

func init() {
	config.Load()

	if config.Loaded && config.GetValueBoolDefault("gole.emqx.enable", false) {
		err := config.GetValueObject("gole.emqx", &Cfg)
		if err != nil {
			logger.Warnf("读取emqx配置异常, %v", err.Error())
			return
		}

		mqtt.DEBUG = emqxLogger{"DEBUG"}
		mqtt.WARN = emqxLogger{"WARN"}
		mqtt.CRITICAL = emqxLogger{"CRITICAL"}
		mqtt.ERROR = emqxLogger{"ERROR"}
	}
}

type emqxLogger struct {
	Level string
}

func (log emqxLogger) Println(v ...interface{}) {
	switch log.Level {
	case "DEBUG":
		logger.Debug(v...)
	case "WARN":
		logger.Warn(v...)
	case "CRITICAL":
		logger.Error(v...)
	case "ERROR":
		logger.Error(v...)
	}
}
func (log emqxLogger) Printf(format string, v ...interface{}) {
	switch log.Level {
	case "DEBUG":
		logger.Debugf(format, v...)
	case "WARN":
		logger.Warnf(format, v...)
	case "CRITICAL":
		logger.Errorf(format, v...)
	case "ERROR":
		logger.Errorf(format, v...)
	}
}

func NewEmqxClient() (mqtt.Client, error) {
	if !config.GetValueBoolDefault("gole.emqx.enable", false) {
		logger.Error("emqx没有配置，请先配置")
		return nil, errors.New("emqx没有配置，请先配置")
	}

	_emqxClient := mqtt.NewClient(localEmqxOptions())
	if token := _emqxClient.Connect(); token.Wait() && token.Error() != nil {
		logger.Errorf("链接emqx client失败, %v", token.Error().Error())
		return nil, token.Error()
	}
	return _emqxClient, nil
}

func localEmqxOptions() *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	for _, server := range Cfg.Servers {
		opts.AddBroker(server)
	}

	if Cfg.ClientId != "" {
		opts.SetClientID(Cfg.ClientId)
	}

	if Cfg.Username != "" {
		opts.SetUsername(Cfg.Username)
	}

	if Cfg.Password != "" {
		opts.SetPassword(Cfg.Password)
	}

	if Cfg.CleanSession != true && config.GetValueString("gole.emqx.clean-session") != "" {
		opts.SetCleanSession(Cfg.CleanSession)
	}

	if Cfg.Order != true && config.GetValueString("gole.emqx.order") != "" {
		opts.SetOrderMatters(Cfg.Order)
	}

	if Cfg.WillEnabled != false && config.GetValueString("gole.emqx.will-enabled") != "" {
		opts.WillEnabled = Cfg.WillEnabled
	}

	if Cfg.WillTopic != "" {
		opts.WillTopic = Cfg.WillTopic
	}

	if Cfg.WillQos != 0 {
		opts.WillQos = Cfg.WillQos
	}

	if Cfg.WillRetained != false && config.GetValueString("gole.emqx.will-retained") != "" {
		opts.WillRetained = Cfg.WillRetained
	}

	if Cfg.ProtocolVersion != 0 {
		opts.ProtocolVersion = Cfg.ProtocolVersion
	}

	if Cfg.KeepAlive != 30 && config.GetValueString("gole.emqx.keep-alive") != "" {
		opts.KeepAlive = Cfg.KeepAlive
	}

	if Cfg.PingTimeout != "10s" && config.GetValueString("gole.emqx.ping-timeout") != "" {
		t, err := t0.ParseDuration(Cfg.PingTimeout)
		if err != nil {
			logger.Warnf("读取配置【gole.emqx.ping-timeout】异常：%v", err.Error())
		} else {
			opts.PingTimeout = t
		}
	}

	if Cfg.ConnectTimeout != "30s" && config.GetValueString("gole.emqx.connect-timeout") != "" {
		t, err := t0.ParseDuration(Cfg.ConnectTimeout)
		if err != nil {
			logger.Warnf("读取配置【gole.emqx.connect-timeout】异常：%v", err.Error())
		} else {
			opts.PingTimeout = t
		}
	}

	if Cfg.MaxReconnectInterval != "10m" && config.GetValueString("gole.emqx.max-reconnect-interval") != "" {
		t, err := t0.ParseDuration(Cfg.MaxReconnectInterval)
		if err != nil {
			logger.Warnf("读取配置【gole.emqx.max-reconnect-interval】异常：%v", err.Error())
		} else {
			opts.MaxReconnectInterval = t
		}
	}

	if Cfg.AutoReconnect != true && config.GetValueString("gole.emqx.auto-reconnect") != "" {
		opts.AutoReconnect = Cfg.AutoReconnect
	}

	if Cfg.ConnectRetryInterval != "30s" && config.GetValueString("gole.emqx.connect-retry-interval") != "" {
		t, err := t0.ParseDuration(Cfg.ConnectRetryInterval)
		if err != nil {
			logger.Warnf("读取配置【gole.emqx.connect-retry-interval】异常：%v", err.Error())
		} else {
			opts.ConnectRetryInterval = t
		}
	}

	if Cfg.ConnectRetry != false && config.GetValueString("gole.emqx.connect-retry") != "" {
		opts.ConnectRetry = Cfg.ConnectRetry
	}

	if Cfg.WriteTimeout != "0" && config.GetValueString("gole.emqx.write-timeout") != "" {
		t, err := t0.ParseDuration(Cfg.WriteTimeout)
		if err != nil {
			logger.Warnf("读取配置【gole.emqx.write-timeout】异常：%v", err.Error())
		} else {
			opts.WriteTimeout = t
		}
	}

	if Cfg.ResumeSubs != false && config.GetValueString("gole.emqx.resume-subs") != "" {
		opts.ResumeSubs = Cfg.ResumeSubs
	}

	if Cfg.MaxResumePubInFlight != 0 {
		opts.MaxResumePubInFlight = Cfg.MaxResumePubInFlight
	}

	if Cfg.AutoAckDisabled != false && config.GetValueString("gole.emqx.auto-ack-disabled") != "" {
		opts.AutoAckDisabled = Cfg.AutoAckDisabled
	}

	return opts
}
