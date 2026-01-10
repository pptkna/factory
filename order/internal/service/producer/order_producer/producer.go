package order_producer

import (
	"context"

	kafkaEventConverter "github.com/pptkna/rocket-factory/order/internal/converter"
	"github.com/pptkna/rocket-factory/order/internal/model"
	def "github.com/pptkna/rocket-factory/order/internal/service"
	"github.com/pptkna/rocket-factory/platform/pkg/kafka"
	"github.com/pptkna/rocket-factory/platform/pkg/logger"
	orderEventsV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/events/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var _ def.OrderPaidProducerService = (*service)(nil)

type service struct {
	orderPaidProducer kafka.Producer
}

func NewService(orderPaidProducer kafka.Producer) *service {
	return &service{
		orderPaidProducer: orderPaidProducer,
	}
}

func (s *service) ProduceOrderPaid(ctx context.Context, event *model.OrderPaidEvent) error {
	msg := &orderEventsV1.OrderPaid{
		EventUuid:     event.EventUUID,
		OrderUuid:     event.OrderUUID,
		UserUuid:      event.UserUUID,
		PaymentMethod: kafkaEventConverter.PaymentMethodToEvent(event.PaymentMethod),
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal OrderPaid", zap.Error(err))
		return err
	}

	err = s.orderPaidProducer.Send(ctx, []byte(event.EventUUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish OrderPaid", zap.Error(err))
		return err
	}

	return nil
}
