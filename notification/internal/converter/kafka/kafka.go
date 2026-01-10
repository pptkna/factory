package kafka

import "github.com/pptkna/rocket-factory/notification/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (*model.OrderPaid, error)
}

type OrderAssembledDecoder interface {
	Decode(data []byte) (*model.OrderAssembled, error)
}
