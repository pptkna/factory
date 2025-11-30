package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/pptkna/rocket-factory/order/internal/converter"
	"github.com/pptkna/rocket-factory/order/internal/model"
	orderV1 "github.com/pptkna/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByOrderUuid(ctx context.Context, params orderV1.GetOrderByOrderUuidParams) (orderV1.GetOrderByOrderUuidRes, error) {
	order, err := a.service.Get(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order with uuid: %s not found", params.OrderUUID.String()),
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	parsedOrder, err := converter.OrderDtoToOrderV1(order)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	return parsedOrder, nil
}
