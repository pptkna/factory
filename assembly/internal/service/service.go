package service

import (
	"context"

	"github.com/pptkna/rocket-factory/assembly/internal/model"
)

type OrderPaidConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderAssembledProducerService interface {
	ProduceOrderAssembled(ctx context.Context, event *model.OrderAssembledEvent) error
}
