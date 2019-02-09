package grpcserver

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/srvc/fail"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Engine interface {
	Serve() error
}

type Server interface {
	RegisterWithServer(grpcSvr *grpc.Server)
}

type engineImp struct {
	config Config
}

func New(opts ...Option) Engine {
	config := Config{
		Port: 5000,
	}
	for _, o := range opts {
		config = o(config)
	}

	return &engineImp{
		config: config,
	}
}

func (engine *engineImp) Serve() error {
	gserver := grpc.NewServer()

	for _, server := range engine.config.Servers {
		server.RegisterWithServer(gserver)
	}

	// Reflection.
	// TODO(@rerost) Move to option
	reflection.Register(gserver)

	// Start server
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Signal(syscall.SIGTERM), os.Signal(syscall.SIGINT))

	network, err := net.Listen("tcp", fmt.Sprintf(":%d", engine.config.Port))

	start := func() <-chan error {
		errCh := make(chan error)

		go func(errCh chan<- error) {
			fmt.Printf("Server is statig at %d\n", engine.config.Port)
			if err = gserver.Serve(network); err != nil {
				errCh <- fail.Wrap(err)
			}

			close(errCh)
		}(errCh)
		return errCh
	}

	select {
	case <-quit:
		gserver.GracefulStop()
	case err := <-start():
		if err != nil {
			return fail.Wrap(err)
		}
	}

	return nil
}
