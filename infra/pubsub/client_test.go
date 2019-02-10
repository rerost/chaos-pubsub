package pubsub_test

import (
	"context"
	"testing"

	"github.com/rerost/chaos-pubsub/infra/pubsub"
)

func TestNewClient(t *testing.T) {
	ctx := context.Background()
	_, err := pubsub.NewClient(ctx, "TEST", []byte{})
	t.Error(err)
}
