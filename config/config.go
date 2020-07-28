package config

type LogConfig interface {
	LogFile() string
	IsDebug() bool
	ServiceID() string
}

type Config struct {
}

func NewConfig() LogConfig {
	return Config{}
}

func (c Config) LogFile() string {
	return "json.json"
}

func (c Config) IsDebug() bool {
	return true
}

func (c Config) ServiceID() string {
	return "server-test"
}
