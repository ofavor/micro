package selector

import "github.com/ofavor/micro-lite/registry"

// Selector interface
type Selector interface {
	Select([]*registry.Service, ...SelectOption) (*registry.Node, error)
}

// NewSelector create new selector
func NewSelector(opts ...Option) Selector {
	return newDefaultSelector(opts...)
}
