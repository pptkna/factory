package repository

import (
	"context"

	"github.com/pptkna/rocket-factory/inventory/internal/model"
)

type PartRepository interface {
	Get(ctx context.Context, uuid string) (model.Part, error)
	GetList(ctx context.Context, filters model.PartFiters) ([]model.Part, error)
}
