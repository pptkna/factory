package v1

import (
	def "github.com/pptkna/rocket-factory/order/internal/client/grpc"
	generatedInventoryClientV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/inventory/v1"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	generatedClient generatedInventoryClientV1.InventoryServiceClient
}

func NewClient(generatedClient generatedInventoryClientV1.InventoryServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
