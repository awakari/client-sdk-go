package subscriptions

import (
	"context"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/awakari/client-sdk-go/api/grpc/limits"
	"github.com/awakari/client-sdk-go/model/subscription"
	"github.com/awakari/client-sdk-go/model/subscription/condition"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestService_Create(t *testing.T) {
	//
	svc := NewService(newClientMock())
	//
	cases := map[string]struct {
		req subscription.Data
		err error
	}{
		"invalid": {
			req: subscription.Data{
				Metadata: subscription.Metadata{
					Description: "invalid",
				},
				Condition: condition.NewKiwiTreeCondition(
					condition.NewKiwiCondition(
						condition.NewKeyCondition(condition.NewCondition(false), "key0"),
						false,
						"ok",
					),
				),
			},
			err: ErrInvalid,
		},
		"locked": {
			req: subscription.Data{
				Metadata: subscription.Metadata{
					Description: "busy",
				},
				Condition: condition.NewKiwiTreeCondition(
					condition.NewKiwiCondition(
						condition.NewKeyCondition(condition.NewCondition(false), ""),
						false,
						"locked",
					),
				),
			},
			err: ErrBusy,
		},
		"fail": {
			req: subscription.Data{
				Metadata: subscription.Metadata{
					Description: "fail",
				},
				Condition: condition.NewKiwiTreeCondition(
					condition.NewKiwiCondition(
						condition.NewKeyCondition(
							condition.NewCondition(false),
							"fail",
						),
						false,
						"fail",
					),
				),
			},
			err: ErrInternal,
		},
		"ok": {
			req: subscription.Data{
				Metadata: subscription.Metadata{
					Description: "my subscription",
				},
				Condition: condition.NewKiwiTreeCondition(
					condition.NewKiwiCondition(
						condition.NewKeyCondition(condition.NewCondition(false), "key0"),
						false,
						"ok",
					),
				),
			},
		},
		"ok group": {
			req: subscription.Data{
				Metadata: subscription.Metadata{
					Description: "my subscription",
				},
				Condition: condition.
					NewBuilder().
					GroupChildren(
						[]condition.Condition{
							condition.
								NewBuilder().
								BuildKiwiTreeCondition(),
						},
					).
					BuildGroupCondition(),
			},
		},
		"auth fail": {
			req: subscription.Data{
				Metadata: subscription.Metadata{
					Description: "fail_auth",
				},
				Condition: condition.NewKiwiTreeCondition(
					condition.NewKiwiCondition(
						condition.NewKeyCondition(condition.NewCondition(false), "key0"),
						false,
						"ok",
					),
				),
			},
			err: auth.ErrAuth,
		},
		"limit reached": {
			req: subscription.Data{
				Metadata: subscription.Metadata{
					Description: "limit_reached",
				},
				Condition: condition.NewKiwiTreeCondition(
					condition.NewKiwiCondition(
						condition.NewKeyCondition(condition.NewCondition(false), "key0"),
						false,
						"ok",
					),
				),
			},
			err: limits.ErrReached,
		},
	}
	//
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			id, err := svc.Create(ctx, "user0", c.req)
			if c.err == nil {
				assert.Nil(t, err)
				assert.NotEmpty(t, id)
			} else {
				assert.ErrorIs(t, err, c.err)
			}
		})
	}
}

