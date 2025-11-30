package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/pptkna/rocket-factory/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, uuid string) error {
	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return model.ErrNotFound
		}
		return fmt.Errorf("failed to get order with OrderUUID: %s, %w", uuid, err)
	}
	if order.Status == model.OrderStatusPaid || order.Status == model.OrderStatusCancelled {
		return model.ErrConflict
	}

	newOrder := model.OrderDto{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUuids,
		TotalPrice: order.TotalPrice,
		Status:     model.OrderStatusCancelled,
	}

	err = s.orderRepository.Update(ctx, newOrder)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return model.ErrNotFound
		}
		return fmt.Errorf("failed to cancel order with OrderUUID: %s, %w", uuid, err)
	}

	return nil
}
