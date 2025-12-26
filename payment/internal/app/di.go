package app

import (
	"context"

	v1 "github.com/pptkna/rocket-factory/payment/internal/api/payment/v1"
	"github.com/pptkna/rocket-factory/payment/internal/service"
	"github.com/pptkna/rocket-factory/payment/internal/service/payment"
	paymentV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1Api paymentV1.PaymentServiceServer

	paymentService service.PaymentService
}

func newDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1Api(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentV1Api == nil {
		d.paymentV1Api = v1.NewAPI(d.PaymentService(ctx))
	}

	return d.paymentV1Api
}

func (d *diContainer) PaymentService(_ context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = payment.NewService()
	}

	return d.paymentService
}
