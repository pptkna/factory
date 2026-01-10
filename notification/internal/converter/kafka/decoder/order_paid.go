package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/pptkna/rocket-factory/notification/internal/model"
	eventsV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewOrderPaidDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (*model.OrderPaid, error) {
	var pb eventsV1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return nil, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return &model.OrderPaid{
		EventUUID:       pb.EventUuid,
		OrderUUID:       pb.OrderUuid,
		UserUUID:        pb.UserUuid,
		PaymentMethod:   pb.PaymentMethod.String(),
		TransactionUUID: pb.TransactionUuid,
	}, nil
}
