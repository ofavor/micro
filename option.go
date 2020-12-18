package micro

import (
	"time"

	"github.com/ofavor/micro-lite/client"
	"github.com/ofavor/micro-lite/internal/log"
	"github.com/ofavor/micro-lite/registry"
	"github.com/ofavor/micro-lite/selector"
	"github.com/ofavor/micro-lite/server"
)

// Options options for micro service
type Options struct {
	Server server.Server
	Client client.Client
}

// Option function to set micro service options
type Option func(opts *Options)

func defaultOptions() Options {
	return Options{
		Server: server.NewServer(),
		Client: client.NewClient(),
	}
}

// LogLevel set log level, supported levels: "debug", "info", "warn", "error", "fatal"
func LogLevel(lv string) Option {
	return func(opts *Options) {
		log.SetLevel(lv)
	}
}

// Name set name
func Name(name string) Option {
	return func(opts *Options) {
		opts.Server.Init(server.Name(name))
	}
}

// Version set version
func Version(ver string) Option {
	return func(opts *Options) {
		opts.Server.Init(server.Version(ver))
	}
}

// Address set address
func Address(addr string) Option {
	return func(opts *Options) {
		opts.Server.Init(server.Address(addr))
	}
}

// Registry set registry
func Registry(reg registry.Registry) Option {
	return func(opts *Options) {
		opts.Server.Init(server.Registry(reg))
		opts.Client.Init(client.Registry(reg))
	}
}

// Selector set selector
func Selector(sel selector.Selector) Option {
	return func(opts *Options) {
		opts.Client.Init(client.Selector(sel))
	}
}

// RegisterInterval set server register interval
func RegisterInterval(d time.Duration) Option {
	return func(opts *Options) {
		opts.Server.Init(server.RegisterInterval(d))
	}
}

// RegisterTTL set server register TTL
func RegisterTTL(d time.Duration) Option {
	return func(opts *Options) {
		opts.Server.Init(server.RegisterTTL(d))
	}
}
