package decoder

import (
	"fmt"

	"github.com/pptkna/rocket-factory/assembly/internal/model"
	orderEventsV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type decoder struct{}

func NewDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (*model.OrderPaidEvent, error) {
	var pb orderEventsV1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return nil, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return &model.OrderPaidEvent{
		EventUUID:       pb.EventUuid,
		OrderUUID:       pb.OrderUuid,
		UserUUID:        pb.UserUuid,
		PaymentMethod:   model.PaymentMethod(pb.PaymentMethod.String()),
		TransactionUUID: pb.TransactionUuid,
	}, nil
}
