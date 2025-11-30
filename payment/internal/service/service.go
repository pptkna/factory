package service

import (
	"context"

	"github.com/pptkna/rocket-factory/payment/internal/model"
)

type PaymentService interface {
	Pay(context context.Context, PayOrderRequest model.PayOrderRequest) model.PayOrderResponse
}
