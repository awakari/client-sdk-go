package model

import (
	"io"
)

type WriteStream[T any] interface {
	io.Closer
	WriteBatch(items []T) (ackCount uint32, err error)
}
