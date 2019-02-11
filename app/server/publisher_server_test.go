package server_test

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub/pstest"
	"github.com/google/go-cmp/cmp"
	"github.com/rerost/chaos-pubsub/app/server"
	api_pb "google.golang.org/genproto/googleapis/pubsub/v1"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc"
)

func TestPubsubPublish(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	srvr := pstest.NewServer()

	conn, err := grpc.Dial(srvr.Addr, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	topicName := "projects/test/topics/test-topic"

	publisherClient := api_pb.NewPublisherClient(conn)

	server := server.NewPublisherServiceServer(publisherClient)
	{
		result, err := server.CreateTopic(ctx, &api_pb.Topic{Name: topicName})
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(result, &api_pb.Topic{Name: topicName}); diff != "" {
			t.Error(diff)
		}
	}

	{
		result, err := server.ListTopics(ctx, &api_pb.ListTopicsRequest{})
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(result, &api_pb.ListTopicsResponse{Topics: []*api_pb.Topic{{Name: topicName}}}); diff != "" {
			t.Error(diff)
		}
	}

	{
		result, err := server.GetTopic(ctx, &api_pb.GetTopicRequest{Topic: topicName})
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(result, &api_pb.Topic{Name: topicName}); diff != "" {
			t.Error(diff)
		}
	}

	{
		updateTopic := api_pb.Topic{Name: topicName, Labels: map[string]string{"test": "test"}}
		_, err := server.UpdateTopic(ctx, &api_pb.UpdateTopicRequest{Topic: &updateTopic, UpdateMask: &field_mask.FieldMask{Paths: []string{"labels"}}})
		if err != nil {
			t.Error(err)
		}
	}

	{
		message := api_pb.PubsubMessage{
			Data: []byte("test message"),
		}
		result, err := server.Publish(ctx, &api_pb.PublishRequest{Topic: topicName, Messages: []*api_pb.PubsubMessage{&message}})
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(result, &api_pb.PublishResponse{MessageIds: []string{"m0"}}); diff != "" {
			t.Error(diff)
		}
	}
}
