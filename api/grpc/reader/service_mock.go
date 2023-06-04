package reader

import (
	"context"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/awakari/client-sdk-go/model"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
)

type serviceMock struct {
}

func NewServiceMock() Service {
	return serviceMock{}
}

func (sm serviceMock) OpenReader(ctx context.Context, userId, subId string, batchSize uint32) (rs model.Reader[[]*pb.CloudEvent], err error) {
	switch subId {
	case "fail":
		err = ErrInternal
	case "fail_auth":
		err = auth.ErrAuth
	case "missing":
		err = ErrNotFound
	}
	if err == nil {
		rs = newStreamReaderMock(subId, batchSize)
	}
	return
}
