package converter

import (
	"github.com/pptkna/rocket-factory/payment/internal/model"
	paymentV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/payment/v1"
)

func PayOrderRequestToModel(req *paymentV1.PayOrderRequest) model.PayOrderRequest {
	return model.PayOrderRequest{
		OrderUuid:     req.GetOrderUuid(),
		UserUuid:      req.GetUserUuid(),
		PaymentMethod: PaymentMethodToModel(req.GetPaymentMethod()),
	}
}

func PaymentMethodToModel(paymentMethod paymentV1.PaymentMethod) model.PaymentMethod {
	switch paymentMethod {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN:
		return model.PaymentMethodUnknown
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentV1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}

func PayOrderResponseToProto(res model.PayOrderResponse) *paymentV1.PayOrderResponse {
	return &paymentV1.PayOrderResponse{
		TransactionUuid: res.TransactionUuid,
	}
}
