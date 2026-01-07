package producer

import (
	"context"

	"github.com/pptkna/rocket-factory/assembly/internal/model"
	def "github.com/pptkna/rocket-factory/assembly/internal/service"
	"github.com/pptkna/rocket-factory/platform/pkg/kafka"
	"github.com/pptkna/rocket-factory/platform/pkg/logger"
	orderEventsV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/events/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var _ def.OrderAssembledProducerService = (*service)(nil)

type service struct {
	orderAssembledProducer kafka.Producer
}

func NewService(orderAssembledProducer kafka.Producer) *service {
	return &service{
		orderAssembledProducer: orderAssembledProducer,
	}
}

func (p *service) ProduceOrderAssembled(ctx context.Context, event model.OrderAssembledEvent) error {
	msg := &orderEventsV1.OrderAssembled{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal OrderAssembled", zap.Error(err))
		return err
	}

	err = p.orderAssembledProducer.Send(ctx, []byte(event.EventUUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish OrderAssembled", zap.Error(err))
		return err
	}

	return nil
}
