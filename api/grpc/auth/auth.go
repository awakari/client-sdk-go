package auth

import (
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
)

const keyUserId = "x-awakari-user-id"

var ErrAuth = errors.New("authentication failure")

func SetOutgoingAuthInfo(src context.Context, userId string) (dst context.Context) {
	dst = metadata.AppendToOutgoingContext(src, keyUserId, userId)
	return
}
