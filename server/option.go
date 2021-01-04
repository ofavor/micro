package server

import (
	"time"

	"github.com/google/uuid"
	"github.com/ofavor/micro-lite/registry"
)

// Options for server
type Options struct {
	ID               string
	Name             string
	Version          string
	Address          string
	Registry         registry.Registry
	RegisterInterval time.Duration
	RegisterTTL      time.Duration
	HdlrWrappers     []HandlerWrapper
}

// Option function to set server options
type Option func(opts *Options)

func defaultOptions() Options {
	return Options{
		ID:               uuid.New().String(),
		Name:             "server",
		Version:          "1.0.0",
		Address:          ":8888",
		Registry:         registry.NewRegistry(),
		RegisterInterval: 30 * time.Second,
		RegisterTTL:      60 * time.Second,
		HdlrWrappers:     []HandlerWrapper{},
	}
}

// ID set id
func ID(id string) Option {
	return func(opts *Options) {
		opts.ID = id
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

// Registry set registry
func Registry(reg registry.Registry) Option {
	return func(opts *Options) {
		opts.Registry = reg
	}
}

// RegistryAddrs set registry addresses
func RegistryAddrs(addrs []string) Option {
	return func(opts *Options) {
		opts.Registry.Init(registry.Addrs(addrs))
	}
}

// RegisterInterval set register interval
func RegisterInterval(d time.Duration) Option {
	return func(opts *Options) {
		opts.RegisterInterval = d
	}
}

// RegisterTTL set register TTL
func RegisterTTL(d time.Duration) Option {
	return func(opts *Options) {
		opts.RegisterTTL = d
	}
}

// WrapHandler add handler wrapper
func WrapHandler(w HandlerWrapper) Option {
	return func(opts *Options) {
		opts.HdlrWrappers = append(opts.HdlrWrappers, w)
	}
}
