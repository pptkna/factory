package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type ApiConfig interface {
	Address() string
	ReadTimeout() time.Duration
	ShutDownTimeout() time.Duration
}

type InventoryGRPCConfig interface {
	Address() string
}

type PaymentGRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	Address() string
	MigrationDirectory() string
}
