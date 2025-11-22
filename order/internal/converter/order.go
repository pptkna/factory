package converter

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pptkna/rocket-factory/order/internal/model"
	orderV1 "github.com/pptkna/rocket-factory/shared/pkg/openapi/order/v1"
	"github.com/samber/lo"
)

func OrderDtoToModel(orderDto *orderV1.OrderDto) model.OrderDto {
	partUuids := make([]string, len(orderDto.GetPartUuids()))
	for i, u := range orderDto.GetPartUuids() {
		partUuids[i] = u.String()
	}

	var transactionUUID *string
	if ok := orderDto.GetTransactionUUID().Set; ok {
		transactionUUID = lo.ToPtr(orderDto.GetTransactionUUID().Value.String())
	}

	var paymentMethod *model.PaymentMethod
	if ok := orderDto.GetPaymentMethod().Set; ok {
		paymentMethod = lo.ToPtr(PaymentMethodToModel(orderDto.GetPaymentMethod().Value))
	}

	return model.OrderDto{
		OrderUUID:       orderDto.GetOrderUUID().String(),
		UserUUID:        orderDto.GetUserUUID().String(),
		PartUuids:       partUuids,
		TotalPrice:      orderDto.GetTotalPrice(),
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          OrderStatusToModel(orderDto.Status),
	}
}

func PaymentMethodToModel(paymentMethod orderV1.PaymentMethod) model.PaymentMethod {
	switch paymentMethod {
	case orderV1.PaymentMethodCARD:
		return model.PaymentMethodCard
	case orderV1.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case orderV1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCreditCard
	default:
		return model.PaymentMethodUnknown
	}
}

func OrderStatusToModel(status orderV1.OrderStatus) model.OrderStatus {
	switch status {
	case orderV1.OrderStatusPENDINGPAYMENT:
		return model.OrderStatusPendingPayment
	case orderV1.OrderStatusPAID:
		return model.OrderStatusPaid
	default:
		return model.OrderStatusCancelled
	}
}

func CreateOrderRequestToModel(req *orderV1.CreateOrderRequest) model.CreateOrderRequest {
	partUuids := make([]string, len(req.PartUuids))
	for i, u := range req.PartUuids {
		partUuids[i] = u.String()
	}

	return model.CreateOrderRequest{
		UserUUID:  req.UserUUID.String(),
		PartUuids: partUuids,
	}
}

func OrderDtoToOrderV1(orderDto model.OrderDto) (*orderV1.OrderDto, error) {
	orderUuid, err := uuid.Parse(orderDto.OrderUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse orderUuid")
	}

	userUuid, err := uuid.Parse(orderDto.UserUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse userUuid")
	}

	partUuids := make([]uuid.UUID, len(orderDto.PartUuids))
	for i, pu := range orderDto.PartUuids {
		parsedUuid, err := uuid.Parse(pu)
		if err != nil {
			return nil, fmt.Errorf("failed to parse partUuid")
		}

		partUuids[i] = parsedUuid
	}

	var transactionUUID orderV1.OptUUID
	if orderDto.TransactionUUID != nil {
		parsedUuid, err := uuid.Parse(*orderDto.TransactionUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse transactionUUID")
		}

		transactionUUID.SetTo(parsedUuid)
	}

	var paymentMethod orderV1.OptPaymentMethod
	if orderDto.PaymentMethod != nil {
		paymentMethod.SetTo(PaymentMethodToOrderV1(*orderDto.PaymentMethod))
	}

	return &orderV1.OrderDto{
		OrderUUID:       orderUuid,
		UserUUID:        userUuid,
		PartUuids:       partUuids,
		TotalPrice:      orderDto.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          OrderStatusToOrderV1(orderDto.Status),
	}, nil
}

func PaymentMethodToOrderV1(paymentMethod model.PaymentMethod) orderV1.PaymentMethod {
	switch paymentMethod {
	case model.PaymentMethodCard:
		return orderV1.PaymentMethodCARD
	case model.PaymentMethodSBP:
		return orderV1.PaymentMethodSBP
	case model.PaymentMethodCreditCard:
		return orderV1.PaymentMethodCREDITCARD
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}

func OrderStatusToOrderV1(status model.OrderStatus) orderV1.OrderStatus {
	switch status {
	case model.OrderStatusPendingPayment:
		return orderV1.OrderStatusPENDINGPAYMENT
	case model.OrderStatusPaid:
		return orderV1.OrderStatusPAID
	default:
		return orderV1.OrderStatusCANCELLED
	}
}
