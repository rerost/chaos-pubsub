package app

import (
	"github.com/rerost/chaos-pubsub/app/server"
	"github.com/rerost/chaos-pubsub/grpcserver"
)

// Run starts the grapiserver.
func Run() error {
	s := grpcserver.New(
		grpcserver.WithServers(
			server.NewPublisherServiceServer(),
			server.NewSubscriberServiceServer(),
		),
	)
	return s.Serve()
}
