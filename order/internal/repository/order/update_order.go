package order

import (
	"context"

	"github.com/pptkna/rocket-factory/order/internal/model"
	repoConverter "github.com/pptkna/rocket-factory/order/internal/repository/converter"
)

func (r *repository) UpdateOrder(_ context.Context, orderDto model.OrderDto) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.orders[orderDto.OrderUUID]
	if !exists {
		return model.ErrNotFound
	}

	r.orders[orderDto.OrderUUID] = repoConverter.OrderDtoToRepoModel(orderDto)

	return nil
}
