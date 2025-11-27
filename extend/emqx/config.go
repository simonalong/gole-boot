package emqx

var Cfg EmqxConfig

type EmqxServer struct {
	User *EmqxServerUser // username and password information
	Host string          // host or host:port
}

type EmqxServerUser struct {
	Username string
	Password string
}

type EmqxConfig struct {
	Servers              []string
	ClientId             string
	Username             string
	Password             string
	CleanSession         bool
	Order                bool
	WillEnabled          bool
	WillTopic            string
	WillQos              byte
	WillRetained         bool
	ProtocolVersion      uint
	KeepAlive            int64
	PingTimeout          string
	ConnectTimeout       string
	MaxReconnectInterval string
	AutoReconnect        bool
	ConnectRetryInterval string
	ConnectRetry         bool
	WriteTimeout         string
	ResumeSubs           bool
	MaxResumePubInFlight int
	AutoAckDisabled      bool
}
