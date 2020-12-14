package server

// Options for server
type Options struct {
	Name    string
	Version string
	Address string
}

// Option function to set server options
type Option func(opts *Options)

func defaultOptions() Options {
	return Options{
		Name:    "server",
		Version: "latest",
		Address: ":8888",
	}
}

// Name set name
func Name(name string) Option {
	return func(opts *Options) {
		opts.Name = name
	}
}

// Version set version
func Version(ver string) Option {
	return func(opts *Options) {
		opts.Version = ver
	}
}

// Address set address
func Address(addr string) Option {
	return func(opts *Options) {
		opts.Address = addr
	}
}
