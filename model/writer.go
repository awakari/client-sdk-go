package model

import (
	"io"
)

type Writer[T any] interface {
	io.Closer
	WriteBatch(items []T) (ackCount uint32, err error)
}
