package foo

import (
	"context"

	"github.com/ofavor/micro-lite/server"

	"github.com/ofavor/micro-lite/client"
)

type FooService interface {
	Bar(ctx context.Context, req *Request, opts ...client.CallOption) (*Response, error)
}

func NewFooService(c client.Client) FooService {
	return &fooService{
		c: c,
	}
}

type fooService struct {
	c client.Client
}

func (s *fooService) Bar(ctx context.Context, req *Request, opts ...client.CallOption) (*Response, error) {
	in := client.NewRequest("", "Foo.Bar", req)
	out := new(Response)
	if err := s.c.Call(ctx, in, out); err != nil {
		return nil, err
	}
	return out, nil
}

type FooHandler interface {
	Bar(ctx context.Context, req *Request, rsp *Response) error
}

func RegisterFooHandler(s server.Server, h FooHandler) {
	hdr := server.NewHandler("Foo", h)
	s.Handle(hdr)
}
