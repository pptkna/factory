package order_consumer

import (
	"context"
	"time"

	"github.com/pptkna/rocket-factory/assembly/internal/model"
	"github.com/pptkna/rocket-factory/platform/pkg/kafka"
	"github.com/pptkna/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	orderPaidEvent, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid event", zap.Error(err))
		return err
	}

	start := time.Now()

	time.Sleep(10 * time.Second)

	buildTimeSec := int64(time.Since(start).Seconds())

	orderAssembledEvent := &model.OrderAssembledEvent{
		EventUUID:    orderPaidEvent.EventUUID,
		OrderUUID:    orderPaidEvent.OrderUUID,
		UserUUID:     orderPaidEvent.UserUUID,
		BuildTimeSec: buildTimeSec,
	}

	s.orderAssembledProducer.ProduceOrderAssembled(ctx, orderAssembledEvent)

	return nil
}
