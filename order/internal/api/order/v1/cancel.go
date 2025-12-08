package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/pptkna/rocket-factory/order/internal/model"
	orderV1 "github.com/pptkna/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) PostOrderCancel(ctx context.Context, params orderV1.PostOrderCancelParams) (orderV1.PostOrderCancelRes, error) {
	err := a.service.Cancel(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return &orderV1.ConflictError{
				Code:    404,
				Message: fmt.Sprintf("Order with UUID: %s not found", params.OrderUUID.String()),
			}, nil
		}
		if errors.Is(err, model.ErrConflict) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("Order with UUID: %s already payed or closed", params.OrderUUID.String()),
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("Internal server error, %v", err),
		}, nil
	}

	return &orderV1.PostOrderCancelNoContent{}, nil
}
