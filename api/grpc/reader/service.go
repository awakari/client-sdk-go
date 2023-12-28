package reader

import (
	"context"
	"errors"
	"fmt"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/awakari/client-sdk-go/model"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type Service interface {
	OpenReader(ctx context.Context, userId, subId string, batchSize uint32) (rs model.Reader[[]*pb.CloudEvent], err error)
	OpenAckReader(ctx context.Context, userId, subId string, batchSize uint32) (r model.AckReader[[]*pb.CloudEvent], err error)
}

type service struct {
	client ServiceClient
}

var ErrInternal = errors.New("internal failure")

var ErrInvalidRequest = errors.New("invalid request")

var ErrNotFound = errors.New("subscription not found")

var ErrUnavailable = errors.New("unavailable")

func NewService(client ServiceClient) Service {
	return service{
		client: client,
	}
}

func (svc service) OpenReader(ctx context.Context, userId, subId string, batchSize uint32) (rs model.Reader[[]*pb.CloudEvent], err error) {
	ctx = auth.SetOutgoingAuthInfo(ctx, userId)
	var stream Service_ReadClient
	stream, err = svc.client.Read(ctx)
	if err == nil {
		reqStart := ReadRequest{
			Command: &ReadRequest_Start{
				Start: &ReadCommandStart{
					SubId:     subId,
					BatchSize: batchSize,
				},
			},
		}
		err = stream.Send(&reqStart)
	}
	if err == nil {
		rs = newStreamReader(stream)
	}
	err = decodeError(err)
	return
}

func (svc service) OpenAckReader(ctx context.Context, userId, subId string, batchSize uint32) (r model.AckReader[[]*pb.CloudEvent], err error) {
	ctx = auth.SetOutgoingAuthInfo(ctx, userId)
	var stream Service_ReadClient
	stream, err = svc.client.Read(ctx)
	if err == nil {
		reqStart := ReadRequest{
			Command: &ReadRequest_Start{
				Start: &ReadCommandStart{
					SubId:     subId,
					BatchSize: batchSize,
				},
			},
		}
		err = stream.Send(&reqStart)
	}
	if err == nil {
		r = newStreamAckReader(stream)
	}
	err = decodeError(err)
	return
}

func decodeError(src error) (dst error) {
	switch {
	case src == nil:
	case status.Code(src) == codes.OK:
	case src == io.EOF:
		dst = src
	case status.Code(src) == codes.DeadlineExceeded:
		dst = context.DeadlineExceeded
	case status.Code(src) == codes.InvalidArgument:
		dst = fmt.Errorf("%w: %s", ErrInvalidRequest, src)
	case status.Code(src) == codes.Unauthenticated:
		dst = fmt.Errorf("%w: %s", auth.ErrAuth, src)
	case status.Code(src) == codes.NotFound:
		dst = fmt.Errorf("%w: %s", ErrNotFound, src)
	case status.Code(src) == codes.Unavailable:
		dst = fmt.Errorf("%w: %s", ErrUnavailable, src)
	default:
		dst = fmt.Errorf("%w: %s", ErrInternal, src)
	}
	return
}
