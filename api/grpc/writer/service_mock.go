package writer

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

func (sm serviceMock) OpenStream(ctx context.Context, userId string) (ws model.WriteStream[*pb.CloudEvent], err error) {
	switch userId {
	case "fail":
		err = ErrInternal
	case "fail_auth":
		err = auth.ErrAuth
	}
	if err == nil {
		ws = newWriteStreamMock()
	}
	return
}
