package part

import (
	"context"
	"fmt"

	"github.com/pptkna/rocket-factory/inventory/internal/model"
	repoConverter "github.com/pptkna/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/pptkna/rocket-factory/inventory/internal/repository/model"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *repository) GetList(ctx context.Context, filters *model.PartFiters) ([]*model.Part, error) {
	mongoFilter := bson.M{}

	if filters != nil {
		if len(filters.Uuids) > 0 {
			mongoFilter["uuid"] = bson.M{"$in": filters.Uuids}
		}
		if len(filters.Categories) > 0 {
			mongoFilter["category"] = bson.M{"$in": filters.Categories}
		}
		if len(filters.ManufacturerCountries) > 0 {
			mongoFilter["manufacturer.country"] = bson.M{"$in": filters.ManufacturerCountries}
		}
		if len(filters.Names) > 0 {
			mongoFilter["name"] = bson.M{"$in": filters.Names}
		}
		if len(filters.Tags) > 0 {
			// $all проверяет, что массив содержит все указанные элементы
			mongoFilter["tags"] = bson.M{"$all": filters.Tags}
		}
	}

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch parts: %w", err)
	}
	defer cursor.Close(ctx)

	var parts []*repoModel.Part
	if err := cursor.All(ctx, &parts); err != nil {
		return nil, fmt.Errorf("failed to decode parts: %w", err)
	}

	if len(parts) == 0 {
		return nil, model.ErrNotFound
	}

	var result []*model.Part
	for _, p := range parts {
		result = append(result, lo.ToPtr(repoConverter.PartToModel(*p)))
	}

	return result, nil
}
