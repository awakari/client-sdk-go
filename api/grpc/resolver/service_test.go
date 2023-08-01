package resolver

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
		msgs      []*pb.CloudEvent
		ackCount0 uint32
		ackCount1 uint32
		err       error
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
			ackCount0: 3,
			ackCount1: 3,
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
			ackCount0: 1,
			err:       ErrInternal,
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
			ackCount0: 2,
			err:       limits.ErrReached,
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
			err: auth.ErrAuth,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			w, err := svc.OpenWriter(context.TODO(), "user0")
			require.Nil(t, err)
			var ackCount uint32
			ackCount, err = w.WriteBatch(c.msgs)
			assert.Equal(t, c.ackCount0, ackCount)
			require.Nil(t, err)
			ackCount, err = w.WriteBatch(c.msgs)
			assert.Equal(t, c.ackCount1, ackCount)
			assert.ErrorIs(t, err, c.err)
			err = w.Close()
		})
	}
}
