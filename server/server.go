package server

import (
	"context"

	"github.com/ofavor/micro-lite/internal/transport"
	"github.com/ofavor/micro-lite/registry"
)

// Server interface
type Server interface {
	// ID get server id
	ID() string

	// Init server with option
	Init(Option)

	// Start the server
	Start() error

	// Stop the server
	Stop() error

	// Handle register handler to handle rpc request
	Handle(Handler) error
}

// Handler interface
type Handler interface {
	// Name get handler name
	Name() string

	// Target get handler target object
	Target() interface{}

	// Endpoints get handler endpoints
	Endpoints() []*registry.Endpoint
}

// HandlerFunc handler function
type HandlerFunc func(ctx context.Context, req *transport.Request, rsp interface{}) error

// HandlerWrapper wrapper for handler function
type HandlerWrapper func(HandlerFunc) HandlerFunc

// NewServer create new server
func NewServer(opts ...Option) Server {
	return newGRPCServer(opts...)
}

// NewHandler create new handler
func NewHandler(name string, target interface{}) Handler {
	return newHandler(name, target)
}
