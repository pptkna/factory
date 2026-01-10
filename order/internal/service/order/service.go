package order

import (
	client "github.com/pptkna/rocket-factory/order/internal/client/grpc"
	"github.com/pptkna/rocket-factory/order/internal/repository"
	def "github.com/pptkna/rocket-factory/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository repository.OrderRepository

	inventoryClient client.InventoryClient
	paymentClient   client.PaymentClient

	orderPaidProducerService def.OrderPaidProducerService
}

func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient client.InventoryClient,
	paymentClient client.PaymentClient,
	orderPaidProducerService def.OrderPaidProducerService,
) *service {
	return &service{
		orderRepository: orderRepository,

		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,

		orderPaidProducerService: orderPaidProducerService,
	}
}
