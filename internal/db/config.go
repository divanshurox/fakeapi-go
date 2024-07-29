package db

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewEmptyConfig() *Config {
	return &Config{}
}

func (config *Config) WithHost(host string) *Config {
	config.Host = host
	return config
}

func (config *Config) WithPort(port string) *Config {
	config.Port = port
	return config
}

func (config *Config) WithUsername(usr string) *Config {
	config.Username = usr
	return config
}

func (config *Config) WithPassword(pwd string) *Config {
	config.Password = pwd
	return config
}
