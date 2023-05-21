package writer

import (
	"context"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type streamMock struct {
	lastReqMsgs []*pb.CloudEvent
}

func newStreamMock() Service_SubmitMessagesClient {
	return &streamMock{}
}

func (sm *streamMock) Send(req *SubmitMessagesRequest) (err error) {
	sm.lastReqMsgs = req.Msgs
	return
}

func (sm *streamMock) Recv() (resp *SubmitMessagesResponse, err error) {
	resp = &SubmitMessagesResponse{}
	for _, msg := range sm.lastReqMsgs {
		switch msg.Id {
		case "fail":
			err = status.Error(codes.Internal, "internal failure")
		case "fail_auth":
			err = status.Error(codes.Unauthenticated, "authentication failure")
		case "limit_reached":
			err = status.Error(codes.ResourceExhausted, "usage limit reached")
		}
		if err == nil {
			resp.AckCount++
		} else {
			resp.Err = err.Error()
			break
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
