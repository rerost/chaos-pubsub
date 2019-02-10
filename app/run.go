package app

import (
	"context"

	"github.com/rerost/chaos-pubsub/app/server"
	"github.com/rerost/chaos-pubsub/infra/pubsub"
	"github.com/rerost/chaos-pubsub/lib/grpcserver"
	"github.com/rerost/chaos-pubsub/lib/interceptor/logger"
	"github.com/srvc/fail"
	"google.golang.org/grpc"
)

// Run starts the grapiserver.
func Run() error {
	conn, err := grpc.Dial("localhost:8085", grpc.WithInsecure())
	if err != nil {
		return fail.Wrap(err)
	}
	defer conn.Close()

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "test", []byte{})
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
