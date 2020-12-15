package server

// Server interface
type Server interface {
	Init(Option)
	Start() error
	Stop() error
	Handle(Handler) error
}

// Handler interface
type Handler interface {
	Name() string
	Target() interface{}
	Endpoints() []interface{}
}

// NewServer create new server
func NewServer(opts ...Option) Server {
	return newGRPCServer(opts...)
}

// NewHandler create new handler
func NewHandler(name string, target interface{}) Handler {
	return newHandler(name, target)
}
