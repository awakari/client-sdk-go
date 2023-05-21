package limits

import (
	"context"
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

func (cm clientMock) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (resp *GetResponse, err error) {
	resp = &GetResponse{}
	md, _ := metadata.FromOutgoingContext(ctx)
	userId := md.Get("x-awakari-user-id")
	switch userId[0] {
	case "fail":
		err = status.Error(codes.Internal, "internal failure")
	case "fail_auth":
		err = status.Error(codes.Unauthenticated, "authentication failure")
	default:
		resp.Count = 1
		resp.UserId = userId[0]
	}
	return
}
