package model

import "time"

type PaymentMethod string

const (
	PaymentMethodUnknown       PaymentMethod = "UNKNOWN"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type OrderDto struct {
	OrderUUID       string
	UserUUID        string
	PartUuids       []string
	TotalPrice      float32
	TransactionUUID *string
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}

type PayParams struct {
	OrderUUID     string
	PaymentMethod PaymentMethod
}

type Category string

const (
	CategoryUnspecified Category = "UNSPECIFIED"
	CategoryUnknown     Category = "UNKNOWN"
	CategoryEngine      Category = "ENGINE"
	CategoryFuel        Category = "FUEL"
	CategoryPorthole    Category = "PORTHOLE"
	CategoryWing        Category = "WING"
)

type PartFiters struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

type Dimensions struct {
	Length float32
	Width  float32
	Height float32
	Weight float32
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Value struct {
	String *string
	Int64  *int64
	Double *float64
	Bool   *bool
}

type Part struct {
	Uuid          string
	Name          string
	Description   string
	Price         float32
	StockQuantity int
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]*Value
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

type CreateOrderRequest struct {
	UserUUID  string
	PartUuids []string
}
