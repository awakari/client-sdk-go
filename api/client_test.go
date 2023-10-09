package api

import (
	"context"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/awakari/client-sdk-go/api/grpc/limits"
	"github.com/awakari/client-sdk-go/api/grpc/permits"
	"github.com/awakari/client-sdk-go/api/grpc/reader"
	"github.com/awakari/client-sdk-go/api/grpc/resolver"
	"github.com/awakari/client-sdk-go/api/grpc/subject"
	"github.com/awakari/client-sdk-go/api/grpc/subscriptions"
	"github.com/awakari/client-sdk-go/model/subscription"
	"github.com/awakari/client-sdk-go/model/subscription/condition"
	"github.com/awakari/client-sdk-go/model/usage"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClient_ReadUsage(t *testing.T) {
	cases := map[string]struct {
		svcPermits permits.Service
		subj       usage.Subject
		userId     string
		u          usage.Usage
		err        error
	}{
		"ok": {
			svcPermits: permits.NewServiceMock(),
			subj:       usage.SubjectSubscriptions,
			userId:     "user0",
			u: usage.Usage{
				Count:      1,
				CountTotal: 2,
				Since:      time.Date(2023, 05, 07, 04, 57, 20, 0, time.UTC),
			},
		},
		"api disabled": {
			subj:   usage.SubjectSubscriptions,
			userId: "user0",
			err:    ErrApiDisabled,
		},
		"invalid subject": {
			svcPermits: permits.NewServiceMock(),
			subj:       usage.SubjectUndefined,
			userId:     "user0",
			err:        subject.ErrInvalidSubject,
		},
		"fail": {
			svcPermits: permits.NewServiceMock(),
			subj:       usage.SubjectPublishEvents,
			userId:     "fail",
			err:        permits.ErrInternal,
		},
		"fail auth": {
			svcPermits: permits.NewServiceMock(),
			subj:       usage.SubjectSubscriptions,
			userId:     "fail_auth",
			err:        auth.ErrAuth,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			cl := client{
				svcPermits: c.svcPermits,
			}
			u, err := cl.ReadUsage(context.TODO(), c.userId, c.subj)
			assert.Equal(t, c.u, u)
			assert.ErrorIs(t, err, c.err)
			assert.Nil(t, cl.Close())
		})
	}
}

func TestClient_ReadUsageLimit(t *testing.T) {
	cases := map[string]struct {
		svcLimits limits.Service
		subj      usage.Subject
		userId    string
		l         usage.Limit
		err       error
	}{
		"ok": {
			svcLimits: limits.NewServiceMock(),
			subj:      usage.SubjectSubscriptions,
			userId:    "user0",
			l: usage.Limit{
				Count:  2,
				UserId: "user0",
			},
		},
		"with expiration": {
			svcLimits: limits.NewServiceMock(),
			subj:      usage.SubjectSubscriptions,
			userId:    "with_expiration",
			l: usage.Limit{
				Count:   2,
				UserId:  "with_expiration",
				Expires: time.Date(2345, 10, 1, 20, 21, 35, 0, time.UTC),
			},
		},
		"both group and user missing": {
			svcLimits: limits.NewServiceMock(),
			subj:      usage.SubjectPublishEvents,
			userId:    "group_missing",
		},
		"group present, user missing": {
			svcLimits: limits.NewServiceMock(),
			subj:      usage.SubjectSubscriptions,
			userId:    "user_missing",
			l: usage.Limit{
				Count: 1,
			},
		},
		"api disabled": {
			subj:   usage.SubjectSubscriptions,
			userId: "user0",
			err:    ErrApiDisabled,
		},
		"invalid subject": {
			svcLimits: limits.NewServiceMock(),
			subj:      usage.SubjectUndefined,
			userId:    "user0",
			err:       subject.ErrInvalidSubject,
		},
		"fail": {
			svcLimits: limits.NewServiceMock(),
			subj:      usage.SubjectPublishEvents,
			userId:    "fail",
			err:       limits.ErrInternal,
		},
		"fail auth": {
			svcLimits: limits.NewServiceMock(),
			subj:      usage.SubjectSubscriptions,
			userId:    "fail_auth",
			err:       auth.ErrAuth,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			cl := client{
				svcLimits: c.svcLimits,
			}
			l, err := cl.ReadUsageLimit(context.TODO(), c.userId, c.subj)
			assert.Equal(t, c.l, l)
			assert.ErrorIs(t, err, c.err)
			assert.Nil(t, cl.Close())
		})
	}
}

