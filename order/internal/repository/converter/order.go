package converter

import (
	"github.com/pptkna/rocket-factory/order/internal/model"
	repoModel "github.com/pptkna/rocket-factory/order/internal/repository/model"
	"github.com/samber/lo"
)

func OrderDtoToRepoModel(orderDto model.OrderDto) *repoModel.OrderDto {
	var transactionUUID *string
	if orderDto.TransactionUUID != nil {
		transactionUUID = lo.ToPtr(*orderDto.TransactionUUID)
	}

	var paymentMethod *repoModel.PaymentMethod
	if orderDto.PaymentMethod != nil {
		paymentMethod = lo.ToPtr(PaymentMethodToRepoModel(*orderDto.PaymentMethod))
	}

	return &repoModel.OrderDto{
		OrderUUID:       orderDto.OrderUUID,
		UserUUID:        orderDto.UserUUID,
		PartUuids:       orderDto.PartUuids,
		TotalPrice:      orderDto.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          OrderStatusToRepoModel(orderDto.Status),
	}
}

func PaymentMethodToRepoModel(paymentMethod model.PaymentMethod) repoModel.PaymentMethod {
	switch paymentMethod {
	case model.PaymentMethodCard:
		return repoModel.PaymentMethodCard
	case model.PaymentMethodSBP:
		return repoModel.PaymentMethodSBP
	case model.PaymentMethodCreditCard:
		return repoModel.PaymentMethodCreditCard
	case model.PaymentMethodInvestorMoney:
		return repoModel.PaymentMethodInvestorMoney
	default:
		return repoModel.PaymentMethodUnknown
	}
}

func OrderStatusToRepoModel(orderStatus model.OrderStatus) repoModel.OrderStatus {
	switch orderStatus {
	case model.OrderStatusPendingPayment:
		return repoModel.OrderStatusPendingPayment
	case model.OrderStatusPaid:
		return repoModel.OrderStatusPaid
	default:
		return repoModel.OrderStatusCancelled
	}
}

func OrderDtoToModel(orderDto *repoModel.OrderDto) model.OrderDto {
	var transactionUUID *string
	if orderDto.TransactionUUID != nil {
		transactionUUID = lo.ToPtr(*orderDto.TransactionUUID)
	}

	var paymentMethod *model.PaymentMethod
	if orderDto.PaymentMethod != nil {
		paymentMethod = lo.ToPtr(PaymentMethodToModel(*orderDto.PaymentMethod))
	}

	return model.OrderDto{
		OrderUUID:       orderDto.OrderUUID,
		UserUUID:        orderDto.UserUUID,
		PartUuids:       orderDto.PartUuids,
		TotalPrice:      orderDto.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          OrderStatusToModel(orderDto.Status),
	}
}

func PaymentMethodToModel(paymentMethod repoModel.PaymentMethod) model.PaymentMethod {
	switch paymentMethod {
	case repoModel.PaymentMethodCard:
		return model.PaymentMethodCard
	case repoModel.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case repoModel.PaymentMethodCreditCard:
		return model.PaymentMethodCreditCard
	case repoModel.PaymentMethodInvestorMoney:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnknown
	}
}

func OrderStatusToModel(orderStatus repoModel.OrderStatus) model.OrderStatus {
	switch orderStatus {
	case repoModel.OrderStatusPendingPayment:
		return model.OrderStatusPendingPayment
	case repoModel.OrderStatusPaid:
		return model.OrderStatusPaid
	default:
		return model.OrderStatusCancelled
	}
}
