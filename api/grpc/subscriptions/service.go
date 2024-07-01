package subscriptions

import (
	"context"
	"errors"
	"fmt"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/awakari/client-sdk-go/api/grpc/limits"
	"github.com/awakari/client-sdk-go/model/subscription"
	"github.com/awakari/client-sdk-go/model/subscription/condition"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service interface {

	// Create a subscription with the specified fields.
	Create(ctx context.Context, userId string, subData subscription.Data) (id string, err error)

	// Read returns the subscription specified by the id. Returns ErrNotFound if subscription is missing.
	Read(ctx context.Context, userId, subId string) (subData subscription.Data, err error)

	// Update the subscription.Data
	Update(ctx context.Context, userId, subId string, subData subscription.Data) (err error)

	// Delete the specified subscription all associated conditions those not in use by any other subscription.
	// Returns ErrNotFound if a subscription with the specified id is missing.
	Delete(ctx context.Context, userId, subId string) (err error)

	// Search returns all subscription ids matching the query.
	Search(ctx context.Context, userId string, q subscription.Query, cursor subscription.Cursor) (ids []string, err error)
}

type service struct {
	client ServiceClient
}

// ErrNotFound indicates the subscription is missing in the storage and can not be read/updated/deleted.
var ErrNotFound = errors.New("subscription was not found")

// ErrInternal indicates some unexpected internal failure.
var ErrInternal = errors.New("internal failure")

// ErrInvalid indicates the invalid subscription.
var ErrInvalid = errors.New("invalid subscription")

var ErrUnavailable = errors.New("unavailable")

func NewService(client ServiceClient) Service {
	return service{
		client: client,
	}
}

func (svc service) Create(ctx context.Context, userId string, subData subscription.Data) (id string, err error) {
	ctx = auth.SetOutgoingAuthInfo(ctx, userId)
	req := CreateRequest{
		Cond:        encodeCondition(subData.Condition),
		Description: subData.Description,
		Enabled:     subData.Enabled,
		Expires:     timestamppb.New(subData.Expires.UTC()),
		Public:      subData.Public,
	}
	var resp *CreateResponse
	resp, err = svc.client.Create(ctx, &req)
	err = decodeError(err)
	if err == nil {
		id = resp.Id
	}
	return
}

func (svc service) Read(ctx context.Context, userId, subId string) (subData subscription.Data, err error) {
	ctx = auth.SetOutgoingAuthInfo(ctx, userId)
	req := ReadRequest{
		Id: subId,
	}
	var resp *ReadResponse
	resp, err = svc.client.Read(ctx, &req)
	err = decodeError(err)
	if err == nil {
		subData.Condition, err = decodeCondition(resp.Cond)
		subData.Description = resp.Description
		subData.Enabled = resp.Enabled
		if resp.Expires != nil {
			subData.Expires = resp.Expires.AsTime()
		}
		if resp.Created != nil {
			subData.Created = resp.Created.AsTime()
		}
		if resp.Updated != nil {
			subData.Created = resp.Updated.AsTime()
		}
		subData.Public = resp.Public
		subData.Followers = resp.Followers
	}
	return
}

func (svc service) Update(ctx context.Context, userId, subId string, data subscription.Data) (err error) {
	ctx = auth.SetOutgoingAuthInfo(ctx, userId)
	if data.Condition == nil {
		err = fmt.Errorf("%w: missing condition", ErrInvalid)
	}
	if err == nil {
		req := UpdateRequest{
			Id:          subId,
			Description: data.Description,
			Enabled:     data.Enabled,
			Expires:     timestamppb.New(data.Expires.UTC()),
			Cond:        encodeCondition(data.Condition),
			Public:      data.Public,
		}
		_, err = svc.client.Update(ctx, &req)
		err = decodeError(err)
	}
	return
}

func (svc service) Delete(ctx context.Context, userId, subId string) (err error) {
	ctx = auth.SetOutgoingAuthInfo(ctx, userId)
	req := DeleteRequest{
		Id: subId,
	}
	_, err = svc.client.Delete(ctx, &req)
	err = decodeError(err)
	return
}

func (svc service) Search(ctx context.Context, userId string, q subscription.Query, cursor subscription.Cursor) (ids []string, err error) {
	ctx = auth.SetOutgoingAuthInfo(ctx, userId)
	switch q.Public {
	case true:
		req := SearchRequest{
			Cursor: &Cursor{
				Followers: cursor.Followers,
				Id:        cursor.Id,
			},
			Limit:   q.Limit,
			Order:   Order(q.Order),
			Pattern: q.Pattern,
			Sort:    Sort(q.Sort),
		}
		var resp *SearchResponse
		resp, err = svc.client.Search(ctx, &req)
		if resp != nil {
			ids = resp.Ids
		}
	default:
		req := SearchOwnRequest{
			Cursor:  cursor.Id,
			Limit:   q.Limit,
			Order:   Order(q.Order),
			Pattern: q.Pattern,
		}
		var resp *SearchOwnResponse
		resp, err = svc.client.SearchOwn(ctx, &req)
		if resp != nil {
			ids = resp.Ids
		}
	}
	err = decodeError(err)
	return
}

func encodeCondition(src condition.Condition) (dst *Condition) {
	dst = &Condition{
		Not: src.IsNot(),
	}
	switch c := src.(type) {
	case condition.GroupCondition:
		var dstGroup []*Condition
		for _, childSrc := range c.GetGroup() {
			childDst := encodeCondition(childSrc)
			dstGroup = append(dstGroup, childDst)
		}
		dst.Cond = &Condition_Gc{
			Gc: &GroupCondition{
				Logic: GroupLogic(c.GetLogic()),
				Group: dstGroup,
			},
		}
	case condition.TextCondition:
		dst.Cond = &Condition_Tc{
			Tc: &TextCondition{
				Key:   c.GetKey(),
				Term:  c.GetTerm(),
				Exact: c.IsExact(),
			},
		}
	case condition.NumberCondition:
		dstOp := encodeNumOp(c.GetOperation())
		dst.Cond = &Condition_Nc{
			Nc: &NumberCondition{
				Key: c.GetKey(),
				Op:  dstOp,
				Val: c.GetValue(),
			},
		}
	}
	return
}

func encodeNumOp(src condition.NumOp) (dst Operation) {
	switch src {
	case condition.NumOpGt:
		dst = Operation_Gt
	case condition.NumOpGte:
		dst = Operation_Gte
	case condition.NumOpEq:
		dst = Operation_Eq
	case condition.NumOpLte:
		dst = Operation_Lte
	case condition.NumOpLt:
		dst = Operation_Lt
	default:
		dst = Operation_Undefined
	}
	return
}

func decodeCondition(src *Condition) (dst condition.Condition, err error) {
	gc, nc, tc := src.GetGc(), src.GetNc(), src.GetTc()
	switch {
	case gc != nil:
		var group []condition.Condition
		var childDst condition.Condition
		for _, childSrc := range gc.Group {
			childDst, err = decodeCondition(childSrc)
			if err != nil {
				break
			}
			group = append(group, childDst)
		}
		if err == nil {
			dst = condition.NewGroupCondition(
				condition.NewCondition(src.Not),
				condition.GroupLogic(gc.GetLogic()),
				group,
			)
		}
	case nc != nil:
		dstOp := decodeNumOp(nc.Op)
		dst = condition.NewNumberCondition(
			condition.NewKeyCondition(condition.NewCondition(src.Not), nc.GetKey()),
			dstOp, nc.Val,
		)
	case tc != nil:
		dst = condition.NewTextCondition(
			condition.NewKeyCondition(condition.NewCondition(src.Not), tc.GetKey()),
			tc.GetTerm(), tc.GetExact(),
		)
	default:
		err = fmt.Errorf("%w: unsupported condition type", ErrInternal)
	}
	return
}

func decodeNumOp(src Operation) (dst condition.NumOp) {
	switch src {
	case Operation_Gt:
		dst = condition.NumOpGt
	case Operation_Gte:
		dst = condition.NumOpGte
	case Operation_Eq:
		dst = condition.NumOpEq
	case Operation_Lte:
		dst = condition.NumOpLte
	case Operation_Lt:
		dst = condition.NumOpLt
	default:
		dst = condition.NumOpUndefined
	}
	return
}

func decodeError(src error) (dst error) {
	switch {
	case src == nil:
	default:
		s, isGrpcErr := status.FromError(src)
		switch {
		case !isGrpcErr:
			dst = src
		case s.Code() == codes.OK:
		case s.Code() == codes.NotFound:
			dst = fmt.Errorf("%w: %s", ErrNotFound, src)
		case s.Code() == codes.InvalidArgument:
			dst = fmt.Errorf("%w: %s", ErrInvalid, src)
		case s.Code() == codes.ResourceExhausted:
			dst = fmt.Errorf("%w: %s", limits.ErrReached, src)
		case s.Code() == codes.Unauthenticated:
			dst = fmt.Errorf("%w: %s", auth.ErrAuth, src)
		case s.Code() == codes.Unavailable:
			dst = fmt.Errorf("%w: %s", ErrUnavailable, src)
		default:
			dst = fmt.Errorf("%w: %s", ErrInternal, src)
		}
	}
	return
}
