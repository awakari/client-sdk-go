package permits

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

func (cm clientMock) GetUsage(ctx context.Context, req *GetUsageRequest, opts ...grpc.CallOption) (resp *GetUsageResponse, err error) {
	resp = &GetUsageResponse{}
	md, _ := metadata.FromOutgoingContext(ctx)
	userId := md.Get("x-awakari-user-id")
	switch userId[0] {
	case "fail":
		err = status.Error(codes.Internal, "internal failure")
	case "fail_auth":
		err = status.Error(codes.Unauthenticated, "authentication failure")
	default:
		resp.Count = 1
		resp.CountTotal = 2
		resp.Since = timestamppb.New(time.Date(2023, 05, 07, 04, 57, 20, 0, time.UTC))
	}
	return
}
