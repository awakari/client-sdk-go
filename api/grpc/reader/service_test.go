package reader

import (
	"context"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestService_Read(t *testing.T) {
	svc := NewService(newClientMock())
	cases := map[string]struct {
		userId    string
		subId     string
		batchSize uint32
		msgs      []*pb.CloudEvent
		errOpen   error
		errRead   error
	}{
		"ok": {
			userId:    "user0",
			subId:     "sub0",
			batchSize: 2,
			msgs: []*pb.CloudEvent{
				{
					Id:          "msg0",
					Source:      "source0",
					SpecVersion: "specversion0",
					Type:        "type0",
					Attributes:  map[string]*pb.CloudEventAttributeValue{},
					Data: &pb.CloudEvent_TextData{
						TextData: "data0",
					},
				},
				{
					Id:          "msg1",
					Source:      "source0",
					SpecVersion: "specversion0",
					Type:        "type0",
					Attributes:  map[string]*pb.CloudEventAttributeValue{},
					Data: &pb.CloudEvent_TextData{
						TextData: "data0",
					},
				},
			},
		},
		"invalid batch size": {
			userId:    "user0",
			subId:     "sub0",
			batchSize: 0,
			errRead:   ErrInvalidRequest,
		},
		"fail": {
			userId:    "user0",
			subId:     "fail",
			batchSize: 2,
			errRead:   ErrInternal,
		},
		"fail auth": {
			userId:    "fail_auth",
			subId:     "sub0",
			batchSize: 2,
			errOpen:   io.EOF,
			errRead:   auth.ErrAuth,
		},
		"missing subscription": {
			userId:    "user0",
			subId:     "missing",
			batchSize: 2,
			errRead:   ErrNotFound,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			ctx := context.TODO()
			r, err := svc.OpenReader(ctx, c.userId, c.subId, c.batchSize)
			assert.ErrorIs(t, err, c.errOpen)
			if c.errOpen == nil {
				var msgs []*pb.CloudEvent
				msgs, err = r.Read()
				assert.Equal(t, c.msgs, msgs)
				assert.ErrorIs(t, err, c.errRead)
				err = r.Close()
				assert.Nil(t, err)
			}
		})
	}
}