func TestClient_ReadMessages(t *testing.T) {
	cases := map[string]struct {
		svcReader reader.Service
		subId     string
		batchSize uint32
		msgs      []*pb.CloudEvent
		err0      error
		err1      error
	}{
		"ok": {
			svcReader: reader.NewServiceMock(),
			subId:     "sub0",
			batchSize: 3,
			msgs: []*pb.CloudEvent{
				{
					Id:          "msg0",
					Source:      "source0",
					SpecVersion: "specversion0",
					Type:        "type0",
					Attributes:  map[string]*pb.CloudEventAttributeValue{},
					Data: &pb.CloudEvent_TextData{
						TextData: "data",
					},
				},
				{
					Id:          "msg1",
					Source:      "source0",
					SpecVersion: "specversion0",
					Type:        "type0",
					Attributes:  map[string]*pb.CloudEventAttributeValue{},
					Data: &pb.CloudEvent_TextData{
						TextData: "data",
					},
				},
				{
					Id:          "msg2",
					Source:      "source0",
					SpecVersion: "specversion0",
					Type:        "type0",
					Attributes:  map[string]*pb.CloudEventAttributeValue{},
					Data: &pb.CloudEvent_TextData{
						TextData: "data",
					},
				},
			},
		},
		"api disabled": {
			subId:     "sub0",
			batchSize: 3,
			err0:      ErrApiDisabled,
		},
		"fail": {
			svcReader: reader.NewServiceMock(),
			subId:     "fail",
			batchSize: 3,
			err0:      reader.ErrInternal,
		},
		"fail auth": {
			svcReader: reader.NewServiceMock(),
			subId:     "fail_auth",
			batchSize: 3,
			err0:      auth.ErrAuth,
		},
		"fail read": {
			svcReader: reader.NewServiceMock(),
			subId:     "fail_read",
			batchSize: 3,
			err1:      reader.ErrInternal,
		},
		"sub missing": {
			svcReader: reader.NewServiceMock(),
			subId:     "missing",
			batchSize: 3,
			err0:      reader.ErrNotFound,
		},
		"invalid batch size": {
			svcReader: reader.NewServiceMock(),
			subId:     "sub0",
			batchSize: 0,
			err1:      reader.ErrInvalidRequest,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			cl := client{
				svcReader: c.svcReader,
			}
			rs, err := cl.OpenMessagesReader(context.TODO(), "user0", c.subId, c.batchSize)
			assert.ErrorIs(t, err, c.err0)
			if err == nil {
				var msgs []*pb.CloudEvent
				msgs, err = rs.Read()
				assert.Equal(t, c.msgs, msgs)
				assert.ErrorIs(t, err, c.err1)
				assert.Nil(t, rs.Close())
			}
			assert.Nil(t, cl.Close())
		})
	}
}

func TestClient_ReadMessages_Ack(t *testing.T) {
	cases := map[string]struct {
		svcReader reader.Service
		subId     string
		batchSize uint32
		msgs      []*pb.CloudEvent
		err0      error
		err1      error
	}{
		"ok": {
			svcReader: reader.NewServiceMock(),
			subId:     "sub0",
			batchSize: 3,
			msgs: []*pb.CloudEvent{
				{
					Id:          "msg0",
					Source:      "source0",
					SpecVersion: "specversion0",
					Type:        "type0",
					Attributes:  map[string]*pb.CloudEventAttributeValue{},
					Data: &pb.CloudEvent_TextData{
						TextData: "data",
					},
				},
				{
					Id:          "msg1",
					Source:      "source0",
					SpecVersion: "specversion0",
					Type:        "type0",
					Attributes:  map[string]*pb.CloudEventAttributeValue{},
					Data: &pb.CloudEvent_TextData{
						TextData: "data",
					},
				},
				{
					Id:          "msg2",
					Source:      "source0",
					SpecVersion: "specversion0",
					Type:        "type0",
					Attributes:  map[string]*pb.CloudEventAttributeValue{},
					Data: &pb.CloudEvent_TextData{
						TextData: "data",
					},
				},
			},
		},
		"api disabled": {
			subId:     "sub0",
			batchSize: 3,
			err0:      ErrApiDisabled,
		},
		"fail": {
			svcReader: reader.NewServiceMock(),
			subId:     "fail",
			batchSize: 3,
			err0:      reader.ErrInternal,
		},
		"fail auth": {
			svcReader: reader.NewServiceMock(),
			subId:     "fail_auth",
			batchSize: 3,
			err0:      auth.ErrAuth,
		},
		"fail read": {
			svcReader: reader.NewServiceMock(),
			subId:     "fail_read",
			batchSize: 3,
			err1:      reader.ErrInternal,
		},
		"sub missing": {
			svcReader: reader.NewServiceMock(),
			subId:     "missing",
			batchSize: 3,
			err0:      reader.ErrNotFound,
		},
		"invalid batch size": {
			svcReader: reader.NewServiceMock(),
			subId:     "sub0",
			batchSize: 0,
			err1:      reader.ErrInvalidRequest,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			cl := client{
				svcReader: c.svcReader,
			}
			rs, err := cl.OpenMessagesAckReader(context.TODO(), "user0", c.subId, c.batchSize)
			assert.ErrorIs(t, err, c.err0)
			if err == nil {
				var msgs []*pb.CloudEvent
				msgs, err = rs.Read()
				assert.Equal(t, c.msgs, msgs)
				assert.ErrorIs(t, err, c.err1)
				assert.Nil(t, rs.Ack(uint32(len(msgs))))
				assert.Nil(t, rs.Close())
			}
			assert.Nil(t, cl.Close())
		})
	}
}

