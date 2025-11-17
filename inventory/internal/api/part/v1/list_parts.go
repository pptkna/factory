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
