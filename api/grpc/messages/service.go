package messages

import (
	"context"
	"errors"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/awakari/client-sdk-go/model"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
)

type Service interface {
	Read(ctx context.Context, userId, subId string) (rs model.ReadStream[*pb.CloudEvent], err error)
}

type service struct {
	client ServiceClient
}

var ErrInternal = errors.New("internal failure")

var ErrNotFound = errors.New("not found")

func NewService(client ServiceClient) Service {
	return service{
		client: client,
	}
}

func (svc service) Read(ctx context.Context, userId, subId string) (rs model.ReadStream[*pb.CloudEvent], err error) {
	ctx = auth.SetOutgoingAuthInfo(ctx, userId)
	var stream Service_ReceiveClient
	stream, err = svc.client.Receive(ctx)
	if err == nil {
		rs = newReadStream(stream)
		reqStart := ReceiveRequest{
			Command: &ReceiveRequest_Start{
				Start: &ReceiveCommandStart{
					SubId: subId,
				},
			},
		}
		err = stream.Send(&reqStart)
	}
	return
}
