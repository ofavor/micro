package micro

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ofavor/micro-lite/client"
	"github.com/ofavor/micro-lite/internal/log"
	"github.com/ofavor/micro-lite/server"
)

// Service interface
type Service interface {
	// Client get client instance
	Client() client.Client

	// Server get server instance
	Server() server.Server

	// Run the service
	Run() error
}

func newService(opts ...Option) Service {
	options := defaultOptions()
	for _, o := range opts {
		o(&options)
	}
	return &service{
		opts: options,
	}
}

type service struct {
	opts Options
}

func (s *service) Client() client.Client {
	return s.opts.Client
}

func (s *service) Server() server.Server {
	return s.opts.Server
}

func (s *service) Run() error {
	log.Info("Service is running ...")

	// start the server
	if err := s.opts.Server.Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	select {
	case <-ch:
	}

	if err := s.opts.Server.Stop(); err != nil {
		return err
	}

	log.Info("Service is terminated")
	return nil
}
