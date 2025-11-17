package v1

import (
	"context"

	"github.com/pptkna/rocket-factory/payment/internal/converter"
	paymentV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	res := a.paymentService.PayOrder(ctx, converter.PayOrderRequestToModel(req))

	return converter.PayOrderResponseToProto(res), nil
}
