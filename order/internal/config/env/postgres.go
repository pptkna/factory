package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	Host               string `env:"POSTGRES_HOST,required"`
	Port               string `env:"EXTERNAL_POSTGRES_PORT,required"`
	User               string `env:"POSTGRES_USER,required"`
	Password           string `env:"POSTGRES_PASSWORD,required"`
	DB                 string `env:"POSTGRES_DB,required"`
	SSLMode            string `env:"POSTGRES_SSL_MODE,required"`
	MigrationDirectory string `env:"MIGRATION_DIRECTORY,required"`
}

type postgresConfig struct {
	raw postgresEnvConfig
}

func NewPostgresConfig() (*postgresConfig, error) {
	var raw postgresEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &postgresConfig{raw: raw}, nil
}

func (cfg *postgresConfig) Address() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.raw.Host, cfg.raw.Port, cfg.raw.User, cfg.raw.Password, cfg.raw.DB, cfg.raw.SSLMode)
}

func (cfg *postgresConfig) MigrationDirectory() string {
	return cfg.raw.MigrationDirectory
}
