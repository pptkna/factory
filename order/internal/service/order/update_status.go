package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/pptkna/rocket-factory/order/internal/model"
	"github.com/pptkna/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) UpdateStatus(ctx context.Context, orderUuid string, status model.OrderStatus) error {
	order, err := s.orderRepository.Get(ctx, orderUuid)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return model.ErrNotFound
		}
		return fmt.Errorf("failed to get order with OrderUUID: %s, %w", orderUuid, err)
	}

	newOrder := &model.OrderDto{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUuids,
		TotalPrice: order.TotalPrice,
		Status:     status,
	}

	err = s.orderRepository.Update(ctx, newOrder)
	if err != nil {
		logger.Error(ctx, "update status error",
			zap.String("orderUUID", orderUuid),
			zap.String("status", string(status)),
			zap.Error(err),
		)
		if errors.Is(err, model.ErrNotFound) {
			return model.ErrNotFound
		}
		return fmt.Errorf("failed to cancel order with OrderUUID: %s, %w", orderUuid, err)
	}

	return nil
}
