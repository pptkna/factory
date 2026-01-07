package model

type PaymentMethod string

const (
	PaymentMethodUnspecified   PaymentMethod = "UNSPECIFIED"
	PaymentMethodUnknown       PaymentMethod = "UNKNOWN"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

type OrderPaidEvent struct {
	EventUUID       string
	OrderUUID       string
	UserUUID        string
	PaymentMethod   PaymentMethod
	TransactionUUID string
}

type OrderAssembledEvent struct {
	EventUUID    string
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
}
