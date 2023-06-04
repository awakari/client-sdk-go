package reader

import (
	"context"
	"fmt"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
)

type clientMock struct {
}

func newClientMock() ServiceClient {
	return clientMock{}
}

func (cm clientMock) Read(ctx context.Context, opts ...grpc.CallOption) (Service_ReadClient, error) {
	md, _ := metadata.FromOutgoingContext(ctx)
	userId := md.Get("x-awakari-user-id")
	return newStreamMock(userId[0]), nil
}

type streamMock struct {
	userId    string
	subId     string
	batchSize uint32
}

func newStreamMock(userId string) Service_ReadClient {
	return &streamMock{
		userId: userId,
	}
}

func (sm *streamMock) Send(req *ReadRequest) (err error) {
	start, ack := req.GetStart(), req.GetAck()
	switch {
	case sm.userId == "fail_auth":
		err = io.EOF
	case start != nil:
		sm.subId = start.SubId
		sm.batchSize = start.BatchSize
	case ack != nil:
	}
	return
}

func (sm *streamMock) Recv() (resp *ReadResponse, err error) {
	switch {
	case sm.subId == "fail":
		err = status.Error(codes.Internal, "internal failure")
	case sm.subId == "missing":
		err = status.Error(codes.NotFound, "subscription not found")
	case sm.batchSize == 0:
		err = status.Error(codes.InvalidArgument, "batch size should be > 0")
	case sm.userId == "fail_auth":
		err = status.Error(codes.Unauthenticated, "auth failure")
	default:
		resp = &ReadResponse{}
		for i := uint32(0); i < sm.batchSize; i++ {
			msg := &pb.CloudEvent{
				Id:          fmt.Sprintf("msg%d", i),
				Source:      "source0",
				SpecVersion: "specversion0",
				Type:        "type0",
				Attributes:  map[string]*pb.CloudEventAttributeValue{},
				Data: &pb.CloudEvent_TextData{
					TextData: "data0",
				},
			}
			resp.Msgs = append(resp.Msgs, msg)
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
