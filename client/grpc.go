package client

import (
	"context"

	"github.com/ofavor/micro-lite/internal/log"
	"github.com/ofavor/micro-lite/internal/transport"
	"google.golang.org/protobuf/proto"
)

type grpcClient struct {
	opts Options

	connPool *pool
}

func newGRPCClient(opts ...Option) Client {
	options := defaultOptions()
	for _, o := range opts {
		o(&options)
	}
	return &grpcClient{
		opts:     options,
		connPool: newPool(),
	}
}

func (c *grpcClient) Call(ctx context.Context, req Request, rsp proto.Message, opts ...CallOption) error {
	log.Debug("Client call:", req.Endpoint())
	services, err := c.opts.Registry.GetService(req.Service())
	if err != nil {
		return err
	}
	node, err := c.opts.Selector.Select(services)
	if err != nil {
		return err
	}
	conn, err := c.connPool.GetConn(node.Address)
	if err != nil {
		return err
	}
	gc := transport.NewMicroClient(conn)
	data, err := proto.Marshal(req.Data())
	if err != nil {
		return err
	}
	in := &transport.Request{
		Service:  req.Service(),
		Endpoint: req.Endpoint(),
		Data:     data,
	}

	ret, err := gc.HandleRequest(ctx, in)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(ret.Data, rsp)
	if err != nil {
		return err
	}
	return nil
}

type grpcRequest struct {
	service  string
	endpoint string
	data     proto.Message
}

func newGRPCRequest(service string, endpoint string, req proto.Message) Request {
	return &grpcRequest{
		service:  service,
		endpoint: endpoint,
		data:     req,
	}
}

func (r *grpcRequest) Service() string {
	return r.service
}

func (r *grpcRequest) Endpoint() string {
	return r.endpoint
}

func (r *grpcRequest) Data() proto.Message {
	return r.data
}
