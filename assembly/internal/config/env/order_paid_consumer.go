package env

import "github.com/caarlos0/env/v11"

type orderPaidConsumerEnvConfig struct {
	topicName       string `env:"ORDER_PAID_TOPIC_NAME,required"`
	consumerGroupID string `env:"ORDER_PAID_CONSUMER_GROUP_ID,required"`
}

type orderPaidConsumerConfig struct {
	raw orderPaidConsumerEnvConfig
}

func NewOrderPaidConsumerConfig() (*orderPaidConsumerConfig, error) {
	var raw orderPaidConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderPaidConsumerConfig{raw: raw}, nil
}

func (cfg *orderPaidConsumerConfig) TopicName() string {
	return cfg.raw.topicName
}

func (cfg *orderPaidConsumerConfig) ConsumerGroupID() string {
	return cfg.raw.consumerGroupID
}
