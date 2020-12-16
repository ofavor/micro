package server

import "github.com/ofavor/micro-lite/registry"

// Server interface
type Server interface {
	Init(Option)
	Start() error
	Stop() error
	Handle(Handler) error
}

// Handler interface
type Handler interface {
	Name() string
	Target() interface{}
	Endpoints() []*registry.Endpoint
}

// NewServer create new server
func NewServer(opts ...Option) Server {
	return newGRPCServer(opts...)
}

// NewHandler create new handler
func NewHandler(name string, target interface{}) Handler {
	return newHandler(name, target)
}
