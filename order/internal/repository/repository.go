package repository

import (
	"context"

	"github.com/pptkna/rocket-factory/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, orderDto model.OrderDto) error
	GetOrder(ctx context.Context, uuid string) (model.OrderDto, error)
	UpdateOrder(ctx context.Context, orderDto model.OrderDto) error
}
