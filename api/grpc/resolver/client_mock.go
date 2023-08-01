package resolver

import (
	"context"
	"google.golang.org/grpc"
)

type clientMock struct {
}

func newClientMock() ServiceClient {
	return clientMock{}
}

func (cm clientMock) SubmitMessages(ctx context.Context, opts ...grpc.CallOption) (Service_SubmitMessagesClient, error) {
	return newStreamMock(), nil
}
