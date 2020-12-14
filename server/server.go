package server

import (
	"context"
)

// Server interface
type Server interface {
	Init(Option)
	Start() error
	Stop() error
	AddHandler(Handler) error
}

// Handler interface
type Handler interface {
	Invoke(context.Context, []byte) ([]byte, error)
}

// NewServer create new server
func NewServer(opts ...Option) Server {
	return newGRPCServer(opts...)
}
