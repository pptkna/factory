package order_consumer

import (
	"context"
	"math/rand"
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

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
		zap.String("order_uuid", orderPaidEvent.OrderUUID),
		zap.String("event_uuid", orderPaidEvent.EventUUID),
		zap.String("payment_method", string(orderPaidEvent.PaymentMethod)),
		zap.String("transaction_uuid", orderPaidEvent.TransactionUUID),
	)

	//nolint:gosec
	delay := time.Duration(rand.Intn(10)+1) * time.Second
	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return ctx.Err()
	}

	orderAssembledEvent := &model.OrderAssembledEvent{
		EventUUID:    orderPaidEvent.EventUUID,
		OrderUUID:    orderPaidEvent.OrderUUID,
		UserUUID:     orderPaidEvent.UserUUID,
		BuildTimeSec: int64(delay / time.Second),
	}

	err = s.orderAssembledProducer.ProduceOrderAssembled(ctx, orderAssembledEvent)
	if err != nil {
		logger.Error(ctx, "Failed to produce ship assembled event",
			zap.String("event_uuid", orderAssembledEvent.EventUUID),
			zap.String("order_uuid", orderAssembledEvent.OrderUUID),
			zap.Int64("build_time_sec", orderAssembledEvent.BuildTimeSec),
			zap.Error(err),
		)
		return err
	}

	return nil
}
