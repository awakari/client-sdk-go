package permits

import (
	"context"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/awakari/client-sdk-go/api/grpc/subject"
	"github.com/awakari/client-sdk-go/model/usage"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestService_GetUsage(t *testing.T) {
	svc := NewService(newClientMock())
	cases := map[string]struct {
		userId  string
		subject usage.Subject
		out     usage.Usage
		err     error
	}{
		"ok": {
			userId:  "user0",
			subject: usage.SubjectPublishMessages,
			out: usage.Usage{
				Count:      1,
				CountTotal: 2,
				Since:      time.Date(2023, 05, 07, 04, 57, 20, 0, time.UTC),
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
			lim, err := svc.GetUsage(context.TODO(), c.userId, c.subject)
			assert.Equal(t, c.out, lim)
			assert.ErrorIs(t, err, c.err)
		})
	}
}
