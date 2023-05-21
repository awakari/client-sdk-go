package messages

import (
	"github.com/awakari/client-sdk-go/model"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
)

type readStreamMock struct {
	userId string
	subId  string
}

func newReadStreamMock(userId, subId string) (rs model.ReadStream[*pb.CloudEvent]) {
	return readStreamMock{
		userId: userId,
		subId:  subId,
	}
}

func (r readStreamMock) Close() error {
	return nil
}

func (r readStreamMock) Read() (msg *pb.CloudEvent, err error) {
	switch r.subId {
	case "fail_read":
		err = ErrInternal
	case "missing":
		err = ErrNotFound
	default:
		msg = &pb.CloudEvent{
			Id:          "msg0",
			Source:      "source0",
			SpecVersion: "specversion0",
			Type:        "type0",
			Attributes:  map[string]*pb.CloudEventAttributeValue{},
			Data: &pb.CloudEvent_TextData{
				TextData: "data",
			},
		}
	}
	return
}
