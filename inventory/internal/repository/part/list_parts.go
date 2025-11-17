package part

import (
	"context"
	"slices"

	"github.com/pptkna/rocket-factory/inventory/internal/model"
	repoConverter "github.com/pptkna/rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) ListParts(_ context.Context, filters model.PartFiters) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.parts) == 0 {
		return []model.Part{}, model.ErrNotFound
	}

	repoFilters := repoConverter.PartFitersToRepoModel(filters)

	listParts := make([]model.Part, 0, len(r.parts))

	for _, p := range r.parts {
		if len(repoFilters.Uuids) > 0 && !slices.Contains(repoFilters.Uuids, p.Uuid) {
			continue
		}

		if len(repoFilters.Categories) > 0 && !slices.Contains(repoFilters.Categories, p.Category) {
			continue
		}

		if len(repoFilters.ManufacturerCountries) > 0 && !slices.Contains(repoFilters.ManufacturerCountries, p.Manufacturer.Country) {
			continue
		}

		if len(repoFilters.Names) > 0 && !slices.Contains(repoFilters.Names, p.Name) {
			continue
		}

		if len(repoFilters.Tags) > 0 {
			found := false

			for _, t := range repoFilters.Tags {
				if slices.Contains(p.Tags, t) {
					found = true
					break
				}
			}

			if !found {
				continue
			}
		}

		listParts = append(listParts, repoConverter.PartToModel(p))
	}

	if len(listParts) == 0 {
		return []model.Part{}, model.ErrNotFound
	}

	return listParts, nil
}
