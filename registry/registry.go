package registry

// Registry interface
type Registry interface {
	Register(*Service) error
	Deregister(*Service) error
}

// Service struct
type Service struct {
	Name    string
	Version string
}

// Node struct
type Node struct {
	ID string
}

// NewRegistry create new registry
func NewRegistry(opts ...Option) Registry {
	return newETCDRegistry(opts...)
}
