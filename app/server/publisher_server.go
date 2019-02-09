package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rerost/chaos-pubsub/grpcserver"
	api_pb "google.golang.org/genproto/googleapis/pubsub/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PublisherServiceServer is a composite interface of api_pb.PublisherServiceServer and grapiserver.Server.
type PublisherServiceServer interface {
	api_pb.PublisherServer
	grpcserver.Server
}

// NewPublisherServiceServer creates a new PublisherServiceServer instance.
func NewPublisherServiceServer() PublisherServiceServer {
	return &publisherServiceServerImpl{}
}

type publisherServiceServerImpl struct {
}

func (*publisherServiceServerImpl) CreateTopic(context.Context, *api_pb.Topic) (*api_pb.Topic, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*publisherServiceServerImpl) UpdateTopic(context.Context, *api_pb.UpdateTopicRequest) (*api_pb.Topic, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*publisherServiceServerImpl) Publish(context.Context, *api_pb.PublishRequest) (*api_pb.PublishResponse, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*publisherServiceServerImpl) GetTopic(context.Context, *api_pb.GetTopicRequest) (*api_pb.Topic, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*publisherServiceServerImpl) ListTopics(context.Context, *api_pb.ListTopicsRequest) (*api_pb.ListTopicsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*publisherServiceServerImpl) ListTopicSubscriptions(context.Context, *api_pb.ListTopicSubscriptionsRequest) (*api_pb.ListTopicSubscriptionsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*publisherServiceServerImpl) ListTopicSnapshots(context.Context, *api_pb.ListTopicSnapshotsRequest) (*api_pb.ListTopicSnapshotsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*publisherServiceServerImpl) DeleteTopic(context.Context, *api_pb.DeleteTopicRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
