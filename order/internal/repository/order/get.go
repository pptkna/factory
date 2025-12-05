package order

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pptkna/rocket-factory/order/internal/model"
	"github.com/samber/lo"
)

func (r *repository) Get(ctx context.Context, uuid string) (*model.OrderDto, error) {
	var orderUUID string
	var userUUID string
	var partUuids []string
	var totalPrice float32
	var transactionUUID *string
	var paymentMethod sql.NullString
	var status string
	var createdAt time.Time

	orderQuery := `
		SELECT order_uuid, user_uuid, part_uuids, total_price, transaction_uuid, payment_method, status, created_at
		FROM orders
		WHERE order_uuid = $1
	`
	err := r.db.QueryRowContext(ctx, orderQuery, uuid).Scan(&orderUUID, &userUUID, &partUuids, &totalPrice, &transactionUUID, &paymentMethod, &status, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	var paymentMethodModel *model.PaymentMethod
	if paymentMethod.Valid {
		paymentMethodModel = lo.ToPtr(model.PaymentMethod(paymentMethod.String))
	}

	return &model.OrderDto{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUuids,
		TotalPrice:      totalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethodModel,
		Status:          model.OrderStatus(status),
		CreatedAt:       createdAt,
	}, nil
}
