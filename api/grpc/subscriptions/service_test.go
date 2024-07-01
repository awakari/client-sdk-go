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
				Description: "invalid",
				Condition: condition.NewTextCondition(
					condition.NewKeyCondition(condition.NewCondition(false), "key0"),
					"ok", false,
				),
			},
			err: ErrInvalid,
		},
		"locked": {
			req: subscription.Data{
				Description: "busy",
				Condition: condition.NewTextCondition(
					condition.NewKeyCondition(condition.NewCondition(false), ""),
					"locked", false,
				),
			},
			err: ErrUnavailable,
		},
		"fail": {
			req: subscription.Data{
				Description: "fail",
				Condition: condition.NewTextCondition(
					condition.NewKeyCondition(
						condition.NewCondition(false),
						"fail",
					),
					"fail", false,
				),
			},
			err: ErrInternal,
		},
		"ok": {
			req: subscription.Data{
				Description: "my subscription",
				Condition: condition.NewTextCondition(
					condition.NewKeyCondition(condition.NewCondition(false), "key0"),
					"ok", false,
				),
			},
		},
		"ok group": {
			req: subscription.Data{
				Description: "my subscription",
				Condition: condition.
					NewBuilder().
					All(
						[]condition.Condition{
							condition.
								NewBuilder().
								BuildTextCondition(),
						},
					).
					BuildGroupCondition(),
			},
		},
		"auth fail": {
			req: subscription.Data{
				Description: "fail_auth",
				Condition: condition.NewTextCondition(
					condition.NewKeyCondition(condition.NewCondition(false), "key0"),
					"ok", false,
				),
			},
			err: auth.ErrAuth,
		},
		"limit reached": {
			req: subscription.Data{
				Description: "limit_reached",
				Condition: condition.NewTextCondition(
					condition.NewKeyCondition(condition.NewCondition(false), "key0"),
					"ok", false,
				),
			},
			err: limits.ErrReached,
		},
		"ok with expiration": {
			req: subscription.Data{
				Description: "my subscription",
				Expires:     time.Now(),
				Condition: condition.NewTextCondition(
					condition.NewKeyCondition(condition.NewCondition(false), "key0"),
					"ok", false,
				),
			},
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
				Description: "subscription",
				Enabled:     true,
				Expires:     time.Date(2023, 10, 4, 11, 44, 55, 0, time.UTC),
				Condition: condition.
					NewBuilder().
					Any(
						[]condition.Condition{
							condition.
								NewBuilder().
								Not().
								AttributeKey("k0").
								AnyOfWords("p0").
								BuildTextCondition(),
							condition.
								NewBuilder().
								AttributeKey("k1").
								GreaterThan(-42.1).
								BuildNumberCondition(),
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
				assert.Equal(t, c.sd.Description, sd.Description)
				assert.Equal(t, c.sd.Enabled, sd.Enabled)
				assert.Equal(t, c.sd.Expires, sd.Expires)
				assert.True(t, conditionsDataEqual(c.sd.Condition, sd.Condition))
			} else {
				assert.ErrorIs(t, err, c.err)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	//
	svc := NewService(newClientMock())
	//
	cases := map[string]struct {
		id  string
		err error
		sd  subscription.Data
	}{
		"not found": {
			id: "missing",
			sd: subscription.Data{
				Condition: condition.NewTextCondition(
					condition.NewKeyCondition(
						condition.NewCondition(false),
						"key1",
					),
					"term1",
					true,
				),
			},
			err: ErrNotFound,
		},
		"ok": {
			id: "sub0",
			sd: subscription.Data{
				Description: "my subscription",
				Enabled:     false,
				Expires:     time.Now(),
				Condition: condition.NewGroupCondition(
					condition.NewCondition(false),
					condition.GroupLogicOr,
					[]condition.Condition{
						condition.NewNumberCondition(
							condition.NewKeyCondition(
								condition.NewCondition(true),
								"key0",
							),
							condition.NumOpGt,
							42,
						),
						condition.NewTextCondition(
							condition.NewKeyCondition(
								condition.NewCondition(false),
								"key1",
							),
							"term1",
							true,
						),
					},
				),
			},
		},
		"fail": {
			id: "fail",
			sd: subscription.Data{
				Condition: condition.NewTextCondition(
					condition.NewKeyCondition(
						condition.NewCondition(false),
						"key1",
					),
					"term1",
					true,
				),
			},
			err: ErrInternal,
		},
		"missing condition": {
			id:  "sub0",
			err: ErrInvalid,
		},
		"fail auth": {
			id: "fail_auth",
			sd: subscription.Data{
				Condition: condition.NewTextCondition(
					condition.NewKeyCondition(
						condition.NewCondition(false),
						"key1",
					),
					"term1",
					true,
				),
			},
			err: auth.ErrAuth,
		},
	}
	//
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			err := svc.Update(ctx, "user0", c.id, c.sd)
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

func TestService_Search(t *testing.T) {
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
			ids, err := svc.Search(ctx, "user0", subscription.Query{}, subscription.Cursor{
				Id: c.cursor,
			})
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
			default:
				equal = false
			}
		case condition.TextCondition:
			switch bt := b.(type) {
			case condition.TextCondition:
				equal = at.GetKey() == bt.GetKey() && at.GetTerm() == bt.GetTerm()
			default:
				equal = false
			}
		case condition.NumberCondition:
			switch bt := b.(type) {
			case condition.NumberCondition:
				equal = at.GetKey() == bt.GetKey() && at.GetOperation() == bt.GetOperation() && at.GetValue() == bt.GetValue()
			default:
				equal = false
			}
		}
	}
	return equal
}
