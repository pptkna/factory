package converter

import (
	"time"

	"github.com/pptkna/rocket-factory/order/internal/model"
	inventoryV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/samber/lo"
)

func PartsFilterToProto(partsFilter model.PartFiters) *inventoryV1.PartsFilter {
	categories := make([]inventoryV1.Category, len(partsFilter.Categories))
	for i, c := range partsFilter.Categories {
		categories[i] = CategoryToProto(c)
	}

	return &inventoryV1.PartsFilter{
		Uuids:                 partsFilter.Uuids,
		Names:                 partsFilter.Names,
		Categories:            categories,
		ManufacturerCountries: partsFilter.ManufacturerCountries,
		Tags:                  partsFilter.Tags,
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

func PartsListToModel(parts []*inventoryV1.Part) []model.Part {
	partsModel := make([]model.Part, len(parts))
	for i, p := range parts {
		partsModel[i] = PartToModel(p)
	}
	return partsModel
}

func PartToModel(part *inventoryV1.Part) model.Part {
	var dimensions *model.Dimensions
	if d := part.Dimensions; d != nil {
		dimensions = &model.Dimensions{
			Length: float32(d.GetLength()),
			Width:  float32(d.GetWidth()),
			Height: float32(d.GetHeight()),
			Weight: float32(d.GetWeight()),
		}
	}

	var manufacturer *model.Manufacturer
	if m := part.Manufacturer; m != nil {
		manufacturer = &model.Manufacturer{
			Name:    m.GetName(),
			Country: m.GetCountry(),
			Website: m.GetWebsite(),
		}
	}

	var createdAt *time.Time
	if part.CreatedAt != nil {
		createdAt = lo.ToPtr(part.CreatedAt.AsTime())
	}

	var updatedAt *time.Time
	if part.UpdatedAt != nil {
		updatedAt = lo.ToPtr(part.UpdatedAt.AsTime())
	}

	return model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         float32(part.Price),
		StockQuantity: int(part.StockQuantity),
		Category:      CategoryToModel(part.Category),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          part.Tags,
		Metadata:      MetadataToModel(part.Metadata),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
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

func MetadataToModel(metadata map[string]*inventoryV1.Value) map[string]*model.Value {
	md := make(map[string]*model.Value, len(metadata))
	for k, m := range metadata {
		md[k] = ValueToModel(m)
	}

	return nil
}

func ValueToModel(v *inventoryV1.Value) *model.Value {
	if v == nil {
		return nil
	}

	out := &model.Value{}

	switch val := v.Kind.(type) {
	case *inventoryV1.Value_StringValue:
		out.String = &val.StringValue
	case *inventoryV1.Value_Int64Value:
		out.Int64 = &val.Int64Value
	case *inventoryV1.Value_DoubleValue:
		out.Double = &val.DoubleValue
	case *inventoryV1.Value_BoolValue:
		out.Bool = &val.BoolValue
	}

	return out
}

func PaymentMethodToProto(paymentMethod model.PaymentMethod) paymentV1.PaymentMethod {
	switch paymentMethod {
	case model.PaymentMethodCard:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN
	}
}
