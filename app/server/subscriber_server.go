package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/izumin5210/grapi/pkg/grapiserver"
	api_pb "github.com/rerost/chaos-pubsub/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SubscriberServiceServer is a composite interface of api_pb.SubscriberServiceServer and grapiserver.Server.
type SubscriberServiceServer interface {
	api_pb.SubscriberServer
	grapiserver.Server
}

// NewSubscriberServiceServer creates a new SubscriberServiceServer instance.
func NewSubscriberServiceServer() SubscriberServiceServer {
	return &subscriberServiceServerImpl{}
}

type subscriberServiceServerImpl struct {
}

func (*subscriberServiceServerImpl) CreateSubscription(context.Context, *api_pb.Subscription) (*api_pb.Subscription, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*subscriberServiceServerImpl) GetSubscription(context.Context, *api_pb.GetSubscriptionRequest) (*api_pb.Subscription, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*subscriberServiceServerImpl) UpdateSubscription(context.Context, *api_pb.UpdateSubscriptionRequest) (*api_pb.Subscription, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*subscriberServiceServerImpl) ListSubscriptions(context.Context, *api_pb.ListSubscriptionsRequest) (*api_pb.ListSubscriptionsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*subscriberServiceServerImpl) DeleteSubscription(context.Context, *api_pb.DeleteSubscriptionRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*subscriberServiceServerImpl) ModifyAckDeadline(context.Context, *api_pb.ModifyAckDeadlineRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*subscriberServiceServerImpl) Acknowledge(context.Context, *api_pb.AcknowledgeRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*subscriberServiceServerImpl) Pull(context.Context, *api_pb.PullRequest) (*api_pb.PullResponse, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*subscriberServiceServerImpl) StreamingPull(api_pb.Subscriber_StreamingPullServer) error {
	return status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*subscriberServiceServerImpl) ModifyPushConfig(context.Context, *api_pb.ModifyPushConfigRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*subscriberServiceServerImpl) GetSnapshot(context.Context, *api_pb.GetSnapshotRequest) (*api_pb.Snapshot, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*subscriberServiceServerImpl) ListSnapshots(context.Context, *api_pb.ListSnapshotsRequest) (*api_pb.ListSnapshotsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*subscriberServiceServerImpl) CreateSnapshot(context.Context, *api_pb.CreateSnapshotRequest) (*api_pb.Snapshot, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*subscriberServiceServerImpl) UpdateSnapshot(context.Context, *api_pb.UpdateSnapshotRequest) (*api_pb.Snapshot, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*subscriberServiceServerImpl) DeleteSnapshot(context.Context, *api_pb.DeleteSnapshotRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
func (*subscriberServiceServerImpl) Seek(context.Context, *api_pb.SeekRequest) (*api_pb.SeekResponse, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
