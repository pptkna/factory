package service

import (
	"context"

	"github.com/pptkna/rocket-factory/inventory/internal/model"
)

type PartService interface {
	GetPart(context context.Context, uuid string) (model.Part, error)
	ListParts(context context.Context, filters model.PartFiters) ([]model.Part, error)
}
