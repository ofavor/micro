package registry

// Registry interface
type Registry interface {
	// Initialize registry with option
	Init(Option)

	// Register service to registry
	Register(*Service, ...Option) error

	// Deregister service from registry
	Deregister(*Service) error

	// GetService by specifying service name
	GetService(string) ([]*Service, error)
}

// Service struct
type Service struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints []*Endpoint       `json:"endpoints"`
	Nodes     []*Node           `json:"nodes"`
}

// Node struct
type Node struct {
	ID       string            `json:"id"`
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
}

// Endpoint struct
type Endpoint struct {
	Name     string            `json:"name"`
	Request  *Value            `json:"request"`
	Response *Value            `json:"response"`
	Metadata map[string]string `json:"metadata"`
}

// Value struct
type Value struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Values []*Value `json:"values"`
}

// NewRegistry create new registry
func NewRegistry(opts ...Option) Registry {
	return newETCDRegistry(opts...)
}
