package v1

import (
	"context"

	clientConverter "github.com/pptkna/rocket-factory/order/internal/client/converter"
	"github.com/pptkna/rocket-factory/order/internal/model"
	generatedInventoryClientV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartFiters) ([]model.Part, error) {
	res, err := c.generatedClient.ListParts(ctx, &generatedInventoryClientV1.ListPartsRequest{
		Filter: clientConverter.PartsFilterToProto(filter),
	})
	if err != nil {
		return nil, err
	}

	return clientConverter.PartsListToModel(res.Parts), nil
}
