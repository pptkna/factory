package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/pptkna/rocket-factory/order/internal/model"
)

func (s *service) Pay(ctx context.Context, orderUuid string, paymentMethod model.PaymentMethod) (string, error) {
	order, err := s.orderRepository.Get(ctx, orderUuid)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "", model.ErrNotFound
		}
		return "", fmt.Errorf("failed to get order with OrderUUID: %s, %w", orderUuid, err)
	}
	if order.Status == model.OrderStatusPaid || order.Status == model.OrderStatusCancelled {
		return "", model.ErrConflict
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, orderUuid, order.UserUUID, paymentMethod)
	if err != nil {
		return "", err
	}

	newOrder := model.OrderDto{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: &transactionUUID,
		PaymentMethod:   &paymentMethod,
		Status:          model.OrderStatusPaid,
	}

	err = s.orderRepository.Update(ctx, newOrder)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "", model.ErrNotFound
		}
		return "", fmt.Errorf("failed to cancel order with OrderUUID: %s, %w", orderUuid, err)
	}

	return transactionUUID, nil
}
