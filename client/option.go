package client

import (
	"github.com/ofavor/micro-lite/registry"
	"github.com/ofavor/micro-lite/selector"
)

// Options for client
type Options struct {
	Registry registry.Registry
	Selector selector.Selector
}

// Option function to set client options
type Option func(opts *Options)

func defaultOptions() Options {
	return Options{
		Registry: registry.NewRegistry(),
		Selector: selector.NewSelector(),
	}
}

// Registry set registry
func Registry(reg registry.Registry) Option {
	return func(opts *Options) {
		opts.Registry = reg
	}
}

// Selector set selector
func Selector(sel selector.Selector) Option {
	return func(opts *Options) {
		opts.Selector = sel
	}
}
