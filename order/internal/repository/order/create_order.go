package order

import (
	"context"

	"github.com/pptkna/rocket-factory/order/internal/model"
	repoConverter "github.com/pptkna/rocket-factory/order/internal/repository/converter"
)

func (r *repository) CreateOrder(_ context.Context, orderDto model.OrderDto) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exist := r.orders[orderDto.OrderUUID]; exist {
		return model.ErrConflict
	}

	r.orders[orderDto.OrderUUID] = repoConverter.OrderDtoToRepoModel(orderDto)

	return nil
}
