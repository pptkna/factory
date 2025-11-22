package v1

import "github.com/pptkna/rocket-factory/order/internal/service"

type api struct {
	service service.OrderService
}

func NewApi(service service.OrderService) *api {
	return &api{
		service: service,
	}
}
