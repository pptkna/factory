package kafka

import "github.com/pptkna/rocket-factory/order/internal/model"

type OrderAssembledDecoder interface {
	Decode(data []byte) (*model.OrderAssembledEvent, error)
}
