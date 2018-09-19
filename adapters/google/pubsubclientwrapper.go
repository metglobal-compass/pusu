package google

import (
	"cloud.google.com/go/pubsub"
	"context"
)

// Google Cloud Pub/Sub Client wrapper which implements client interface
type pubSubClientWrapper struct {
	client *pubsub.Client
}

func (p *pubSubClientWrapper) Topic(name string) *pubsub.Topic {
	return p.client.Topic(name)
}

func (p *pubSubClientWrapper) TopicExists(ctx context.Context, topic *pubsub.Topic) (bool, error) {
	return topic.Exists(ctx)
}

func (p *pubSubClientWrapper) CreateTopic(ctx context.Context, name string) (*pubsub.Topic, error) {
	return p.client.CreateTopic(ctx, name)
}

func (p *pubSubClientWrapper) Subscription(name string) *pubsub.Subscription {
	return p.client.Subscription(name)
}

func (p *pubSubClientWrapper) SubscriptionExists(ctx context.Context, subscription *pubsub.Subscription) (bool, error) {
	return subscription.Exists(ctx)
}

func (p *pubSubClientWrapper) CreateSubscription(ctx context.Context, name string, config pubsub.SubscriptionConfig) (*pubsub.Subscription, error) {
	return p.client.CreateSubscription(ctx, name, config)
}
