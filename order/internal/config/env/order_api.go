package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type orderApiEnvConfig struct {
	Host            string        `env:"HTTP_HOST,required"`
	Port            string        `env:"HTTP_PORT,required"`
	HTTPReadTimeout time.Duration `env:"HTTP_READ_TIMEOUT,required"`
	ShutDownTimeout time.Duration `env:"ORDER_SHUT_DOWN_TIMEOUT,required"`
}

type orderApiConfig struct {
	raw orderApiEnvConfig
}

func NewOrderApiConfig() (*orderApiConfig, error) {
	var raw orderApiEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderApiConfig{raw: raw}, nil
}

func (cfg *orderApiConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *orderApiConfig) ReadTimeout() time.Duration {
	return cfg.raw.HTTPReadTimeout
}

func (cfg *orderApiConfig) ShutDownTimeout() time.Duration {
	return cfg.raw.ShutDownTimeout
}
