package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

const configEnvPrefix = "server"

type Config struct {
	Port            int    `default:":8080"`
	TargetPrefix    string `required:"true"`
	ChallengeLength int    `default:"8"`
}

func (c Config) Host() string {
	return fmt.Sprintf("0.0.0.0:%d", c.Port)
}

func ParseConfig() (*Config, error) {
	cfg := &Config{}

	if err := envconfig.Process(configEnvPrefix, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
