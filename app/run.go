package app

import (
	"github.com/izumin5210/grapi/pkg/grapiserver"
	"github.com/rerost/chaos-pubsub/app/server"
)

// Run starts the grapiserver.
func Run() error {
	s := grapiserver.New(
		grapiserver.WithDefaultLogger(),
		grapiserver.WithServers(
			server.NewPublisherServiceServer(),
			server.NewSubscriberServiceServer(),
		),
	)
	return s.Serve()
}
