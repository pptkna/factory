package part

import (
	"context"
	"errors"
	"fmt"

	"github.com/pptkna/rocket-factory/inventory/internal/model"
)

func (s *service) ListParts(context context.Context, filters model.PartFiters) ([]model.Part, error) {
	parts, err := s.partRepository.ListParts(context, filters)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return []model.Part{}, model.ErrNotFound
		}
		return []model.Part{}, fmt.Errorf("failed to get list part: %w", err)
	}

	return parts, nil
}
