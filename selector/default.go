package selector

import (
	"errors"

	"github.com/ofavor/micro-lite/registry"
)

type defaultSelector struct {
	opts Options
	i    int
}

func newDefaultSelector(opts ...Option) Selector {
	options := defaultOptions()
	for _, o := range opts {
		o(&options)
	}
	return &defaultSelector{
		opts: options,
	}
}

func (s *defaultSelector) Select(services []*registry.Service) (*registry.Node, error) {
	if len(services) == 0 {
		return nil, errors.New("no valid service")
	}
	service := services[0]
	l := len(service.Nodes)
	if l == 0 {
		return nil, errors.New("no valid node")
	}
	ret := service.Nodes[s.i%l]
	s.i++
	return ret, nil
}
