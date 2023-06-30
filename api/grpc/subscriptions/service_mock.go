package subscriptions

import (
	"context"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/awakari/client-sdk-go/api/grpc/limits"
	"github.com/awakari/client-sdk-go/model/subscription"
	"github.com/awakari/client-sdk-go/model/subscription/condition"
)

type serviceMock struct {
}

func NewServiceMock() Service {
	return serviceMock{}
}

func (sm serviceMock) Create(ctx context.Context, userId string, subData subscription.Data) (id string, err error) {
	switch subData.Description {
	case "fail":
		err = ErrInternal
	case "fail_auth":
		err = auth.ErrAuth
	case "invalid":
		err = ErrInvalid
	case "limit_reached":
		err = limits.ErrReached
	case "busy":
		err = ErrBusy
	default:
		id = "sub0"
	}
	return
}

func (sm serviceMock) Read(ctx context.Context, userId, subId string) (subData subscription.Data, err error) {
	switch subId {
	case "fail":
		err = ErrInternal
	case "fail_auth":
		err = auth.ErrAuth
	case "missing":
		err = ErrNotFound
	default:
		subData.Description = "my subscription"
		subData.Enabled = true
		subData.Condition = condition.
			NewBuilder().
			BuildTextCondition()
	}
	return
}

func (sm serviceMock) Update(ctx context.Context, userId, subId string, subData subscription.Data) (err error) {
	switch subId {
	case "fail":
		err = ErrInternal
	case "fail_auth":
		err = auth.ErrAuth
	case "missing":
		err = ErrNotFound
	}
	return
}

func (sm serviceMock) Delete(ctx context.Context, userId, subId string) (err error) {
	switch subId {
	case "fail":
		err = ErrInternal
	case "fail_auth":
		err = auth.ErrAuth
	case "missing":
		err = ErrNotFound
	}
	return
}

func (sm serviceMock) SearchOwn(ctx context.Context, userId string, limit uint32, cursor string) (ids []string, err error) {
	switch cursor {
	case "":
		ids = []string{
			"sub0",
			"sub1",
		}
	case "fail":
		err = ErrInternal
	case "fail_auth":
		err = auth.ErrAuth
	}
	return
}
