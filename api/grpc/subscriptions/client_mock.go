package subscriptions

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type clientMock struct{}

func newClientMock() ServiceClient {
	return clientMock{}
}

func (cm clientMock) Create(ctx context.Context, req *CreateRequest, opts ...grpc.CallOption) (resp *CreateResponse, err error) {
	resp = &CreateResponse{}
	switch req.Description {
	case "fail":
		err = status.Error(codes.Internal, "internal failure")
	case "fail_auth":
		err = status.Error(codes.Unauthenticated, "authentication failure")
	case "invalid":
		err = status.Error(codes.InvalidArgument, "invalid subscription condition")
	case "limit_reached":
		err = status.Error(codes.ResourceExhausted, "subscriptions count limit reached")
	case "busy":
		err = status.Error(codes.Unavailable, "retry the operation")
	default:
		resp.Id = "sub0"
	}
	return
}

func (cm clientMock) Read(ctx context.Context, req *ReadRequest, opts ...grpc.CallOption) (resp *ReadResponse, err error) {
	resp = &ReadResponse{}
	switch req.Id {
	case "fail":
		err = status.Error(codes.Internal, "internal failure")
	case "fail_auth":
		err = status.Error(codes.Unauthenticated, "authentication failure")
	case "missing":
		err = status.Error(codes.NotFound, "subscription not found")
	default:
		resp.Description = "subscription"
		resp.Enabled = true
		resp.Cond = &ConditionOutput{
			Cond: &ConditionOutput_Gc{
				Gc: &GroupConditionOutput{
					Logic: GroupLogic_Or,
					Group: []*ConditionOutput{
						{
							Not: true,
							Cond: &ConditionOutput_Tc{
								Tc: &TextConditionOutput{
									Key:  "k0",
									Term: "p0",
								},
							},
						},
						{
							Cond: &ConditionOutput_Tc{
								Tc: &TextConditionOutput{
									Key:   "k1",
									Term:  "p1",
									Exact: true,
								},
							},
						},
					},
				},
			},
		}
	}
	return
}

func (cm clientMock) Update(ctx context.Context, req *UpdateRequest, opts ...grpc.CallOption) (resp *UpdateResponse, err error) {
	resp = &UpdateResponse{}
	switch req.Id {
	case "fail":
		err = status.Error(codes.Internal, "internal failure")
	case "fail_auth":
		err = status.Error(codes.Unauthenticated, "authentication failure")
	case "missing":
		err = status.Error(codes.NotFound, "subscription not found")
	}
	return
}

func (cm clientMock) Delete(ctx context.Context, req *DeleteRequest, opts ...grpc.CallOption) (resp *DeleteResponse, err error) {
	resp = &DeleteResponse{}
	switch req.Id {
	case "fail":
		err = status.Error(codes.Internal, "internal failure")
	case "fail_auth":
		err = status.Error(codes.Unauthenticated, "authentication failure")
	case "missing":
		err = status.Error(codes.NotFound, "subscription not found")
	}
	return
}

func (cm clientMock) SearchOwn(ctx context.Context, req *SearchOwnRequest, opts ...grpc.CallOption) (resp *SearchOwnResponse, err error) {
	resp = &SearchOwnResponse{}
	switch req.Cursor {
	case "":
		resp.Ids = []string{
			"sub0",
			"sub1",
		}
	case "fail":
		err = status.Error(codes.Internal, "internal failure")
	case "fail_auth":
		err = status.Error(codes.Unauthenticated, "authentication failure")
	}
	return
}
