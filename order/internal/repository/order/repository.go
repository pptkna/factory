package order

import (
	"sync"

	def "github.com/pptkna/rocket-factory/order/internal/repository"
	"github.com/pptkna/rocket-factory/order/internal/repository/model"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu     sync.RWMutex
	orders map[string]*model.OrderDto
}

func NewRepository() *repository {
	return &repository{
		orders: make(map[string]*model.OrderDto),
	}
}
