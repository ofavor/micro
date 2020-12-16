package registry

import "time"

// Options of registry
type Options struct {
	Addrs []string
	TTL   time.Duration
}

// Option function to set registry options
type Option func(opts *Options)

func defaultOptions() Options {
	return Options{
		Addrs: []string{"127.0.0.1:2379"},
	}
}

// Addrs registry addresses
func Addrs(addrs []string) Option {
	return func(opts *Options) {
		opts.Addrs = addrs
	}
}

// TTL registry ttl
func TTL(ttl time.Duration) Option {
	return func(opts *Options) {
		opts.TTL = ttl
	}
}
