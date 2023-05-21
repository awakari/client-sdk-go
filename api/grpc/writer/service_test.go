package writer

import (
	"context"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/awakari/client-sdk-go/api/grpc/limits"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_Write(t *testing.T) {
	svc := NewService(newClientMock())
	cases := map[string]struct {
		msgs     []*pb.CloudEvent
		ackCount uint32
		err      error
		errMsg   string
	}{
		"ok": {
			msgs: []*pb.CloudEvent{
				{
					Id: "msg0",
				},
				{
					Id: "msg1",
				},
				{
					Id: "msg2",
				},
			},
			ackCount: 3,
		},
		"fail": {
			msgs: []*pb.CloudEvent{
				{
					Id: "msg0",
				},
				{
					Id: "fail",
				},
				{
					Id: "msg2",
				},
			},
			ackCount: 1,
			err:      ErrInternal,
			errMsg:   "internal failure: rpc error: code = Internal desc = internal failure",
		},
		"limit reached": {
			msgs: []*pb.CloudEvent{
				{
					Id: "msg0",
				},
				{
					Id: "msg1",
				},
				{
					Id: "limit_reached",
				},
			},
			ackCount: 2,
			err:      limits.ErrReached,
			errMsg:   "usage limit reached: rpc error: code = ResourceExhausted desc = usage limit reached",
		},
		"fail auth": {
			msgs: []*pb.CloudEvent{
				{
					Id: "fail_auth",
				},
				{
					Id: "msg1",
				},
				{
					Id: "msg2",
				},
			},
			err:    auth.ErrAuth,
			errMsg: "authentication failure: rpc error: code = Unauthenticated desc = authentication failure",
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			ws, err := svc.OpenStream(context.TODO(), "user0")
			require.Nil(t, err)
			var ackCount uint32
			ackCount, err = ws.WriteBatch(c.msgs)
			assert.Equal(t, c.ackCount, ackCount)
			assert.ErrorIs(t, err, c.err)
			if err != nil {
				assert.Equal(t, c.errMsg, err.Error())
			}
		})
	}
}
