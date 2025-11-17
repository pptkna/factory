package converter

import (
	"github.com/pptkna/rocket-factory/inventory/internal/model"
	repoModel "github.com/pptkna/rocket-factory/inventory/internal/repository/model"
	"github.com/samber/lo"
)

func PartFitersToRepoModel(partFiters model.PartFiters) repoModel.PartFiters {
	categories := make([]repoModel.Category, len(partFiters.Categories))
	for i, c := range partFiters.Categories {
		categories[i] = CategoryToRepoModel(c)
	}

	return repoModel.PartFiters{
		Uuids:                 partFiters.Uuids,
		Names:                 partFiters.Names,
		Categories:            categories,
		ManufacturerCountries: partFiters.ManufacturerCountries,
		Tags:                  partFiters.Tags,
	}
}

func PartToRepoModel(part model.Part) repoModel.Part {
	var dimension *repoModel.Dimensions
	if part.Dimensions != nil {
		dimension = lo.ToPtr(DimensionsToRepoModel(*part.Dimensions))
	}

	var manufacturer *repoModel.Manufacturer
	if part.Manufacturer != nil {
		manufacturer = lo.ToPtr(ManufacturerToRepoModel(*part.Manufacturer))
	}

	return repoModel.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToRepoModel(part.Category),
		Dimensions:    dimension,
		Manufacturer:  manufacturer,
		Tags:          part.Tags,
		Metadata:      MetadataToRepoModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func CategoryToRepoModel(category model.Category) repoModel.Category {
	switch category {
	case model.CategoryUnknown:
		return repoModel.CategoryUnknown
	case model.CategoryEngine:
		return repoModel.CategoryEngine
	case model.CategoryFuel:
		return repoModel.CategoryFuel
	case model.CategoryPorthole:
		return repoModel.CategoryPorthole
	case model.CategoryWing:
		return repoModel.CategoryWing
	default:
		return repoModel.CategoryUnspecified
	}
}

func DimensionsToRepoModel(dimensions model.Dimensions) repoModel.Dimensions {
	return repoModel.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ManufacturerToRepoModel(manufacturer model.Manufacturer) repoModel.Manufacturer {
	return repoModel.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func MetadataToRepoModel(metadata map[string]*model.Value) map[string]*repoModel.Value {
	md := make(map[string]*repoModel.Value, len(metadata))
	for k, m := range metadata {
		md[k] = ValueToRepoModel(m)
	}

	return md
}

func ValueToRepoModel(value *model.Value) *repoModel.Value {
	if value != nil {
		return &repoModel.Value{
			String: value.String,
			Int64:  value.Int64,
			Double: value.Double,
			Bool:   value.Bool,
		}
	}

	return nil
}

func PartToModel(part repoModel.Part) model.Part {
	var dimension *model.Dimensions
	if part.Dimensions != nil {
		dimension = lo.ToPtr(DimensionsToModel(*part.Dimensions))
	}

	var manufacturer *model.Manufacturer
	if part.Manufacturer != nil {
		manufacturer = lo.ToPtr(ManufacturerToModel(*part.Manufacturer))
	}

	return model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToModel(part.Category),
		Dimensions:    dimension,
		Manufacturer:  manufacturer,
		Tags:          part.Tags,
		Metadata:      MetadataToModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func CategoryToModel(category repoModel.Category) model.Category {
	switch category {
	case repoModel.CategoryUnknown:
		return model.CategoryUnknown
	case repoModel.CategoryEngine:
		return model.CategoryEngine
	case repoModel.CategoryFuel:
		return model.CategoryFuel
	case repoModel.CategoryPorthole:
		return model.CategoryPorthole
	case repoModel.CategoryWing:
		return model.CategoryWing
	default:
		return model.CategoryUnspecified
	}
}

func DimensionsToModel(dimensions repoModel.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ManufacturerToModel(manufacturer repoModel.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func MetadataToModel(metadata map[string]*repoModel.Value) map[string]*model.Value {
	md := make(map[string]*model.Value, len(metadata))
	for k, m := range metadata {
		md[k] = ValueToModel(m)
	}

	return md
}

func ValueToModel(value *repoModel.Value) *model.Value {
	if value != nil {
		return &model.Value{
			String: value.String,
			Int64:  value.Int64,
			Double: value.Double,
			Bool:   value.Bool,
		}
	}

	return nil
}
