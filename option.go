package micro

import (
	"github.com/ofavor/micro-lite/client"
	"github.com/ofavor/micro-lite/internal/log"
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
