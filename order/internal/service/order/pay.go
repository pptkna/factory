package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/pptkna/rocket-factory/order/internal/model"
	"github.com/pptkna/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

// TODO: Добавить транзакции
func (s *service) Pay(ctx context.Context, orderUuid string, paymentMethod model.PaymentMethod) (string, error) {
	order, err := s.orderRepository.Get(ctx, orderUuid)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "", model.ErrNotFound
		}
		return "", fmt.Errorf("failed to get order with OrderUUID: %s, %w", orderUuid, err)
	}
	if order.Status == model.OrderStatusPaid || order.Status == model.OrderStatusCancelled || order.Status == model.OrderStatusAssembled {
		return "", model.ErrConflict
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, orderUuid, order.UserUUID, paymentMethod)
	if err != nil {
		return "", err
	}

	newOrder := &model.OrderDto{
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
		return "", fmt.Errorf("failed to update order with OrderUUID: %s, %w", orderUuid, err)
	}

	err = s.orderPaidProducerService.ProduceOrderPaid(ctx, &model.OrderPaidEvent{
		EventUUID:       uuid.NewString(),
		OrderUUID:       orderUuid,
		UserUUID:        newOrder.UserUUID,
		PaymentMethod:   paymentMethod,
		TransactionUUID: transactionUUID,
	})
	if err != nil {
		logger.Error(ctx, "faild to produce OrderPaidEvent", zap.Error(err))
		return "", err
	}

	return transactionUUID, nil
}
