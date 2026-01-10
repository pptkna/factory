package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/pptkna/rocket-factory/notification/internal/model"
	eventsV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/events/v1"
)

type orderAssembledDecoder struct{}

func NewOrderAssembledDecoder() *orderAssembledDecoder {
	return &orderAssembledDecoder{}
}

func (d *orderAssembledDecoder) Decode(data []byte) (*model.OrderAssembled, error) {
	var pb eventsV1.OrderAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return nil, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return &model.OrderAssembled{
		EventUUID:    pb.EventUuid,
		OrderUUID:    pb.OrderUuid,
		UserUUID:     pb.UserUuid,
		BuildTimeSec: pb.BuildTimeSec,
	}, nil
}
