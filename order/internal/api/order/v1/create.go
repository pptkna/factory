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

func (a *api) PostOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.PostOrderRes, error) {
	order, err := a.service.Create(ctx, converter.CreateOrderRequestToModel(req))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintln("Not found, %w", err),
			}, nil
		}
		if errors.Is(err, model.ErrConflict) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintln("Not found, %w", err),
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	parsedOrderUuid, err := uuid.Parse(order.OrderUUID)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprint("Internal server error: %w", err),
		}, nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  parsedOrderUuid,
		TotalPrice: order.TotalPrice,
	}, nil
}
