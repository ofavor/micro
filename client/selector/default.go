package selector

import (
	"errors"
	"math/rand"

	"github.com/ofavor/micro-lite/registry"
)

type defaultSelector struct {
	opts Options
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

func (s *defaultSelector) Select(services []*registry.Service, opts ...SelectOption) (*registry.Node, error) {
	sOpts := s.opts.SelectOpts
	for _, so := range opts {
		so(&sOpts)
	}

	for _, filter := range sOpts.Filters {
		services = filter(services)
	}
	if len(services) == 0 {
		return nil, errors.New("no valid service")
	}
	nodes := []*registry.Node{}
	for _, s := range services {
		nodes = append(nodes, s.Nodes...)
	}
	l := len(nodes)
	if l == 0 {
		return nil, errors.New("no valid node")
	}
	ret := nodes[rand.Int()%l]
	return ret, nil
}
