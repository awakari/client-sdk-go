package messages

import (
	"context"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_Read(t *testing.T) {
	svc := NewService(newClientMock())
	cases := map[string]struct {
		userId string
		subId  string
		msgs   []*pb.CloudEvent
		errs   []error
	}{
		"ok": {
			userId: "user0",
			subId:  "sub0",
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
					Id:          "msg0",
					Source:      "source0",
					SpecVersion: "specversion0",
					Type:        "type0",
					Attributes:  map[string]*pb.CloudEventAttributeValue{},
					Data: &pb.CloudEvent_TextData{
						TextData: "data0",
					},
				},
			},
			errs: []error{
				nil,
				nil,
			},
		},
		"fail": {
			userId: "user0",
			subId:  "fail",
			msgs: []*pb.CloudEvent{
				nil,
			},
			errs: []error{
				ErrInternal,
			},
		},
		"fail ack": {
			userId: "user0",
			subId:  "fail_ack",
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
				nil,
			},
			errs: []error{
				nil,
				ErrInternal,
			},
		},
		"fail auth": {
			userId: "user0",
			subId:  "fail_auth",
			msgs: []*pb.CloudEvent{
				nil,
			},
			errs: []error{
				auth.ErrAuth,
			},
		},
		"missing subscription": {
			userId: "user0",
			subId:  "missing",
			msgs: []*pb.CloudEvent{
				nil,
			},
			errs: []error{
				ErrNotFound,
			},
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			ctx := context.TODO()
			rs, err := svc.Read(ctx, c.userId, c.subId)
			require.Nil(t, err)
			defer rs.Close()
			var msg *pb.CloudEvent
			msg, err = rs.Read()
			assert.Equal(t, c.msgs[0], msg)
			assert.ErrorIs(t, err, c.errs[0])
			if err == nil {
				msg, err = rs.Read()
				assert.Equal(t, c.msgs[1], msg)
				assert.ErrorIs(t, err, c.errs[1])
			}
		})
	}
}
