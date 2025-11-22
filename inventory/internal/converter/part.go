package converter

import (
	"github.com/pptkna/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/inventory/v1"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PartsFilterToModel(partsFilter *inventoryV1.PartsFilter) model.PartFiters {
	if partsFilter == nil {
		return model.PartFiters{}
	}

	categories := make([]model.Category, len(partsFilter.Categories))
	for i, c := range partsFilter.Categories {
		categories[i] = CategoryToModel(c)
	}

	return model.PartFiters{
		Uuids:                 partsFilter.Uuids,
		Names:                 partsFilter.Names,
		Categories:            categories,
		ManufacturerCountries: partsFilter.ManufacturerCountries,
		Tags:                  partsFilter.Tags,
	}
}

func CategoryToModel(category inventoryV1.Category) model.Category {
	switch category {
	case inventoryV1.Category_CATEGORY_UNKNOWN:
		return model.CategoryUnknown
	case inventoryV1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventoryV1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventoryV1.Category_CATEGORY_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnspecified
	}
}

func PartToProto(part model.Part) *inventoryV1.Part {
	var dimensions *inventoryV1.Dimensions
	if part.Dimensions != nil {
		dimensions = lo.ToPtr(DimensionsToProto(*part.Dimensions))
	}

	var manufacturer *inventoryV1.Manufacturer
	if part.Manufacturer != nil {
		manufacturer = lo.ToPtr(ManufacturerToProto(*part.Manufacturer))
	}

	var createdAt *timestamppb.Timestamp
	if part.CreatedAt != nil {
		createdAt = timestamppb.New(*part.CreatedAt)
	}

	var updatedAt *timestamppb.Timestamp
	if part.UpdatedAt != nil {
		updatedAt = timestamppb.New(*part.UpdatedAt)
	}

	return &inventoryV1.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         float64(part.Price),
		StockQuantity: int64(part.StockQuantity),
		Category:      CategoryToProto(part.Category),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          part.Tags,
		Metadata:      MetadataToProto(part.Metadata),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func CategoryToProto(category model.Category) inventoryV1.Category {
	switch category {
	case model.CategoryUnknown:
		return inventoryV1.Category_CATEGORY_UNKNOWN
	case model.CategoryEngine:
		return inventoryV1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return inventoryV1.Category_CATEGORY_UNSPECIFIED
	}
}

func DimensionsToProto(dimension model.Dimensions) inventoryV1.Dimensions {
	return inventoryV1.Dimensions{
		Length: float64(dimension.Length),
		Width:  float64(dimension.Width),
		Height: float64(dimension.Height),
		Weight: float64(dimension.Weight),
	}
}

func ManufacturerToProto(manufacturer model.Manufacturer) inventoryV1.Manufacturer {
	return inventoryV1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func MetadataToProto(metadata map[string]*model.Value) map[string]*inventoryV1.Value {
	md := make(map[string]*inventoryV1.Value, len(metadata))
	for k, m := range metadata {
		md[k] = ValueToProto(m)
	}

	return nil
}

func ValueToProto(value *model.Value) *inventoryV1.Value {
	if value != nil {
		v := &inventoryV1.Value{}

		switch {
		case value.String != nil:
			v.Kind = &inventoryV1.Value_StringValue{StringValue: *value.String}
		case value.Double != nil:
			v.Kind = &inventoryV1.Value_DoubleValue{DoubleValue: *value.Double}
		case value.Int64 != nil:
			v.Kind = &inventoryV1.Value_Int64Value{Int64Value: *value.Int64}
		case value.Bool != nil:
			v.Kind = &inventoryV1.Value_BoolValue{BoolValue: *value.Bool}
		}

		return v
	}

	return nil
}
