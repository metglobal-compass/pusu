package google

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/metglobal-compass/pusu"
	"time"
)

const (
	endpointPattern = "%s/_handlers/topics/%s/subscribers/%s"
)

type cloudAdder struct {
	client client
	host   string
}

// Implementation of internal Creator interface for Google Adapter
func (t *cloudAdder) CreateSubscription(subscription pusu.Subscription) error {
	// Use single context
	ctx := context.Background()

	// Get topic instance
	topic := t.client.Topic(subscription.Topic())

	// Check if topic existss
	topicExists, err := t.client.TopicExists(ctx, topic)
	if err != nil {
		return err
	}

	// If topic does not exists, create in cloud
	if !topicExists {
		topic, err = t.client.CreateTopic(ctx, subscription.Topic())
		if err != nil {
			return err
		}
	}

	// Create subscription instance
	clientSubscription := t.client.Subscription(subscription.Name())

	// Check if subscription exists
	exists, err := t.client.SubscriptionExists(ctx, clientSubscription)
	if err != nil {
		return err
	}

	// If subscription does not exists, create it in cloud
	if !exists {
		subscriptionConfig := pubsub.SubscriptionConfig{
			Topic:       topic,
			AckDeadline: 10 * time.Second,
			PushConfig: pubsub.PushConfig{
				Endpoint: fmt.Sprintf(endpointPattern, t.host, subscription.Topic(), subscription.Name()),
			},
		}
		_, err := t.client.CreateSubscription(ctx, subscription.Name(), subscriptionConfig)
		if err != nil {
			return err
		}
	}

	return nil
}
