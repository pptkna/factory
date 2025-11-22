package payment

import (
	"context"

	clientConverter "github.com/pptkna/rocket-factory/order/internal/client/converter"
	"github.com/pptkna/rocket-factory/order/internal/model"
	generatedPaymentClientV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUuid, userUuid string, paymentMethod model.PaymentMethod) (string, error) {
	res, err := c.generatedClient.PayOrder(ctx, &generatedPaymentClientV1.PayOrderRequest{
		OrderUuid:     orderUuid,
		UserUuid:      userUuid,
		PaymentMethod: clientConverter.PaymentMethodToProto(paymentMethod),
	})
	if err != nil {
		return "", err
	}

	return res.GetTransactionUuid(), nil
}
