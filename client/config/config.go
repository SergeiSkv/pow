package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

const configEnvPrefix = "client"

type Config struct {
	ServerHost   string `default:"pow-server"`
	ServerPort   int    `default:"8080"`
	TargetPrefix string `required:"true"`
}

func (c Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.ServerHost, c.ServerPort)
}

func ParseConfig() (*Config, error) {
	cfg := &Config{}

	if err := envconfig.Process(configEnvPrefix, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
