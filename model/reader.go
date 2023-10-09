package model

import (
	"io"
)

type Reader[T any] interface {
	io.Closer
	Read() (item T, err error)
}

type AckReader[T any] interface {
	Reader[T]
	Ack(count uint32) (err error)
}
