package subscriptions

import (
	"context"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/awakari/client-sdk-go/api/grpc/limits"
	"github.com/awakari/client-sdk-go/model/subscription"
	"github.com/awakari/client-sdk-go/model/subscription/condition"
	"time"
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
		err = ErrUnavailable
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
		subData.Expires = time.Date(2023, 10, 4, 11, 44, 55, 0, time.UTC)
		subData.Condition = condition.
			NewBuilder().
			Any([]condition.Condition{
				condition.
					NewBuilder().
					BuildTextCondition(),
				condition.
					NewBuilder().
					LessThanOrEqual(42).
					BuildNumberCondition(),
			}).
			BuildGroupCondition()
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

func (sm serviceMock) Search(ctx context.Context, userId string, q subscription.Query, cursor subscription.Cursor) (ids []string, err error) {
	switch cursor.Id {
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
