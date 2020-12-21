package client

import (
	"time"

	"github.com/ofavor/micro-lite/client/selector"
	"github.com/ofavor/micro-lite/registry"
)

// Options for client
type Options struct {
	Registry registry.Registry
	Selector selector.Selector

	CallOpts CallOptions
}

// Option function to set client options
type Option func(opts *Options)

// CallOptions rpc call options
type CallOptions struct {
	SelectOpts     []selector.SelectOption
	RequestTimeout time.Duration
	Retry          int
}

// CallOption function to set rpc call options
type CallOption func(opts *CallOptions)

func defaultOptions() Options {
	return Options{
		Registry: registry.NewRegistry(),
		Selector: selector.NewSelector(),
		CallOpts: CallOptions{
			Retry: 1,
		},
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

// WithSelectOption set select option for call
func WithSelectOption(so selector.SelectOption) CallOption {
	return func(cos *CallOptions) {
		cos.SelectOpts = append(cos.SelectOpts, so)
	}
}

// WithRequestTimeout set request timeout option for call
func WithRequestTimeout(t time.Duration) CallOption {
	return func(cos *CallOptions) {
		cos.RequestTimeout = t
	}
}

// WithRetry set request failure retry count option for call
func WithRetry(r int) CallOption {
	return func(cos *CallOptions) {
		cos.Retry = r
	}
}
