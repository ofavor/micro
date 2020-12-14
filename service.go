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
	Client() client.Client
	Server() server.Server
	Run() error
}

func newService(opts ...Option) Service {
	options := defaultOptions()
	for _, o := range opts {
		o(&options)
	}
	return &service{
		Options: options,
	}
}

type service struct {
	Options
}

func (s *service) Client() client.Client {
	return s.Options.Client
}

func (s *service) Server() server.Server {
	return s.Options.Server
}

func (s *service) Run() error {
	log.Info("Service is running ...")

	if err := s.Options.Server.Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	select {
	case <-ch:
	}

	if err := s.Options.Server.Stop(); err != nil {
		return err
	}

	log.Info("Service is terminated")
	return nil
}
