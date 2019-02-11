package fault_test

import (
	"context"
	"net"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"github.com/google/go-cmp/cmp"
	"github.com/rerost/chaos-pubsub/app/server"
	"github.com/rerost/chaos-pubsub/lib/interceptor/fault"
	"google.golang.org/api/option"
	pb "google.golang.org/genproto/googleapis/pubsub/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func TestDuplicateSendPublish(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	info := &grpc.UnaryServerInfo{FullMethod: "/google.pubsub.v1.Publisher/Publish"}

	called := 0
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		// Send method
		called++
		return nil, nil
	}

	fault.UnaryServerInterceptor()(ctx, nil, info, handler)
	if diff := cmp.Diff(called, 2); diff != "" {
		t.Error(diff)
	}
}

func createServer(pubsubServer *pstest.Server, serverOption ...grpc.ServerOption) (*grpc.Server, error) {
	// Prepare Server
	conn, err := grpc.Dial(pubsubServer.Addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	publisherBackend := pb.NewPublisherClient(conn)
	subscriberBackend := pb.NewSubscriberClient(conn)

	srv := grpc.NewServer(
		serverOption...,
	)
	publisherServer := server.NewPublisherServiceServer(publisherBackend)
	subscriberServer := server.NewSubscriberServiceServer(subscriberBackend)
	publisherServer.RegisterWithServer(srv)
	subscriberServer.RegisterWithServer(srv)

	return srv, nil
}

func preparePubSub(ctx context.Context, client *pubsub.Client) (*pubsub.Topic, *pubsub.Subscription, error) {
	var topic *pubsub.Topic
	var subscription *pubsub.Subscription
	var err error

	topic, err = client.CreateTopic(ctx, "test-topic")
	if err != nil {
		return topic, subscription, err
	}

	subscription, err = client.CreateSubscription(ctx, "test-subscription", pubsub.SubscriptionConfig{Topic: topic})
	if err != nil {
		return topic, subscription, err
	}

	return topic, subscription, nil
}

func TestDuplicateSendPublishWithRealPubsub(t *testing.T) {
	t.Parallel()

	pubsubServer := pstest.NewServer()
	defer pubsubServer.Close()

	srv, err := createServer(pubsubServer, grpc.UnaryInterceptor(fault.UnaryServerInterceptor()))
	if err != nil {
		t.Error(err)
	}
	defer srv.Stop()

	lis := bufconn.Listen(1024 * 1024)
	go func(t *testing.T) {
		if err := srv.Serve(lis); err != nil {
			t.Error(err)
		}
	}(t)

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(func(string, time.Duration) (net.Conn, error) { return lis.Dial() }), grpc.WithInsecure())
	pubsubClient, err := pubsub.NewClient(ctx, "project", option.WithGRPCConn(conn))
	if err != nil {
		t.Error(err)
	}

	topic, subscription, err := preparePubSub(ctx, pubsubClient)
	if err != nil {
		t.Error(err)
	}

	topic.Publish(ctx, &pubsub.Message{Data: []byte("test message")})

	countCh := make(chan int, 2)
	go func(ctx context.Context, countCh chan<- int) {
		cctx, cancel := context.WithTimeout(ctx, time.Second*2)
		defer cancel()

		subscription.Receive(cctx, func(ctx context.Context, m *pubsub.Message) {
			m.Ack()

			countCh <- 1
		})

		close(countCh)
	}(ctx, countCh)

	count := 0
	for {
		_, ok := <-countCh
		if !ok {
			if diff := cmp.Diff(count, 2); diff != "" {
				t.Error(diff)
			}
			return
		} else {
			count++
		}
	}
}

func TestDuplicateSendSubscriptionWithRealPubsub(t *testing.T) {
	t.Parallel()
	pubsubServer := pstest.NewServer()
	defer pubsubServer.Close()

	srv, err := createServer(pubsubServer, grpc.StreamInterceptor(fault.StreamServerInterceptor()))
	if err != nil {
		t.Error(err)
	}
	defer srv.Stop()

	lis := bufconn.Listen(1024 * 1024)
	go func(t *testing.T) {
		if err := srv.Serve(lis); err != nil {
			t.Error(err)
		}
	}(t)

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(func(string, time.Duration) (net.Conn, error) { return lis.Dial() }), grpc.WithInsecure())
	pubsubClient, err := pubsub.NewClient(ctx, "project", option.WithGRPCConn(conn))
	if err != nil {
		t.Error(err)
	}

	topic, subscription, err := preparePubSub(ctx, pubsubClient)
	if err != nil {
		t.Error(err)
	}

	topic.Publish(ctx, &pubsub.Message{Data: []byte("test message")})
	countCh := make(chan int, 2)
	go func(ctx context.Context, countCh chan<- int) {
		cctx, cancel := context.WithTimeout(ctx, time.Second*2)
		defer cancel()

		subscription.Receive(cctx, func(ctx context.Context, m *pubsub.Message) {
			m.Ack()

			countCh <- 1
		})

		close(countCh)
	}(ctx, countCh)

	count := 0
	for {
		_, ok := <-countCh
		if !ok {
			if diff := cmp.Diff(count, 2); diff != "" {
				t.Error(diff)
			}
			return
		} else {
			count++
		}
	}
}
