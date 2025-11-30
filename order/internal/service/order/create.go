package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/pptkna/rocket-factory/order/internal/model"
)

func (s *service) Create(ctx context.Context, req model.CreateOrderRequest) (model.OrderDto, error) {
	parts, err := s.inventoryClient.ListParts(ctx, model.PartFiters{
		Uuids: req.PartUuids,
	})
	if err != nil {
		return model.OrderDto{}, err
	}
	if len(parts) != len(req.PartUuids) {
		return model.OrderDto{}, model.ErrNotFound
	}

	orderUuid := uuid.New()

	var totalPrice float32
	for _, p := range parts {
		totalPrice += p.Price
	}

	order := model.OrderDto{
		OrderUUID:  orderUuid.String(),
		UserUUID:   req.UserUUID,
		PartUuids:  req.PartUuids,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPendingPayment,
	}

	err = s.orderRepository.Create(ctx, order)
	if err != nil {
		if errors.Is(err, model.ErrConflict) {
			return model.OrderDto{}, model.ErrConflict
		}
		return model.OrderDto{}, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}
