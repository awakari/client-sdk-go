package model

import (
	"io"
)

type ReadStream[T any] interface {
	io.Closer
	Read() (item T, err error)
}
