package reader

import (
	"fmt"
	"github.com/awakari/client-sdk-go/model"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
)

type streamReaderMock struct {
	subId     string
	batchSize uint32
}

func newStreamReaderMock(subId string, batchSize uint32) (r model.Reader[[]*pb.CloudEvent]) {
	return streamReaderMock{
		subId:     subId,
		batchSize: batchSize,
	}
}

func (r streamReaderMock) Close() error {
	return nil
}

func (r streamReaderMock) Read() (msgs []*pb.CloudEvent, err error) {
	switch {
	case r.batchSize == 0:
		err = ErrInvalidRequest
	case r.subId == "fail_read":
		err = ErrInternal
	default:
		for i := uint32(0); i < r.batchSize; i++ {
			msg := &pb.CloudEvent{
				Id:          fmt.Sprintf("msg%d", i),
				Source:      "source0",
				SpecVersion: "specversion0",
				Type:        "type0",
				Attributes:  map[string]*pb.CloudEventAttributeValue{},
				Data: &pb.CloudEvent_TextData{
					TextData: "data",
				},
			}
			msgs = append(msgs, msg)
		}
	}
	return
}
