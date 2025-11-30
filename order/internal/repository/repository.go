package repository

import (
	"context"

	"github.com/pptkna/rocket-factory/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, orderDto model.OrderDto) error
	Get(ctx context.Context, uuid string) (model.OrderDto, error)
	Update(ctx context.Context, orderDto model.OrderDto) error
}
