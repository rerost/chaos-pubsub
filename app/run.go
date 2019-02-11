package app

import (
	"context"
	"os"

	"github.com/rerost/chaos-pubsub/app/server"
	"github.com/rerost/chaos-pubsub/infra/pubsub"
	"github.com/rerost/chaos-pubsub/lib/grpcserver"
	"github.com/rerost/chaos-pubsub/lib/interceptor/logger"
	"github.com/srvc/fail"
)

// Run starts the grapiserver.
func Run() error {
	ctx := context.Background()
	projectName := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectName == "" {
		fail.New("Please set env GOOGLE_CLOUD_PROJECT")
	}
	client, err := pubsub.NewClient(ctx, projectName)
	if err != nil {
		return fail.Wrap(err)
	}

	s := grpcserver.New(
		grpcserver.WithServers(
			server.NewPublisherServiceServer(client.PublisherClient),
			server.NewSubscriberServiceServer(client.SubscriberClient),
		),
		grpcserver.WithGrpcServerUnaryInterceptors(
			logger.UnaryServerInterceptor(),
		),
		grpcserver.WithGrpcServerStreamInterceptors(
			logger.StreamServerInterceptor(),
		),
	)
	return fail.Wrap(s.Serve())
}
