package server

import (
	"github.com/google/uuid"
	"github.com/ofavor/micro-lite/registry"
)

// Options for server
type Options struct {
	ID       string
	Name     string
	Version  string
	Address  string
	Registry registry.Registry
}

// Option function to set server options
type Option func(opts *Options)

func defaultOptions() Options {
	return Options{
		ID:       uuid.New().String(),
		Name:     "server",
		Version:  "latest",
		Address:  ":8888",
		Registry: registry.NewRegistry(),
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
