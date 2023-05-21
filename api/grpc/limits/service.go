package limits

import (
	"context"
	"errors"
	"fmt"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/awakari/client-sdk-go/api/grpc/subject"
	"github.com/awakari/client-sdk-go/model/usage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	Get(ctx context.Context, userId string, subj usage.Subject) (l usage.Limit, err error)
}

type service struct {
	client ServiceClient
}

var ErrInternal = errors.New("internal failure")

var ErrReached = errors.New("usage limit reached")

func NewService(client ServiceClient) Service {
	return service{
		client: client,
	}
}

func (svc service) Get(ctx context.Context, userId string, subj usage.Subject) (l usage.Limit, err error) {
	ctx = auth.SetOutgoingAuthInfo(ctx, userId)
	req := GetRequest{}
	req.Subj, err = subject.Encode(subj)
	var resp *GetResponse
	if err == nil {
		resp, err = svc.client.Get(ctx, &req)
		err = decodeError(err)
	}
	if err == nil {
		l.Count = resp.Count
		l.UserId = resp.UserId
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
		case s.Code() == codes.Unauthenticated:
			dst = fmt.Errorf("%w: %s", auth.ErrAuth, src)
		default:
			dst = fmt.Errorf("%w: %s", ErrInternal, src)
		}
	}
	return
}
