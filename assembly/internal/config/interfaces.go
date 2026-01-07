package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type KafkaConfig interface {
	Brokers() string
}

type OrderAssembledProducerConfig interface {
	TopicName() string
}

type OrderPaidConsumerConfig interface {
	TopicName() string
	ConsumerGroupID() string
}
