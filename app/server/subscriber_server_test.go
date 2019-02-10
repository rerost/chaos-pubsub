package server_test

import (
	"context"
	"net"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"github.com/google/go-cmp/cmp"
	"github.com/rerost/chaos-pubsub/app/server"
	"github.com/srvc/fail"
	"google.golang.org/api/option"
	api_pb "google.golang.org/genproto/googleapis/pubsub/v1"
	"google.golang.org/grpc"
)

func TestPubsubSubscribe(t *testing.T) {
	ctx := context.Background()
	srvr := pstest.NewServer()

	conn, err := grpc.Dial(srvr.Addr, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	topicName := "projects/test/topics/test-topic"
	subscriptionName := "projects/test/subscriptions/test-sub"

	subscriberClient := api_pb.NewSubscriberClient(conn)
	publisherClient := api_pb.NewPublisherClient(conn)

	publishServer := server.NewPublisherServiceServer(publisherClient)
	server := server.NewSubscriberServiceServer(subscriberClient)

	// Prepare topic
	_, err = publishServer.CreateTopic(ctx, &api_pb.Topic{Name: topicName})
	if err != nil {
		t.Error(err)
	}

	{
		subscription := api_pb.Subscription{Topic: topicName, Name: subscriptionName, AckDeadlineSeconds: 30}
		_, err = server.CreateSubscription(ctx, &subscription)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestPubsubSubscriberWithRealClient(t *testing.T) {
	ctx := context.Background()
	port := "5000"
	go func(t *testing.T) {
		if err := startRealServer(":" + port); err != nil {
			t.Error(t)
		}
	}(t)

	conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	pubsubClient, err := pubsub.NewClient(ctx, "test1", option.WithGRPCConn(conn))
	if err != nil {
		t.Error(t)
		t.Fail()
	}

	topicName := "test-topic"
	subscriptionName := "test-sub"

	topic, err := pubsubClient.CreateTopic(ctx, topicName)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	sub, err := pubsubClient.CreateSubscription(ctx, subscriptionName, pubsub.SubscriptionConfig{Topic: topic})
	if err != nil {
		t.Error(err)
	}

	{
		topic.Publish(ctx, &pubsub.Message{Data: []byte("test message")})
		cctx, _ := context.WithTimeout(ctx, time.Second*5)
		err := sub.Receive(cctx, func(ctx context.Context, message *pubsub.Message) {
			message.Ack()
			if diff := cmp.Diff(message.Data, []byte("test message")); diff != "" {
				t.Error(diff)
			}
		})

		if err != nil {
			t.Error(err)
		}
	}
}

func startRealServer(port string) error {
	srvr := pstest.NewServer()
	conn, err := grpc.Dial(srvr.Addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	subscriberClient := api_pb.NewSubscriberClient(conn)
	publisherClient := api_pb.NewPublisherClient(conn)

	publishServer := server.NewPublisherServiceServer(publisherClient)
	subscribeServer := server.NewSubscriberServiceServer(subscriberClient)

	network, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	gserver := grpc.NewServer()
	publishServer.RegisterWithServer(gserver)
	subscribeServer.RegisterWithServer(gserver)

	errCh := make(chan error)
	go func() {
		err := gserver.Serve(network)
		if err != nil {
			errCh <- fail.Wrap(err)
		}
	}()

	err = <-errCh
	if err != nil {
		return err
	}
	return nil
}
