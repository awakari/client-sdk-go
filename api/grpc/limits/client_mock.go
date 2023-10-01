package limits

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type clientMock struct {
}

func newClientMock() ServiceClient {
	return clientMock{}
}

func (cm clientMock) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (resp *GetResponse, err error) {
	resp = &GetResponse{}
	md, _ := metadata.FromOutgoingContext(ctx)
	userId := md.Get("x-awakari-user-id")
	switch userId[0] {
	case "fail":
		err = status.Error(codes.Internal, "internal failure")
	case "fail_auth":
		err = status.Error(codes.Unauthenticated, "authentication failure")
	case "with_expiration":
		resp.Count = 1
		resp.UserId = userId[0]
		resp.Expires = timestamppb.New(time.Date(2023, 10, 1, 20, 21, 35, 0, time.UTC))
	default:
		resp.Count = 1
		resp.UserId = userId[0]
	}
	return
}
