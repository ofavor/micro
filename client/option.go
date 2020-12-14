package client

// Options for client
type Options struct {
}

// Option function to set client options
type Option func(opts *Options)

func defaultOptions() Options {
	return Options{}
}
