package app

import (
	"context"
	"os"
	"strconv"

	"github.com/rerost/chaos-pubsub/app/server"
	"github.com/rerost/chaos-pubsub/infra/pubsub"
	"github.com/rerost/chaos-pubsub/lib/grpcserver"
	"github.com/rerost/chaos-pubsub/lib/interceptor/fault"
	"github.com/rerost/chaos-pubsub/lib/interceptor/logger"
	"github.com/srvc/fail"
)

// Run starts the grapiserver.
func Run() error {
	ctx := context.Background()
	port := int64(5000)
	if _port := os.Getenv("APP_PORT"); _port != "" {
		p, err := strconv.ParseInt(_port, 10, 64)
		if err != nil {
			return err
		}
		port = p
	}

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
			fault.UnaryServerInterceptor(),
		),
		grpcserver.WithGrpcServerStreamInterceptors(
			logger.StreamServerInterceptor(),
			fault.StreamServerInterceptor(),
		),
		grpcserver.WithPort(port),
	)
	return fail.Wrap(s.Serve())
}
