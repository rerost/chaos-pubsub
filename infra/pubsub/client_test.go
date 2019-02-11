package pubsub_test

import (
	"context"
	"runtime"
	"testing"

	real_pubsub "cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"github.com/google/go-cmp/cmp"
	"github.com/rerost/chaos-pubsub/infra/pubsub"
	"google.golang.org/api/option"
	api_pb "google.golang.org/genproto/googleapis/pubsub/v1"
	"google.golang.org/grpc"
)

func TestNewClientWithGC(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	s := pstest.NewServer()
	defer s.Close()

	conn, err := grpc.Dial(s.Addr, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}

	// Prepare topic/subscription
	{
		pubsubClient, err := real_pubsub.NewClient(ctx, "dummy", option.WithGRPCConn(conn))
		if err != nil {
			t.Error(err)
		}

		topic, err := pubsubClient.CreateTopic(ctx, "dummy-topic")
		if err != nil {
			t.Error(err)
		}

		_, err = pubsubClient.CreateSubscription(ctx, "dummy-subscription", real_pubsub.SubscriptionConfig{Topic: topic})
		if err != nil {
			t.Error(err)
		}
	}

	client, err := pubsub.NewClient(ctx, "dummy", option.WithGRPCConn(conn))
	if err != nil {
		t.Error(err)
	}

	pc := client.PublisherClient
	sc := client.SubscriberClient

	presult, err := pc.ListTopics(ctx, &api_pb.ListTopicsRequest{Project: "projects/dummy"})
	if err != nil {
		t.Error(err)
	}
	if presult == nil {
		t.Error("Want not nil")
	}

	sresult, err := sc.ListSubscriptions(ctx, &api_pb.ListSubscriptionsRequest{Project: "projects/dummy"})
	if err != nil {
		t.Error(err)
	}
	if sresult == nil {
		t.Error("Want not nil")
	}

	// Run GC for check reflection
	runtime.GC()

	presultAfterGC, err := pc.ListTopics(ctx, &api_pb.ListTopicsRequest{Project: "projects/dummy"})
	if err != nil {
		t.Error(err)
	}

	sresultAfterGC, err := sc.ListSubscriptions(ctx, &api_pb.ListSubscriptionsRequest{Project: "projects/dummy"})
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(presult, presultAfterGC); diff != "" {
		t.Error(diff)
	}

	if diff := cmp.Diff(sresult, sresultAfterGC); diff != "" {
		t.Error(diff)
	}
}
