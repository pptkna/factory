package payment

import (
	def "github.com/pptkna/rocket-factory/order/internal/client/grpc"
	generatedPaymentClientV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	generatedClient generatedPaymentClientV1.PaymentServiceClient
}

func NewClient(generatedClient generatedPaymentClientV1.PaymentServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
