package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/pptkna/rocket-factory/order/internal/converter"
	"github.com/pptkna/rocket-factory/order/internal/model"
	orderV1 "github.com/pptkna/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) PostOrderPay(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PostOrderPayParams) (orderV1.PostOrderPayRes, error) {
	transactionUUID, err := a.service.Pay(ctx, params.OrderUUID.String(), converter.PaymentMethodToModel(req.PaymentMethod))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order with uuid: %s not found", params.OrderUUID),
			}, nil
		}
		if errors.Is(err, model.ErrConflict) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("Order with UUID: %s already payed or cancelled", params.OrderUUID.String()),
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	parsedTtransactionUUID, err := uuid.Parse(transactionUUID)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprint("Internal server error: %w", err),
		}, nil
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: parsedTtransactionUUID,
	}, nil
}
