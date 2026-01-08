package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/pptkna/rocket-factory/assembly/internal/config"
	kafkaConverter "github.com/pptkna/rocket-factory/assembly/internal/converter/kafka"
	decoder "github.com/pptkna/rocket-factory/assembly/internal/converter/kafka/decoder"
	"github.com/pptkna/rocket-factory/assembly/internal/service"
	orderConsumer "github.com/pptkna/rocket-factory/assembly/internal/service/consumer/order_consumer"
	orderProducer "github.com/pptkna/rocket-factory/assembly/internal/service/producer/order_producer"
	"github.com/pptkna/rocket-factory/platform/pkg/closer"
	wrappedKafka "github.com/pptkna/rocket-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/pptkna/rocket-factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/pptkna/rocket-factory/platform/pkg/kafka/producer"
	"github.com/pptkna/rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/pptkna/rocket-factory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	orderPaidConsumerService      service.OrderPaidConsumerService
	orderAssembledProducerService service.OrderAssembledProducerService

	consumerGroup sarama.ConsumerGroup

	syncProducer           sarama.SyncProducer
	orderAssembledProducer wrappedKafka.Producer

	orderPaidDecoder  kafkaConverter.OrderPaidDecoder
	orderPaidConsumer wrappedKafka.Consumer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderPaidConsumerService() service.OrderPaidConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = orderConsumer.NewService(d.OrderPaidConsumer(), d.OrderPaidDecoder(), d.OrderAssembledProducerService())
	}

	return d.orderPaidConsumerService
}

func (d *diContainer) OrderAssembledProducerService() service.OrderAssembledProducerService {
	if d.orderAssembledProducerService == nil {
		d.orderAssembledProducerService = orderProducer.NewService(d.OrderAssembledProducer())
	}

	return d.orderAssembledProducerService
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		c, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.ConsumerGroupID(),
			// TODO заменить OrderPaidConsumer на consumer
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = c
	}

	return d.consumerGroup
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			// TODO заменить OrderAssembledProducer на producer
			config.AppConfig().OrderAssembledProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) OrderAssembledProducer() wrappedKafka.Producer {
	if d.orderAssembledProducer == nil {
		d.orderAssembledProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderAssembledProducer.TopicName(),
			logger.Logger(),
		)
	}

	return d.orderAssembledProducer
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.TopicName(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderPaidConsumer
}
