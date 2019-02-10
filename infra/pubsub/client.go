package pubsub

import (
	"context"
	"reflect"
	"unsafe"

	"cloud.google.com/go/pubsub"
	"github.com/srvc/fail"
	"google.golang.org/api/option"
	api_pb "google.golang.org/genproto/googleapis/pubsub/v1"
	"google.golang.org/grpc"
)

type Client struct {
	PublisherClient  api_pb.PublisherClient
	SubscriberClient api_pb.SubscriberClient
}

func NewClient(ctx context.Context, project string, keyfile []byte) (Client, error) {
	var pubsubClient *pubsub.Client
	var err error
	if len(keyfile) == 0 {
		pubsubClient, err = pubsub.NewClient(ctx, project)
	} else {
		pubsubClient, err = pubsub.NewClient(ctx, project, option.WithCredentialsJSON(keyfile))
	}

	if err != nil {
		return Client{}, fail.Wrap(err)
	}

	// TODO(rerost) **REMOVE EVIL**
	v := reflect.ValueOf(*pubsubClient)
	pubc := reflect.Indirect(v.FieldByName("pubc"))
	subc := reflect.Indirect(v.FieldByName("subc"))
	publisherClient := reflect.Indirect(pubc.FieldByName("conn"))
	subscriberClient := reflect.Indirect(subc.FieldByName("conn"))
	publisherClientConn := reflect.NewAt(publisherClient.Type(), unsafe.Pointer(publisherClient.UnsafeAddr())).Interface().(*grpc.ClientConn)
	subscriberClientConn := reflect.NewAt(publisherClient.Type(), unsafe.Pointer(subscriberClient.UnsafeAddr())).Interface().(*grpc.ClientConn)

	return Client{
		PublisherClient:  api_pb.NewPublisherClient(publisherClientConn),
		SubscriberClient: api_pb.NewSubscriberClient(subscriberClientConn),
	}, nil
}
