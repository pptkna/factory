package part

import (
	"context"

	"github.com/pptkna/rocket-factory/inventory/internal/model"
	repoConverter "github.com/pptkna/rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoPart, ok := r.parts[uuid]
	if !ok {
		return model.Part{}, model.ErrNotFound
	}

	return repoConverter.PartToModel(repoPart), nil
}
