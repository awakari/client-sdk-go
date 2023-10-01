package limits

import (
	"context"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/awakari/client-sdk-go/api/grpc/subject"
	"github.com/awakari/client-sdk-go/model/usage"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestService_Get(t *testing.T) {
	svc := NewService(newClientMock())
	cases := map[string]struct {
		userId  string
		subject usage.Subject
		lim     usage.Limit
		err     error
	}{
		"ok": {
			userId:  "user0",
			subject: usage.SubjectPublishEvents,
			lim: usage.Limit{
				Count:  1,
				UserId: "user0",
			},
		},
		"with expiration": {
			userId:  "with_expiration",
			subject: usage.SubjectPublishEvents,
			lim: usage.Limit{
				Count:   1,
				UserId:  "with_expiration",
				Expires: time.Date(2023, 10, 1, 20, 21, 35, 0, time.UTC),
			},
		},
		"invalid subject": {
			userId:  "user0",
			subject: usage.SubjectUndefined,
			err:     subject.ErrInvalidSubject,
		},
		"fail": {
			userId:  "fail",
			subject: usage.SubjectSubscriptions,
			err:     ErrInternal,
		},
		"fail auth": {
			userId:  "fail_auth",
			subject: usage.SubjectPublishEvents,
			err:     auth.ErrAuth,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			lim, err := svc.Get(context.TODO(), c.userId, c.subject)
			assert.Equal(t, c.lim, lim)
			assert.ErrorIs(t, err, c.err)
		})
	}
}
