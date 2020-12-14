package micro

// NewService create new service with options
func NewService(opts ...Option) Service {
	return newService(opts...)
}
