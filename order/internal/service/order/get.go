package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/pptkna/rocket-factory/order/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (model.OrderDto, error) {
	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return model.OrderDto{}, model.ErrNotFound
		}
		return model.OrderDto{}, fmt.Errorf("failed to get order with OrderUUID: %s, %w", uuid, err)
	}

	return order, nil
}
