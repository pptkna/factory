package service

import (
	"context"

	"github.com/pptkna/rocket-factory/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, orderDto *model.CreateOrderRequest) (*model.OrderDto, error)
	Get(ctx context.Context, uuid string) (*model.OrderDto, error)
	Pay(ctx context.Context, orderUuid string, paymentMethod model.PaymentMethod) (string, error)
	Cancel(ctx context.Context, uuid string) error
}
