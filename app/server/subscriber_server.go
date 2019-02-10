package server

import (
	"context"
	"io"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rerost/chaos-pubsub/lib/grpcserver"
	api_pb "google.golang.org/genproto/googleapis/pubsub/v1"
)

// SubscriberServiceServer is a composite interface of api_pb.SubscriberServiceServer and grapiserver.Server.
type SubscriberServiceServer interface {
	api_pb.SubscriberServer
	grpcserver.Server
}

// NewSubscriberServiceServer creates a new SubscriberServiceServer instance.
func NewSubscriberServiceServer(rawClient api_pb.SubscriberClient) SubscriberServiceServer {
	return &subscriberServiceServerImpl{
		rawClient: rawClient,
	}
}

type subscriberServiceServerImpl struct {
	rawClient api_pb.SubscriberClient
}

func (server *subscriberServiceServerImpl) CreateSubscription(ctx context.Context, subscription *api_pb.Subscription) (*api_pb.Subscription, error) {
	return server.rawClient.CreateSubscription(ctx, subscription)
}
func (server *subscriberServiceServerImpl) GetSubscription(ctx context.Context, getSubscriptionRequest *api_pb.GetSubscriptionRequest) (*api_pb.Subscription, error) {
	return server.rawClient.GetSubscription(ctx, getSubscriptionRequest)
}
func (server *subscriberServiceServerImpl) UpdateSubscription(ctx context.Context, updateSubscriptionRequest *api_pb.UpdateSubscriptionRequest) (*api_pb.Subscription, error) {
	return server.rawClient.UpdateSubscription(ctx, updateSubscriptionRequest)
}
func (server *subscriberServiceServerImpl) ListSubscriptions(ctx context.Context, listSubscriptionsRequest *api_pb.ListSubscriptionsRequest) (*api_pb.ListSubscriptionsResponse, error) {
	return server.rawClient.ListSubscriptions(ctx, listSubscriptionsRequest)
}
func (server *subscriberServiceServerImpl) DeleteSubscription(ctx context.Context, deleteSubscriptionRequest *api_pb.DeleteSubscriptionRequest) (*empty.Empty, error) {
	return server.rawClient.DeleteSubscription(ctx, deleteSubscriptionRequest)
}
func (server *subscriberServiceServerImpl) ModifyAckDeadline(ctx context.Context, modifyAckDeadlineRequest *api_pb.ModifyAckDeadlineRequest) (*empty.Empty, error) {
	return server.rawClient.ModifyAckDeadline(ctx, modifyAckDeadlineRequest)
}
func (server *subscriberServiceServerImpl) Acknowledge(ctx context.Context, acknowledgeRequest *api_pb.AcknowledgeRequest) (*empty.Empty, error) {
	return server.rawClient.Acknowledge(ctx, acknowledgeRequest)
}
func (server *subscriberServiceServerImpl) Pull(ctx context.Context, pullRequest *api_pb.PullRequest) (*api_pb.PullResponse, error) {
	return server.rawClient.Pull(ctx, pullRequest)
}
func (server *subscriberServiceServerImpl) StreamingPull(stream api_pb.Subscriber_StreamingPullServer) error {
	streamingPullClient, err := server.rawClient.StreamingPull(stream.Context())
	if err != nil {
		return err
	}

	cctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	errCh1 := make(chan error)
	errCh2 := make(chan error)
	// Client -> Real Cloud Pub/Sub
	go func(ctx context.Context, errCh chan<- error, streamingPullServer api_pb.Subscriber_StreamingPullServer, streamingPullClient api_pb.Subscriber_StreamingPullClient) {
		for {
			select {
			case <-ctx.Done():
				break
			default:
				received, err := streamingPullServer.Recv()
				if err != nil {
					errCh <- err
					break
				}
				if err = streamingPullClient.Send(received); err != nil {
					errCh <- err
					break
				}
			}
		}
	}(cctx, errCh1, stream, streamingPullClient)

	// Real Cloud Pub/Sub -> Client
	go func(ctx context.Context, errCh chan<- error, streamingPullServer api_pb.Subscriber_StreamingPullServer, streamingPullClient api_pb.Subscriber_StreamingPullClient) {
		for {
			select {
			case <-ctx.Done():
				break
			default:
				received, err := streamingPullClient.Recv()
				if err != nil {
					errCh <- err
					break
				}
				if err = streamingPullServer.Send(received); err != nil {
					errCh <- err
					break
				}
			}
		}
	}(cctx, errCh2, stream, streamingPullClient)

	select {
	case err := <-errCh1:
		if err == io.EOF {
			return streamingPullClient.CloseSend()
		} else {
			return err
		}
	case err := <-errCh2:
		return err
	}
}
func (server *subscriberServiceServerImpl) ModifyPushConfig(ctx context.Context, modifyPushConfigRequest *api_pb.ModifyPushConfigRequest) (*empty.Empty, error) {
	return server.rawClient.ModifyPushConfig(ctx, modifyPushConfigRequest)
}
func (server *subscriberServiceServerImpl) GetSnapshot(ctx context.Context, getSnapshotRequest *api_pb.GetSnapshotRequest) (*api_pb.Snapshot, error) {
	return server.rawClient.GetSnapshot(ctx, getSnapshotRequest)
}
func (server *subscriberServiceServerImpl) ListSnapshots(ctx context.Context, listSnapshotsRequest *api_pb.ListSnapshotsRequest) (*api_pb.ListSnapshotsResponse, error) {
	return server.rawClient.ListSnapshots(ctx, listSnapshotsRequest)
}
func (server *subscriberServiceServerImpl) CreateSnapshot(ctx context.Context, createSnapshotRequest *api_pb.CreateSnapshotRequest) (*api_pb.Snapshot, error) {
	return server.rawClient.CreateSnapshot(ctx, createSnapshotRequest)
}
func (server *subscriberServiceServerImpl) UpdateSnapshot(ctx context.Context, updateSnapshotRequest *api_pb.UpdateSnapshotRequest) (*api_pb.Snapshot, error) {
	return server.rawClient.UpdateSnapshot(ctx, updateSnapshotRequest)
}
func (server *subscriberServiceServerImpl) DeleteSnapshot(ctx context.Context, deleteSnapshotRequest *api_pb.DeleteSnapshotRequest) (*empty.Empty, error) {
	return server.rawClient.DeleteSnapshot(ctx, deleteSnapshotRequest)
}
func (server *subscriberServiceServerImpl) Seek(ctx context.Context, seekRequest *api_pb.SeekRequest) (*api_pb.SeekResponse, error) {
	return server.rawClient.Seek(ctx, seekRequest)
}
