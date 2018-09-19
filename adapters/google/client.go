package google

import (
	"cloud.google.com/go/pubsub"
	"context"
)

type client interface {
	Topic(name string) *pubsub.Topic
	TopicExists(ctx context.Context, topic *pubsub.Topic) (bool, error)
	CreateTopic(ctx context.Context, name string) (*pubsub.Topic, error)
	Subscription(name string) *pubsub.Subscription
	SubscriptionExists(ctx context.Context, subscription *pubsub.Subscription) (bool, error)
	CreateSubscription(ctx context.Context, name string, config pubsub.SubscriptionConfig) (*pubsub.Subscription, error)
}
