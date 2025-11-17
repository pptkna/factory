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

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.partService.GetPart(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
		}
		return nil, status.Errorf(codes.Internal, "internal error with UUID %s", req.GetUuid())
	}

	return &inventoryV1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
