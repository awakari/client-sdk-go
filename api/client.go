package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/awakari/client-sdk-go/api/grpc/limits"
	"github.com/awakari/client-sdk-go/api/grpc/permits"
	"github.com/awakari/client-sdk-go/api/grpc/reader"
	"github.com/awakari/client-sdk-go/api/grpc/resolver"
	"github.com/awakari/client-sdk-go/api/grpc/subscriptions"
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

	// OpenMessagesWriter opens the batch message writer. A client should close it once done.
	OpenMessagesWriter(ctx context.Context, userId string) (w model.Writer[*pb.CloudEvent], err error)

	// OpenMessagesReader opens batch message reader. A client should close it once done.
	OpenMessagesReader(ctx context.Context, userId, subId string, batchSize uint32) (r model.Reader[[]*pb.CloudEvent], err error)

	// OpenMessagesAckReader opens batch message reader that requires an explicit ack. A client should close it once done.
	OpenMessagesAckReader(ctx context.Context, userId, subId string, batchSize uint32) (r model.AckReader[[]*pb.CloudEvent], err error)

	// Subscriptions

	// CreateSubscription with the specified fields.
	CreateSubscription(ctx context.Context, userId string, subData subscription.Data) (id string, err error)

	// ReadSubscription specified by the id. Returns ErrNotFound if subscription is missing.
	ReadSubscription(ctx context.Context, userId, subId string) (subData subscription.Data, err error)

	// UpdateSubscription replaces the existing subscription.Data fields.
	UpdateSubscription(ctx context.Context, userId, subId string, subData subscription.Data) (err error)

	// DeleteSubscription and all associated conditions those not in use by any other subscription.
	// Returns ErrNotFound if a subscription with the specified id is missing.
	DeleteSubscription(ctx context.Context, userId, subId string) (err error)

	// SearchSubscriptions returns all subscription ids those have the requested user id.
	SearchSubscriptions(ctx context.Context, userId string, q subscription.Query, cursor subscription.Cursor) (ids []string, err error)
}

type client struct {
	connLimits  *grpc.ClientConn
	connReader  *grpc.ClientConn
	connPermits *grpc.ClientConn
	connSubs    *grpc.ClientConn
	connWriter  *grpc.ClientConn
	svcLimits   limits.Service
	svcReader   reader.Service
	svcPermits  permits.Service
	svcSubs     subscriptions.Service
	svcWriter   resolver.Service
}

var ErrApiDisabled = errors.New("the API call is not enabled for this client")

func (c client) Close() (err error) {
	if c.connLimits != nil {
		err = errors.Join(err, c.connLimits.Close())
	}
	if c.connReader != nil {
		err = errors.Join(err, c.connReader.Close())
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

func (c client) OpenMessagesWriter(ctx context.Context, userId string) (ws model.Writer[*pb.CloudEvent], err error) {
	if c.svcWriter == nil {
		err = fmt.Errorf("%w: OpenMessagesWriter(...)", ErrApiDisabled)
	} else {
		ws, err = c.svcWriter.OpenWriter(ctx, userId)
	}
	return
}

func (c client) OpenMessagesReader(ctx context.Context, userId, subId string, batchSize uint32) (rs model.Reader[[]*pb.CloudEvent], err error) {
	if c.svcReader == nil {
		err = fmt.Errorf("%w: OpenMessagesReader(...)", ErrApiDisabled)
	} else {
		rs, err = c.svcReader.OpenReader(ctx, userId, subId, batchSize)
	}
	return
}

func (c client) OpenMessagesAckReader(ctx context.Context, userId, subId string, batchSize uint32) (r model.AckReader[[]*pb.CloudEvent], err error) {
	if c.svcReader == nil {
		err = fmt.Errorf("%w: OpenMessagesAckReader(...)", ErrApiDisabled)
	} else {
		r, err = c.svcReader.OpenAckReader(ctx, userId, subId, batchSize)
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

func (c client) UpdateSubscription(ctx context.Context, userId, subId string, subData subscription.Data) (err error) {
	if c.svcSubs == nil {
		err = fmt.Errorf("%w: UpdateSubscription(...)", ErrApiDisabled)
	} else {
		err = c.svcSubs.Update(ctx, userId, subId, subData)
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

func (c client) SearchSubscriptions(ctx context.Context, userId string, q subscription.Query, cursor subscription.Cursor) (ids []string, err error) {
	if c.svcSubs == nil {
		err = fmt.Errorf("%w: SearchSubscriptions(...)", ErrApiDisabled)
	} else {
		ids, err = c.svcSubs.Search(ctx, userId, q, cursor)
	}
	return
}
