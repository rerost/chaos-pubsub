package server

import (
	"context"
	"math/rand"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rerost/chaos-pubsub/lib/grpcserver"
	api_pb "google.golang.org/genproto/googleapis/pubsub/v1"
)

// PublisherServiceServer is a composite interface of api_pb.PublisherServiceServer and grapiserver.Server.
type PublisherServiceServer interface {
	api_pb.PublisherServer
	grpcserver.Server
}

// NewPublisherServiceServer creates a new PublisherServiceServer instance.
func NewPublisherServiceServer(rawClient api_pb.PublisherClient) PublisherServiceServer {
	return &publisherServiceServerImpl{
		rawClient: rawClient,
	}
}

type publisherServiceServerImpl struct {
	rawClient api_pb.PublisherClient
}

func (server *publisherServiceServerImpl) CreateTopic(ctx context.Context, topic *api_pb.Topic) (*api_pb.Topic, error) {
	return server.rawClient.CreateTopic(ctx, topic)
}
func (server *publisherServiceServerImpl) UpdateTopic(ctx context.Context, updateTopicRequest *api_pb.UpdateTopicRequest) (*api_pb.Topic, error) {
	return server.rawClient.UpdateTopic(ctx, updateTopicRequest)
}
func (server *publisherServiceServerImpl) Publish(ctx context.Context, publishRequest *api_pb.PublishRequest) (*api_pb.PublishResponse, error) {
	res, err := server.rawClient.Publish(ctx, publishRequest)
	var wg sync.WaitGroup
	n := rand.Intn(10) + 1
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(server *publisherServiceServerImpl) {
			server.rawClient.Publish(ctx, publishRequest)
			wg.Done()
		}(server)
	}

	wg.Wait()
	return res, err
}
func (server *publisherServiceServerImpl) GetTopic(ctx context.Context, getTopicRequest *api_pb.GetTopicRequest) (*api_pb.Topic, error) {
	return server.rawClient.GetTopic(ctx, getTopicRequest)
}
func (server *publisherServiceServerImpl) ListTopics(ctx context.Context, listTopicRequest *api_pb.ListTopicsRequest) (*api_pb.ListTopicsResponse, error) {
	return server.rawClient.ListTopics(ctx, listTopicRequest)
}
func (server *publisherServiceServerImpl) ListTopicSubscriptions(ctx context.Context, listTopicSubscritpionsRequest *api_pb.ListTopicSubscriptionsRequest) (*api_pb.ListTopicSubscriptionsResponse, error) {
	return server.rawClient.ListTopicSubscriptions(ctx, listTopicSubscritpionsRequest)
}
func (server *publisherServiceServerImpl) ListTopicSnapshots(ctx context.Context, listTopicSnapshotsRequest *api_pb.ListTopicSnapshotsRequest) (*api_pb.ListTopicSnapshotsResponse, error) {
	return server.rawClient.ListTopicSnapshots(ctx, listTopicSnapshotsRequest)
}
func (server *publisherServiceServerImpl) DeleteTopic(ctx context.Context, deleteTopicRequest *api_pb.DeleteTopicRequest) (*empty.Empty, error) {
	return server.rawClient.DeleteTopic(ctx, deleteTopicRequest)
}
