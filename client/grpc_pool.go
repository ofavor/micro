package client

import (
	"errors"
	"sync"

	"github.com/ofavor/micro-lite/internal/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

var (
	maxPoolBufferSize = 20
)

type pool struct {
	sync.RWMutex
	conns map[string]*buffer //*grpc.ClientConn
}

type buffer struct {
	sync.RWMutex
	addr  string
	conns []*grpc.ClientConn
}

func newPool() *pool {
	return &pool{
		conns: make(map[string]*buffer), // make(map[string]*grpc.ClientConn),
	}
}

func newBuffer(addr string) *buffer {
	return &buffer{
		addr: addr,
	}
}

func (b *buffer) getConn() (*grpc.ClientConn, error) {
	b.Lock()
	var fc *grpc.ClientConn
	for i, c := range b.conns {
		switch c.GetState() {
		case connectivity.Connecting:
			continue
		case connectivity.Shutdown: // remove connection
			log.Debug("gRPC connection is shutdown, remove it")
			b.conns = append(b.conns[:i], b.conns[i+1:]...)
			continue
		case connectivity.TransientFailure: // remove connection
			log.Debug("gRPC connection transient failure, remove it")
			b.conns = append(b.conns[:i], b.conns[i+1:]...)
			continue
		case connectivity.Ready:
			fc = c
		case connectivity.Idle:
			fc = c
		}
		if fc != nil {
			b.Unlock()
			log.Debug("Got available gRPC connection from buffer")
			return fc, nil
		}
	}
	b.Unlock()
	if len(b.conns) >= maxPoolBufferSize {
		log.Warnf("Pool buffer size reach the max limitation:%d", maxPoolBufferSize)
		return nil, errors.New("pool buffer size reach the max limitation")
	}
	log.Debug("Create new gRPC connection")
	nc, err := grpc.Dial(
		b.addr,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	b.Lock()
	b.conns = append(b.conns, nc)
	b.Unlock()
	return nc, nil
}

func (p *pool) GetConn(addr string) (*grpc.ClientConn, error) {
	p.RLock()
	cs, ok := p.conns[addr]
	if ok {
		p.RUnlock()
		return cs.getConn()
	}
	p.RUnlock()
	cs = newBuffer(addr)
	p.Lock()
	p.conns[addr] = cs
	p.Unlock()
	return cs.getConn()
}
