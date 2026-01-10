package order_assembled_consumer

import (
	"context"
	"math/rand/v2"
	"time"

	"go.uber.org/zap"

	"github.com/pptkna/rocket-factory/notification/internal/model"
	"github.com/pptkna/rocket-factory/platform/pkg/kafka"
	"github.com/pptkna/rocket-factory/platform/pkg/logger"
)

func (s *service) OrderAssembledHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode order paid event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("event_uuid", event.EventUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	//nolint:gosec
	delay := time.Duration(rand.IntN(10)+1) * time.Second
	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return ctx.Err()
	}

	orderAssembled := &model.OrderAssembled{
		EventUUID:    event.EventUUID,
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: int64(delay / time.Second), // Здесь можно использовать реальное время сборки
	}

	err = s.telegramService.SendAssembledNotification(ctx, orderAssembled)
	if err != nil {
		logger.Error(ctx, "Failed to produce ship assembled event to Telegram",
			zap.String("event_uuid", orderAssembled.EventUUID),
			zap.String("order_uuid", orderAssembled.OrderUUID),
			zap.Int64("build_time_sec", orderAssembled.BuildTimeSec),
			zap.Error(err),
		)
		return err
	}

	return nil
}
