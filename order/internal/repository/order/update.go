package order

import (
	"context"
	"fmt"

	"github.com/pptkna/rocket-factory/order/internal/model"
	repoConverter "github.com/pptkna/rocket-factory/order/internal/repository/converter"
)

// TODO: Добавить транзакции
func (r *repository) Update(ctx context.Context, orderDto *model.OrderDto) error {
	orderDtoRepoModel := repoConverter.OrderDtoToRepoModel(orderDto)

	query := `
		UPDATE orders
		SET user_uuid = $1, part_uuids = $2, total_price = $3, transaction_uuid = $4, payment_method = $5, status = $6, updated_at = $8
		WHERE order_uuid = $9
	`
	result, err := r.db.ExecContext(ctx, query,
		&orderDtoRepoModel.UserUUID,
		&orderDtoRepoModel.PartUuids,
		&orderDtoRepoModel.TotalPrice,
		&orderDtoRepoModel.TransactionUUID,
		&orderDtoRepoModel.PaymentMethod,
		&orderDtoRepoModel.Status,
		&orderDtoRepoModel.CreatedAt,
		&orderDtoRepoModel.UpdatedAt,
		&orderDtoRepoModel.OrderUUID,
	)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return model.ErrNotFound
	}

	return nil
}
