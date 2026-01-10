package order_consumer

import (
	"context"

	"github.com/pptkna/rocket-factory/order/internal/model"
	"github.com/pptkna/rocket-factory/platform/pkg/kafka"
	"github.com/pptkna/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	orderAssembledEvent, err := s.orderAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode orderAssembled event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
		zap.String("order_uuid", orderAssembledEvent.OrderUUID),
		zap.String("event_uuid", orderAssembledEvent.EventUUID),
		zap.Int64("build_time_sec", orderAssembledEvent.BuildTimeSec),
	)

	err = s.orderService.UpdateStatus(ctx, orderAssembledEvent.OrderUUID, model.OrderStatusAssembled)
	if err != nil {
		logger.Error(ctx, "Failed to update order status to completed",
			zap.String("order_uuid", orderAssembledEvent.OrderUUID),
			zap.Error(err),
		)
		return err
	}

	return nil
}
