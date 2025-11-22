package service

import (
	"context"

	"github.com/pptkna/rocket-factory/order/internal/model"
)

type OrderService interface {
	PostOrder(ctx context.Context, orderDto model.CreateOrderRequest) (model.OrderDto, error)
	GetOrderByOrderUuid(ctx context.Context, uuid string) (model.OrderDto, error)
	PostOrderPay(ctx context.Context, orderUuid string, paymentMethod model.PaymentMethod) (string, error)
	PostOrderCancel(ctx context.Context, uuid string) error
}
