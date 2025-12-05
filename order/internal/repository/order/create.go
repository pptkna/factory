package order

import (
	"context"
	"fmt"

	"github.com/pptkna/rocket-factory/order/internal/model"
	repoConverter "github.com/pptkna/rocket-factory/order/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, orderDto *model.OrderDto) error {
	repoOrderDto := repoConverter.OrderDtoToRepoModel(orderDto)
	query := `
		INSERT INTO orders (order_uuid, user_uuid, part_uuids, total_price, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		repoOrderDto.OrderUUID,
		repoOrderDto.UserUUID,
		repoOrderDto.PartUuids,
		repoOrderDto.TotalPrice,
		repoOrderDto.Status,
		repoOrderDto.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}
