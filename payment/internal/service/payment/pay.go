package payment

import (
	"context"

	"github.com/google/uuid"
	"github.com/pptkna/rocket-factory/payment/internal/model"
	"github.com/pptkna/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) Pay(ctx context.Context, req model.PayOrderRequest) model.PayOrderResponse {
	transactionUuid := uuid.New().String()

	logger.Info(ctx, "Payment success", zap.String("transaction_uuid", transactionUuid))

	return model.PayOrderResponse{
		TransactionUuid: transactionUuid,
	}
}
