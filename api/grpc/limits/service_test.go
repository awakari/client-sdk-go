package limits

import (
	"context"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/awakari/client-sdk-go/api/grpc/subject"
	"github.com/awakari/client-sdk-go/model/usage"
	"github.com/stretchr/testify/assert"
	"testing"
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
			subject: usage.SubjectPublishMessages,
			lim: usage.Limit{
				Count:  1,
				UserId: "user0",
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
			subject: usage.SubjectPublishMessages,
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
