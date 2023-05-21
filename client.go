package client_sdk_go

import (
	"context"
	"errors"
	"fmt"
	"github.com/awakari/client-sdk-go/api/grpc/limits"
	"github.com/awakari/client-sdk-go/api/grpc/messages"
	"github.com/awakari/client-sdk-go/api/grpc/permits"
	"github.com/awakari/client-sdk-go/api/grpc/subscriptions"
	"github.com/awakari/client-sdk-go/api/grpc/writer"
	"github.com/awakari/client-sdk-go/model"
	"github.com/awakari/client-sdk-go/model/subscription"
	"github.com/awakari/client-sdk-go/model/usage"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"google.golang.org/grpc"
	"io"
)

type Client interface {
	io.Closer

	// Usage

	// ReadUsage returns the current usage counts by the requested subject.
	ReadUsage(ctx context.Context, userId string, subj usage.Subject) (u usage.Usage, err error)

	// ReadUsageLimit returns the usage limit by the requested subject.
	ReadUsageLimit(ctx context.Context, userId string, subj usage.Subject) (l usage.Limit, err error)

	// Messages

	// WriteMessages opens the stream for publishing the messages.
	WriteMessages(ctx context.Context, userId string) (ws model.WriteStream[*pb.CloudEvent], err error)

	// ReadMessages opens the stream for receiving the messages matching the requested subscription.
	ReadMessages(ctx context.Context, userId, subId string) (rs model.ReadStream[*pb.CloudEvent], err error)

	// Subscriptions

	// CreateSubscription with the specified fields.
	CreateSubscription(ctx context.Context, userId string, subData subscription.Data) (id string, err error)

	// ReadSubscription specified by the id. Returns ErrNotFound if subscription is missing.
	ReadSubscription(ctx context.Context, userId, subId string) (subData subscription.Data, err error)

	// UpdateSubscriptionMetadata updates the mutable part of the subscription.Data
	UpdateSubscriptionMetadata(ctx context.Context, userId, subId string, md subscription.Metadata) (err error)

	// DeleteSubscription and all associated conditions those not in use by any other subscription.
	// Returns ErrNotFound if a subscription with the specified id is missing.
	DeleteSubscription(ctx context.Context, userId, subId string) (err error)

	// SearchSubscriptions returns all subscription ids those have the requested user id.
	SearchSubscriptions(ctx context.Context, userId string, limit uint32, cursor string) (ids []string, err error)
}

type client struct {
	connLimits  *grpc.ClientConn
	connMsgs    *grpc.ClientConn
	connPermits *grpc.ClientConn
	connSubs    *grpc.ClientConn
	connWriter  *grpc.ClientConn
	svcLimits   limits.Service
	svcMsgs     messages.Service
	svcPermits  permits.Service
	svcSubs     subscriptions.Service
	svcWriter   writer.Service
}

var ErrApiDisabled = errors.New("the API call is not enabled for this client")

var _ Client = (*client)(nil)

func (c client) Close() (err error) {
	if c.connLimits != nil {
		err = errors.Join(err, c.connLimits.Close())
	}
	if c.connMsgs != nil {
		err = errors.Join(err, c.connMsgs.Close())
	}
	if c.connPermits != nil {
		err = errors.Join(err, c.connPermits.Close())
	}
	if c.connSubs != nil {
		err = errors.Join(err, c.connSubs.Close())
	}
	if c.connWriter != nil {
		err = errors.Join(err, c.connWriter.Close())
	}
	return
}

func (c client) ReadUsage(ctx context.Context, userId string, subj usage.Subject) (u usage.Usage, err error) {
	if c.svcPermits == nil {
		err = fmt.Errorf("%w: ReadUsage(...)", ErrApiDisabled)
	} else {
		u, err = c.svcPermits.GetUsage(ctx, userId, subj)
	}
	return
}

func (c client) ReadUsageLimit(ctx context.Context, userId string, subj usage.Subject) (l usage.Limit, err error) {
	if c.svcLimits == nil {
		err = fmt.Errorf("%w: ReadUsageLimit(...)", ErrApiDisabled)
	} else {
		l, err = c.svcLimits.Get(ctx, userId, subj)
	}
	return
}

func (c client) WriteMessages(ctx context.Context, userId string) (ws model.WriteStream[*pb.CloudEvent], err error) {
	if c.svcWriter == nil {
		err = fmt.Errorf("%w: WriteMessages(...)", ErrApiDisabled)
	} else {
		ws, err = c.svcWriter.OpenStream(ctx, userId)
	}
	return
}

func (c client) ReadMessages(ctx context.Context, userId, subId string) (rs model.ReadStream[*pb.CloudEvent], err error) {
	if c.svcMsgs == nil {
		err = fmt.Errorf("%w: ReadMessages(...)", ErrApiDisabled)
	} else {
		rs, err = c.svcMsgs.Read(ctx, userId, subId)
	}
	return
}

func (c client) CreateSubscription(ctx context.Context, userId string, subData subscription.Data) (id string, err error) {
	if c.svcSubs == nil {
		err = fmt.Errorf("%w: CreateSubscription(...)", ErrApiDisabled)
	} else {
		id, err = c.svcSubs.Create(ctx, userId, subData)
	}
	return
}

func (c client) ReadSubscription(ctx context.Context, userId, subId string) (subData subscription.Data, err error) {
	if c.svcSubs == nil {
		err = fmt.Errorf("%w: ReadSubscription(...)", ErrApiDisabled)
	} else {
		subData, err = c.svcSubs.Read(ctx, userId, subId)
	}
	return
}

func (c client) UpdateSubscriptionMetadata(ctx context.Context, userId, subId string, md subscription.Metadata) (err error) {
	if c.svcSubs == nil {
		err = fmt.Errorf("%w: UpdateSubscriptionMetadata(...)", ErrApiDisabled)
	} else {
		err = c.svcSubs.UpdateMetadata(ctx, userId, subId, md)
	}
	return
}

func (c client) DeleteSubscription(ctx context.Context, userId, subId string) (err error) {
	if c.svcSubs == nil {
		err = fmt.Errorf("%w: DeleteSubscription(...)", ErrApiDisabled)
	} else {
		err = c.svcSubs.Delete(ctx, userId, subId)
	}
	return
}

func (c client) SearchSubscriptions(ctx context.Context, userId string, limit uint32, cursor string) (ids []string, err error) {
	if c.svcSubs == nil {
		err = fmt.Errorf("%w: SearchSubscriptions(...)", ErrApiDisabled)
	} else {
		ids, err = c.svcSubs.SearchOwn(ctx, userId, limit, cursor)
	}
	return
}
