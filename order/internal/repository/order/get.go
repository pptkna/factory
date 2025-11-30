package order

import (
	"context"
	"fmt"

	"github.com/pptkna/rocket-factory/order/internal/model"
	repoConverter "github.com/pptkna/rocket-factory/order/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, uuid string) (model.OrderDto, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orderDto, exists := r.orders[uuid]
	if !exists {
		return model.OrderDto{}, model.ErrNotFound
	}
	if orderDto == nil {
		return model.OrderDto{}, fmt.Errorf("failed to get order with OrderUUID: %s", uuid)
	}

	return repoConverter.OrderDtoToModel(orderDto), nil
}
