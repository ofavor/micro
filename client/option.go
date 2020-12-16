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
