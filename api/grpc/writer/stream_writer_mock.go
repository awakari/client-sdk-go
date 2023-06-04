package writer

import (
	"github.com/awakari/client-sdk-go/api/grpc/limits"
	"github.com/awakari/client-sdk-go/model"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
)

type streamWriterMock struct {
}

func newStreamWriterMock() model.Writer[*pb.CloudEvent] {
	return streamWriterMock{}
}

func (w streamWriterMock) Close() error {
	return nil
}

func (w streamWriterMock) WriteBatch(msgs []*pb.CloudEvent) (ackCount uint32, err error) {
	for _, msg := range msgs {
		switch msg.Id {
		case "fail":
			err = ErrInternal
		case "limit_reached":
			err = limits.ErrReached
		}
		if err != nil {
			break
		}
		ackCount++
	}
	return
}
