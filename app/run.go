package app

import (
	"github.com/rerost/chaos-pubsub/app/server"
	"github.com/rerost/chaos-pubsub/grpcserver"
	"github.com/srvc/fail"
	api_pb "google.golang.org/genproto/googleapis/pubsub/v1"
	"google.golang.org/grpc"
)

// Run starts the grapiserver.
func Run() error {
	conn, err := grpc.Dial("localhost:8085", grpc.WithInsecure())
	if err != nil {
		return fail.Wrap(err)
	}
	defer conn.Close()

	publisherClient := api_pb.NewPublisherClient(conn)
	subscriberClient := api_pb.NewSubscriberClient(conn)

	s := grpcserver.New(
		grpcserver.WithServers(
			server.NewPublisherServiceServer(publisherClient),
			server.NewSubscriberServiceServer(subscriberClient),
		),
	)
	return fail.Wrap(s.Serve())
}
