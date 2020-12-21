package client

import (
	"context"

	"google.golang.org/protobuf/proto"
)

// Client interface
type Client interface {
	// Init client with option function
	Init(Option)

	// Call remote function
	Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error
}

// Request interface
type Request interface {
	ID() string
	Service() string
	Endpoint() string
	Data() proto.Message
}

// Response interface
type Response interface {
	ID() string
	Data() proto.Message
}

// NewClient create new client
func NewClient(opts ...Option) Client {
	return newGRPCClient(opts...)
}

// NewRequest create new request
func NewRequest(service string, method string, req proto.Message) Request {
	return newGRPCRequest(service, method, req)
}