func TestClient_WriteMessages(t *testing.T) {
	cases := map[string]struct {
		svcWriter resolver.Service
		userId    string
		err0      error
		msgs      []*pb.CloudEvent
		ackCount  uint32
		err1      error
	}{
		"api disabled": {
			userId: "user0",
			err0:   ErrApiDisabled,
		},
		"fail open stream": {
			svcWriter: resolver.NewServiceMock(),
			userId:    "fail",
			err0:      resolver.ErrInternal,
		},
		"fail auth": {
			svcWriter: resolver.NewServiceMock(),
			userId:    "fail_auth",
			err0:      auth.ErrAuth,
		},
		"fail write": {
			svcWriter: resolver.NewServiceMock(),
			userId:    "user0",
			msgs: []*pb.CloudEvent{
				{
					Id: "msg0",
				},
				{
					Id: "fail",
				},
			},
			ackCount: 1,
			err1:     resolver.ErrInternal,
		},
		"limit reached": {
			svcWriter: resolver.NewServiceMock(),
			userId:    "user0",
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
			err1:     limits.ErrReached,
		},
		"ok": {
			svcWriter: resolver.NewServiceMock(),
			userId:    "user0",
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
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			cl := client{
				svcWriter: c.svcWriter,
			}
			ws, err := cl.OpenMessagesWriter(context.TODO(), c.userId)
			assert.ErrorIs(t, err, c.err0)
			if err == nil {
				var ackCount uint32
				ackCount, err = ws.WriteBatch(c.msgs)
				assert.Equal(t, c.ackCount, ackCount)
				assert.ErrorIs(t, err, c.err1)
				assert.Nil(t, ws.Close())
			}
			assert.Nil(t, cl.Close())
		})
	}
}

func TestClient_CreateSubscription(t *testing.T) {
	cases := map[string]struct {
		svcSubs subscriptions.Service
		descr   string
		id      string
		err     error
	}{
		"ok": {
			svcSubs: subscriptions.NewServiceMock(),
			descr:   "my subscription",
			id:      "sub0",
		},
		"subs API not set": {
			descr: "my subscription",
			err:   ErrApiDisabled,
		},
		"fail": {
			svcSubs: subscriptions.NewServiceMock(),
			descr:   "fail",
			err:     subscriptions.ErrInternal,
		},
		"fail auth": {
			svcSubs: subscriptions.NewServiceMock(),
			descr:   "fail_auth",
			err:     auth.ErrAuth,
		},
		"invalid": {
			svcSubs: subscriptions.NewServiceMock(),
			descr:   "invalid",
			err:     subscriptions.ErrInvalid,
		},
		"limit reached": {
			svcSubs: subscriptions.NewServiceMock(),
			descr:   "limit_reached",
			err:     limits.ErrReached,
		},
		"busy": {
			svcSubs: subscriptions.NewServiceMock(),
			descr:   "busy",
			err:     subscriptions.ErrBusy,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			cl := client{
				svcSubs: c.svcSubs,
			}
			ctx := context.TODO()
			subData := subscription.Data{
				Description: c.descr,
				Expires:     time.Now(),
			}
			id, err := cl.CreateSubscription(ctx, "user0", subData)
			assert.Equal(t, c.id, id)
			assert.ErrorIs(t, err, c.err)
			assert.Nil(t, cl.Close())
		})
	}
}

