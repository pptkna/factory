package kafka

import "github.com/pptkna/rocket-factory/assembly/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (*model.OrderPaidEvent, error)
}
