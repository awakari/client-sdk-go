package permits

import (
	"context"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/awakari/client-sdk-go/api/grpc/subject"
	"github.com/awakari/client-sdk-go/model/usage"
	"time"
)

type serviceMock struct {
}

func NewServiceMock() Service {
	return serviceMock{}
}

func (sm serviceMock) GetUsage(ctx context.Context, userId string, subj usage.Subject) (u usage.Usage, err error) {
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
		default:
			u.Count = 1
			u.CountTotal = 2
			u.Since = time.Date(2023, 05, 07, 04, 57, 20, 0, time.UTC)
		}
	}
	return
}
