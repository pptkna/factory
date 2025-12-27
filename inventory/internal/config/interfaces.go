package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type InventoryGRPCConfig interface {
	Port() string
	Address() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}
