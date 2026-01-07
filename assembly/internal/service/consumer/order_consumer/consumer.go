package order_consumer

import (
	"context"

	"github.com/pptkna/rocket-factory/assembly/internal/config"
	kafkaConverter "github.com/pptkna/rocket-factory/assembly/internal/converter/kafka"
	def "github.com/pptkna/rocket-factory/assembly/internal/service"
	"github.com/pptkna/rocket-factory/platform/pkg/kafka"
	"github.com/pptkna/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

var _ def.OrderPaidConsumerService = (*service)(nil)

type service struct {
	orderPaidConsumer kafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder

	orderAssembledProducer def.OrderAssembledProducerService
}

func NewService(orderPaidConsumer kafka.Consumer, orderPaidDecoder kafkaConverter.OrderPaidDecoder, orderAssembledProducer def.OrderAssembledProducerService) *service {
	return &service{
		orderPaidConsumer:      orderPaidConsumer,
		orderPaidDecoder:       orderPaidDecoder,
		orderAssembledProducer: orderAssembledProducer,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order orderPaidConsumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume error", zap.Error(err), zap.String("topic name", config.AppConfig().OrderPaidConsumer.TopicName()))
		return err
	}

	return nil
}
