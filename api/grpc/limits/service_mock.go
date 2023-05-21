package limits

import (
	"context"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/awakari/client-sdk-go/api/grpc/subject"
	"github.com/awakari/client-sdk-go/model/usage"
)

type serviceMock struct {
}

func NewServiceMock() Service {
	return serviceMock{}
}

func (sm serviceMock) Get(ctx context.Context, userId string, subj usage.Subject) (l usage.Limit, err error) {
	switch subj {
	case usage.SubjectUndefined:
		err = subject.ErrInvalidSubject
	}
	if err == nil {
		switch userId {
		case "fail":
			err = ErrInternal
		case "fail_auth":
			err = auth.ErrAuth
		case "group_missing":
		case "user_missing":
			l.Count = 1
		default:
			l.Count = 2
			l.UserId = userId
		}
	}
	return
}
