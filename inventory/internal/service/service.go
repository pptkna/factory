package service

import (
	"context"

	"github.com/pptkna/rocket-factory/inventory/internal/model"
)

type PartService interface {
	Get(context context.Context, uuid string) (*model.Part, error)
	GetList(context context.Context, filters *model.PartFiters) ([]*model.Part, error)
}
