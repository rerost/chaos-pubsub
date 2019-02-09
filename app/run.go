package app

import (
	"github.com/rerost/chaos-pubsub/app/server"
	"github.com/rerost/chaos-pubsub/grpcserver"
	"github.com/srvc/fail"
)

// Run starts the grapiserver.
func Run() error {
	s := grpcserver.New(
		grpcserver.WithServers(
			server.NewPublisherServiceServer(nil),
			server.NewSubscriberServiceServer(nil),
		),
	)
	return fail.Wrap(s.Serve())
}
