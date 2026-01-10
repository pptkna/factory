package config

import (
	"time"

	"github.com/IBM/sarama"
)

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

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidProducerConfig interface {
	TopicName() string
	Config() *sarama.Config
}

type OrderAssembledConsumerConfig interface {
	TopicName() string
	ConsumerGroupID() string
	Config() *sarama.Config
}
