package converter

import (
	"github.com/pptkna/rocket-factory/order/internal/model"
	orderEventsV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/events/v1"
)

func PaymentMethodToEvent(paymentMethod model.PaymentMethod) orderEventsV1.PaymentMethod {
	switch paymentMethod {
	case model.PaymentMethodCard:
		return orderEventsV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return orderEventsV1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		return orderEventsV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return orderEventsV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	case model.PaymentMethodUnknown:
		return orderEventsV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN
	default:
		return orderEventsV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
