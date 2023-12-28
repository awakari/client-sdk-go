package permits

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

	// GetUsage returns the current group/user spent counts for the specified subject.
	// A client should specify the empty user id to get the group-level value.
	GetUsage(ctx context.Context, userId string, subj usage.Subject) (u usage.Usage, err error)
}

type service struct {
	client ServiceClient
}

var ErrInternal = errors.New("internal failure")

var ErrUnavailable = errors.New("unavailable")

func NewService(client ServiceClient) Service {
	return service{
		client: client,
	}
}

func (svc service) GetUsage(ctx context.Context, userId string, subj usage.Subject) (u usage.Usage, err error) {
	ctx = auth.SetOutgoingAuthInfo(ctx, userId)
	req := GetUsageRequest{}
	req.Subj, err = subject.Encode(subj)
	var resp *GetUsageResponse
	if err == nil {
		resp, err = svc.client.GetUsage(ctx, &req)
		err = decodeError(err)
	}
	if err == nil {
		u.Count = resp.Count
		u.CountTotal = resp.CountTotal
		u.Since = resp.Since.AsTime()
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
		case s.Code() == codes.Unavailable:
			dst = fmt.Errorf("%w: %s", ErrUnavailable, src)
		default:
			dst = fmt.Errorf("%w: %s", ErrInternal, src)
		}
	}
	return
}
