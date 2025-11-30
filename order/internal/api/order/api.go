package order

import (
	"context"

	orderV1 "github.com/pptkna/rocket-factory/shared/pkg/openapi/order/v1"
)

type OrderApi interface {
	Create(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.PostOrderRes, error)
	Get(ctx context.Context, params orderV1.GetOrderByOrderUuidParams) (orderV1.GetOrderByOrderUuidRes, error)
	Pay(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PostOrderPayParams) (orderV1.PostOrderPayRes, error)
	Cancel(ctx context.Context, params orderV1.PostOrderCancelParams) (orderV1.PostOrderCancelRes, error)
	NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode
}