func TestService_Read(t *testing.T) {
	//
	svc := NewService(newClientMock())
	//
	cases := map[string]struct {
		sd  subscription.Data
		err error
	}{
		"missing": {
			err: ErrNotFound,
		},
		"fail": {
			err: ErrInternal,
		},
		"fail_auth": {
			err: auth.ErrAuth,
		},
		"ok": {
			sd: subscription.Data{
				Metadata: subscription.Metadata{
					Description: "subscription",
					Enabled:     true,
				},
				Condition: condition.
					NewBuilder().
					GroupLogic(condition.GroupLogicOr).
					GroupChildren(
						[]condition.Condition{
							condition.
								NewBuilder().
								Negation().
								MatchAttrKey("k0").
								MatchAttrValuePattern("p0").
								BuildKiwiTreeCondition(),
							condition.
								NewBuilder().
								MatchAttrKey("k1").
								MatchAttrValuePattern("p1").
								MatchAttrValuePartial().
								BuildKiwiTreeCondition(),
						},
					).
					BuildGroupCondition(),
			},
		},
	}
	//
	for id, c := range cases {
		t.Run(id, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			sd, err := svc.Read(ctx, "user0", id)
			if c.err == nil {
				assert.Nil(t, err)
				assert.Equal(t, c.sd.Metadata, sd.Metadata)
				assert.True(t, conditionsDataEqual(c.sd.Condition, sd.Condition))
			} else {
				assert.ErrorIs(t, err, c.err)
			}
		})
	}
}

func TestService_UpdateMetadata(t *testing.T) {
	//
	svc := NewService(newClientMock())
	//
	cases := map[string]struct {
		id  string
		err error
		md  subscription.Metadata
	}{
		"not found": {
			id:  "missing",
			err: ErrNotFound,
		},
		"ok": {
			id: "sub0",
			md: subscription.Metadata{
				Description: "my subscription",
				Enabled:     false,
			},
		},
		"fail": {
			id:  "fail",
			err: ErrInternal,
		},
		"fail auth": {
			id:  "fail_auth",
			err: auth.ErrAuth,
		},
	}
	//
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			err := svc.UpdateMetadata(ctx, "user0", c.id, c.md)
			if c.err == nil {
				assert.Nil(t, err)
			} else {
				assert.ErrorIs(t, err, c.err)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	//
	svc := NewService(newClientMock())
	//
	cases := map[string]struct {
		id     string
		err    error
		errMsg string
	}{
		"not found": {
			id:  "missing",
			err: ErrNotFound,
		},
		"ok": {
			id: "sub0",
		},
		"fail": {
			id:  "fail",
			err: ErrInternal,
		},
		"fail auth": {
			id:  "fail_auth",
			err: auth.ErrAuth,
		},
	}
	//
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			err := svc.Delete(ctx, "user0", c.id)
			if c.err == nil {
				assert.Nil(t, err)
			} else {
				assert.ErrorIs(t, err, c.err)
			}
		})
	}
}

func TestService_SearchOwn(t *testing.T) {
	//
	svc := NewService(newClientMock())
	//
	cases := map[string]struct {
		cursor string
		ids    []string
		err    error
	}{
		"ok": {
			ids: []string{
				"sub0",
				"sub1",
			},
		},
		"fail": {
			cursor: "fail",
			err:    ErrInternal,
		},
		"fail auth": {
			cursor: "fail_auth",
			err:    auth.ErrAuth,
		},
	}
	//
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			ids, err := svc.SearchOwn(ctx, "user0", 0, c.cursor)
			assert.Equal(t, ids, c.ids)
			assert.ErrorIs(t, err, c.err)
		})
	}
}

func conditionsDataEqual(a, b condition.Condition) (equal bool) {
	equal = a.IsNot() == b.IsNot()
	if equal {
		switch at := a.(type) {
		case condition.GroupCondition:
			switch bt := b.(type) {
			case condition.GroupCondition:
				equal = at.GetLogic() == bt.GetLogic()
				if equal {
					ag := at.GetGroup()
					bg := bt.GetGroup()
					equal = len(ag) == len(bg)
					if equal {
						for i, ac := range ag {
							equal = conditionsDataEqual(ac, bg[i])
							if !equal {
								break
							}
						}
					}
				}
			case condition.KiwiCondition:
				equal = false
			default:
				equal = false
			}
		case condition.KiwiCondition:
			switch bt := b.(type) {
			case condition.GroupCondition:
				equal = false
			case condition.KiwiCondition:
				equal = at.IsPartial() == bt.IsPartial() && at.GetKey() == bt.GetKey() && at.GetPattern() == bt.GetPattern()
			default:
				equal = false
			}
		}
	}
	return equal
}
