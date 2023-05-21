package messages

import (
	"context"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type clientMock struct {
}

func newClientMock() ServiceClient {
	return clientMock{}
}

func (cm clientMock) Receive(ctx context.Context, opts ...grpc.CallOption) (Service_ReceiveClient, error) {
	return newStreamMock(), nil
}

type streamMock struct {
	subId string
	ack   bool
}

func newStreamMock() Service_ReceiveClient {
	return &streamMock{}
}

func (sm *streamMock) Send(req *ReceiveRequest) (err error) {
	start, ack := req.GetStart(), req.GetAck()
	switch {
	case start != nil:
		sm.subId = start.SubId
	case ack != nil:
		sm.ack = true
	}
	return
}

func (sm *streamMock) Recv() (msg *pb.CloudEvent, err error) {
	switch {
	case sm.ack && sm.subId == "fail_ack":
		err = status.Error(codes.Internal, "internal failure")
	case sm.subId == "fail_auth":
		err = status.Error(codes.Unauthenticated, "authentication failure")
	case sm.subId == "fail":
		err = status.Error(codes.Internal, "internal failure")
	case sm.subId == "missing":
		err = status.Error(codes.NotFound, "subscription was not found")
	default:
		msg = &pb.CloudEvent{
			Id:          "msg0",
			Source:      "source0",
			SpecVersion: "specversion0",
			Type:        "type0",
			Attributes:  map[string]*pb.CloudEventAttributeValue{},
			Data: &pb.CloudEvent_TextData{
				TextData: "data0",
			},
		}
	}
	return
}

func (sm *streamMock) Header() (metadata.MD, error) {
	//TODO implement me
	panic("implement me")
}

func (sm *streamMock) Trailer() metadata.MD {
	//TODO implement me
	panic("implement me")
}

func (sm *streamMock) CloseSend() error {
	return nil
}

func (sm *streamMock) Context() context.Context {
	//TODO implement me
	panic("implement me")
}

func (sm *streamMock) SendMsg(m interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (sm *streamMock) RecvMsg(m interface{}) error {
	//TODO implement me
	panic("implement me")
}
