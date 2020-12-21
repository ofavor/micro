package client

import (
	"context"

	"github.com/google/uuid"
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

func (c *grpcClient) Init(o Option) {
	o(&c.opts)
}

func (c *grpcClient) Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error {
	log.Debug("gRPC client call:", req.Service(), req.Endpoint())
	services, err := c.opts.Registry.GetService(req.Service())
	if err != nil {
		return err
	}

	callOpts := c.opts.CallOpts
	for _, co := range opts {
		co(&callOpts)
	}
	log.Debug("Call options:", callOpts)
	var rErr error
	retries := callOpts.Retry

	for retries > 0 {
		node, err := c.opts.Selector.Select(services, callOpts.SelectOpts...)
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
		if err != nil { // request failed, retry
			retries--
			rErr = err
			continue
		}
		err = proto.Unmarshal(ret.Data, rsp.(proto.Message))
		if err != nil {
			return err
		}
		return nil
	}
	return rErr
}

type grpcRequest struct {
	id       string
	service  string
	endpoint string
	data     proto.Message
}

func newGRPCRequest(service string, endpoint string, req proto.Message) Request {
	return &grpcRequest{
		id:       uuid.New().String(),
		service:  service,
		endpoint: endpoint,
		data:     req,
	}
}

func (r *grpcRequest) ID() string {
	return r.id
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
