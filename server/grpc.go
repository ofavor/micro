package server

import (
	"context"
	"net"

	"github.com/ofavor/micro-lite/internal/log"
	"github.com/ofavor/micro-lite/internal/transport"
	"google.golang.org/grpc"
)

type grpcServer struct {
	opts     Options
	srv      *grpc.Server
	handlers []Handler
}

func newGRPCServer(opts ...Option) Server {
	options := defaultOptions()
	for _, o := range opts {
		o(&options)
	}
	return &grpcServer{
		opts:     options,
		handlers: make([]Handler, 0),
	}
}

func (s *grpcServer) Init(opt Option) {
	opt(&s.opts)
}

func (s *grpcServer) Start() error {
	log.Infof("Trying to listen on TCP: %s", s.opts.Address)
	l, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		return err
	}

	s.srv = grpc.NewServer()
	transport.RegisterMicroServer(s.srv, s)

	go func() {
		if err := s.srv.Serve(l); err != nil {
			log.Error("gRPC server serve error: ", err)
		}
	}()
	return nil
}

func (s *grpcServer) Stop() error {
	return nil
}

func (s *grpcServer) AddHandler(h Handler) error {
	s.handlers = append(s.handlers, h)
	return nil
}

func (s *grpcServer) HandleRequest(ctx context.Context, req *transport.Request) (*transport.Response, error) {
	log.Debug("Got request:", req)
	h := s.handlers[0]
	ret, err := h.Invoke(ctx, req.Data)
	if err != nil {
		return nil, err
	}
	return &transport.Response{Data: ret}, nil
}
