package payment

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/pptkna/rocket-factory/payment/internal/model"
)

func (s *service) PayOrder(ctx context.Context, req model.PayOrderRequest) model.PayOrderResponse {
	transactionUuid := uuid.New().String()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUuid)

	return model.PayOrderResponse{
		TransactionUuid: transactionUuid,
	}
}
