package part

import (
	"context"
	"errors"
	"fmt"

	"github.com/pptkna/rocket-factory/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (*model.Part, error) {
	part, err := s.partRepository.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get part: %w", err)
	}

	return part, nil
}
