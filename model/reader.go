package model

import (
	"io"
)

type Reader[T any] interface {
	io.Closer
	Read() (item T, err error)
}
