package part

import (
	"github.com/pptkna/rocket-factory/inventory/internal/repository"
	def "github.com/pptkna/rocket-factory/inventory/internal/service"
)

var _ def.PartService = (*service)(nil)

type service struct {
	partRepository repository.PartRepository
}

func NewService(repository repository.PartRepository) *service {
	return &service{
		partRepository: repository,
	}
}
