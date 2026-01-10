package order_consumer

import (
	"context"

	"github.com/pptkna/rocket-factory/order/internal/config"
	kafkaConverter "github.com/pptkna/rocket-factory/order/internal/converter/kafka"
	def "github.com/pptkna/rocket-factory/order/internal/service"
	"github.com/pptkna/rocket-factory/platform/pkg/kafka"
	"github.com/pptkna/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

var _ def.OrderAssembledConsumerService = (*service)(nil)

type service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
	orderService           def.OrderService
}

func NewService(
	orderAssembledConsumer kafka.Consumer,
	orderAssembledDecoder kafkaConverter.OrderAssembledDecoder,
	orderService def.OrderService,
) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
		orderService:           orderService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order orderAssembledConsumer service")

	err := s.orderAssembledConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume error", zap.Error(err), zap.String("topic name", config.AppConfig().OrderAssembledConsumer.TopicName()))
		return err
	}

	return nil
}