func TestClient_ReadSubscription(t *testing.T) {
	cases := map[string]struct {
		svcSubs subscriptions.Service
		subId   string
		subData subscription.Data
		err     error
	}{
		"ok": {
			svcSubs: subscriptions.NewServiceMock(),
			subId:   "sub0",
			subData: subscription.Data{
				Description: "my subscription",
				Enabled:     true,
				Expires:     time.Date(2023, 10, 4, 11, 44, 55, 0, time.UTC),
				Condition: condition.
					NewBuilder().
					Any([]condition.Condition{
						condition.
							NewBuilder().BuildTextCondition(),
						condition.
							NewBuilder().
							LessThanOrEqual(42).
							BuildNumberCondition(),
					}).
					BuildGroupCondition(),
			},
		},
		"subs API not set": {
			subId: "sub0",
			err:   ErrApiDisabled,
		},
		"fail": {
			svcSubs: subscriptions.NewServiceMock(),
			subId:   "fail",
			err:     subscriptions.ErrInternal,
		},
		"fail auth": {
			svcSubs: subscriptions.NewServiceMock(),
			subId:   "fail_auth",
			err:     auth.ErrAuth,
		},
		"missing": {
			svcSubs: subscriptions.NewServiceMock(),
			subId:   "missing",
			err:     subscriptions.ErrNotFound,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			cl := client{
				svcSubs: c.svcSubs,
			}
			ctx := context.TODO()
			subData, err := cl.ReadSubscription(ctx, "user0", c.subId)
			assert.Equal(t, c.subData, subData)
			assert.ErrorIs(t, err, c.err)
			assert.Nil(t, cl.Close())
		})
	}
}

func TestClient_UpdateSubscriptionMetadata(t *testing.T) {
	cases := map[string]struct {
		svcSubs subscriptions.Service
		subId   string
		err     error
	}{
		"ok": {
			svcSubs: subscriptions.NewServiceMock(),
			subId:   "sub0",
		},
		"subs API not set": {
			subId: "sub0",
			err:   ErrApiDisabled,
		},
		"fail": {
			svcSubs: subscriptions.NewServiceMock(),
			subId:   "fail",
			err:     subscriptions.ErrInternal,
		},
		"fail auth": {
			svcSubs: subscriptions.NewServiceMock(),
			subId:   "fail_auth",
			err:     auth.ErrAuth,
		},
		"missing": {
			svcSubs: subscriptions.NewServiceMock(),
			subId:   "missing",
			err:     subscriptions.ErrNotFound,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			cl := client{
				svcSubs: c.svcSubs,
			}
			ctx := context.TODO()
			err := cl.UpdateSubscription(ctx, "user0", c.subId, subscription.Data{
				Expires: time.Now(),
			})
			assert.ErrorIs(t, err, c.err)
			assert.Nil(t, cl.Close())
		})
	}
}

func TestClient_DeleteSubscription(t *testing.T) {
	cases := map[string]struct {
		svcSubs subscriptions.Service
		subId   string
		err     error
	}{
		"ok": {
			svcSubs: subscriptions.NewServiceMock(),
			subId:   "sub0",
		},
		"subs API not set": {
			subId: "sub0",
			err:   ErrApiDisabled,
		},
		"fail": {
			svcSubs: subscriptions.NewServiceMock(),
			subId:   "fail",
			err:     subscriptions.ErrInternal,
		},
		"fail auth": {
			svcSubs: subscriptions.NewServiceMock(),
			subId:   "fail_auth",
			err:     auth.ErrAuth,
		},
		"missing": {
			svcSubs: subscriptions.NewServiceMock(),
			subId:   "missing",
			err:     subscriptions.ErrNotFound,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			cl := client{
				svcSubs: c.svcSubs,
			}
			ctx := context.TODO()
			err := cl.DeleteSubscription(ctx, "user0", c.subId)
			assert.ErrorIs(t, err, c.err)
			assert.Nil(t, cl.Close())
		})
	}
}

func TestClient_SearchSubscriptions(t *testing.T) {
	cases := map[string]struct {
		svcSubs subscriptions.Service
		ids     []string
		cursor  string
		err     error
	}{
		"ok0": {
			svcSubs: subscriptions.NewServiceMock(),
			ids: []string{
				"sub0",
				"sub1",
			},
		},
		"ok1": {
			svcSubs: subscriptions.NewServiceMock(),
			cursor:  "sub1",
		},
		"subs API not set": {
			err: ErrApiDisabled,
		},
		"fail": {
			svcSubs: subscriptions.NewServiceMock(),
			cursor:  "fail",
			err:     subscriptions.ErrInternal,
		},
		"fail auth": {
			svcSubs: subscriptions.NewServiceMock(),
			cursor:  "fail_auth",
			err:     auth.ErrAuth,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			cl := client{
				svcSubs: c.svcSubs,
			}
			ctx := context.TODO()
			ids, err := cl.SearchSubscriptions(ctx, "user0", 0, c.cursor)
			assert.Equal(t, c.ids, ids)
			assert.ErrorIs(t, err, c.err)
			assert.Nil(t, cl.Close())
		})
	}
}
