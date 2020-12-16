package client

import "google.golang.org/grpc"

type pool struct {
	conns map[string]*grpc.ClientConn
}

func newPool() *pool {
	return &pool{
		conns: make(map[string]*grpc.ClientConn),
	}
}

func (p *pool) GetConn(addr string) (*grpc.ClientConn, error) {
	c, ok := p.conns[addr]
	if ok {
		return c, nil
	}
	c, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	p.conns[addr] = c
	return c, nil
}
