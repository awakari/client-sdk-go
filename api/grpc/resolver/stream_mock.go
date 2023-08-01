package resolver

import (
	"context"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
)

type streamMock struct {
	lastReqMsgs []*pb.CloudEvent
	err         error
}

func newStreamMock() Service_SubmitMessagesClient {
	return &streamMock{}
}

func (sm *streamMock) Send(req *SubmitMessagesRequest) (err error) {
	switch sm.err {
	case nil:
		sm.lastReqMsgs = req.Msgs
	default:
		err = io.EOF
	}
	return
}

func (sm *streamMock) Recv() (resp *SubmitMessagesResponse, err error) {
	resp = &SubmitMessagesResponse{}
	switch sm.err {
	case nil:
		for _, msg := range sm.lastReqMsgs {
			switch msg.Id {
			case "fail":
				sm.err = status.Error(codes.Internal, "internal failure")
			case "fail_auth":
				sm.err = status.Error(codes.Unauthenticated, "authentication failure")
			case "limit_reached":
				sm.err = status.Error(codes.ResourceExhausted, "usage limit reached")
			default:
				resp.AckCount++
			}
			if sm.err != nil {
				break
			}
		}
	default:
		err = sm.err
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
