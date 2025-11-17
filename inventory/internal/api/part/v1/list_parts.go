package v1

import (
	"context"
	"errors"

	"github.com/pptkna/rocket-factory/inventory/internal/converter"
	"github.com/pptkna/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filter := req.GetFilter()

	parts, err := a.partService.ListParts(ctx, converter.PartsFilterToProto(filter))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "parts not found")
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	partsProto := make([]*inventoryV1.Part, len(parts))
	for i, p := range parts {
		partsProto[i] = converter.PartToProto(p)
	}

	return &inventoryV1.ListPartsResponse{
		Parts: partsProto,
	}, nil
}

// func (s *inventoryService) ListParts(_ context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
// 	s.mu.RLock()

// 	parts := make([]*inventory_v1.Part, 0, len(s.parts))

// 	for _, p := range s.parts {
// 		parts = append(parts, p)
// 	}

// 	s.mu.RUnlock()

// 	filters := req.GetFilter()
// 	if len(filters.Uuids) == 0 && len(filters.Categories) == 0 && len(filters.ManufacturerCountries) == 0 && len(filters.Names) == 0 && len(filters.Tags) == 0 {
// 		return &inventory_v1.ListPartsResponse{Parts: parts}, nil
// 	}

// 	filteredParts := make([]*inventory_v1.Part, 0, len(parts))
// 	for _, p := range parts {
// 		if matchesFilter(p, filters) {
// 			filteredParts = append(filteredParts, p)
// 		}
// 	}

// 	return &inventory_v1.ListPartsResponse{
// 		Parts: filteredParts,
// 	}, nil
// }
