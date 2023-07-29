package resolver

import (
	"context"
	"errors"
	"fmt"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/awakari/client-sdk-go/api/grpc/limits"
	"github.com/awakari/client-sdk-go/model"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	OpenWriter(ctx context.Context, userId string) (w model.Writer[*pb.CloudEvent], err error)
}

type service struct {
	client ServiceClient
}

var ErrInternal = errors.New("internal failure")

func NewService(client ServiceClient) Service {
	return service{
		client: client,
	}
}

func (svc service) OpenWriter(ctx context.Context, userId string) (w model.Writer[*pb.CloudEvent], err error) {
	ctx = auth.SetOutgoingAuthInfo(ctx, userId)
	var stream Service_SubmitMessagesClient
	stream, err = svc.client.SubmitMessages(ctx)
	err = decodeError(err)
	if err == nil {
		w = newStreamWriter(stream)
	}
	return
}

func decodeError(src error) (dst error) {
	switch {
	case src == nil:
	default:
		s, isGrpcErr := status.FromError(src)
		switch {
		case !isGrpcErr:
			dst = src
		case s.Code() == codes.OK:
		case s.Code() == codes.ResourceExhausted:
			dst = fmt.Errorf("%w: %s", limits.ErrReached, src)
		case s.Code() == codes.Unauthenticated:
			dst = fmt.Errorf("%w: %s", auth.ErrAuth, src)
		default:
			dst = fmt.Errorf("%w: %s", ErrInternal, src)
		}
	}
	return
}
