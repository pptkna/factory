package service

import (
	"context"

	"github.com/pptkna/rocket-factory/payment/internal/model"
)

type PaymentService interface {
	PayOrder(context context.Context, PayOrderRequest model.PayOrderRequest) model.PayOrderResponse
}
