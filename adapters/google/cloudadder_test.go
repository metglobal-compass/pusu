package google

import (
	"cloud.google.com/go/pubsub"
	"context"
	"errors"
	"fmt"
	"github.com/metglobal-compass/pusu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCloudAdder_CreateSubscription(t *testing.T) {
	// Create fake mocked client In this case, we try to create a subscription which does not exist in cloud
	fakeClient := new(fakeClient)
	fakeClient.On("Topic", "test").Return(&pubsub.Topic{})
	fakeClient.On("TopicExists", context.Background(), &pubsub.Topic{}).Return(false, nil)
	fakeClient.On("CreateTopic", context.Background(), "test").Return(&pubsub.Topic{}, nil)
	fakeClient.On("Subscription", "testing").Return(&pubsub.Subscription{})
	fakeClient.On("SubscriptionExists", context.Background(), &pubsub.Subscription{}).
		Return(false, nil)
	fakeClient.On("CreateSubscription", context.Background(), "testing", mock.Anything).
		Return(&pubsub.Subscription{}, nil)

	// Call real method
	cloudAdder := new(cloudAdder)
	cloudAdder.client = fakeClient
	cloudAdder.host = "http://localhost"
	err := cloudAdder.CreateSubscription(pusu.NewSubscription("test", "testing", new(dummySubscriber)))

	// Everything is fine, so real method must return nil as error
	assert.Nil(t, err, "cloudAdder.CreateSubscription must return nil when everything is fine.")

	// Check if CreateSubscription is called with proper configuration
	expectedPushConfiguration := pubsub.SubscriptionConfig{
		Topic:       &pubsub.Topic{},
		AckDeadline: 10 * time.Second,
		PushConfig: pubsub.PushConfig{
			Endpoint: fmt.Sprintf(endpointPattern, cloudAdder.host, "test", "testing"),
		},
	}
	fakeClient.AssertCalled(t, "CreateSubscription", mock.Anything, "testing", expectedPushConfiguration)

	// All methods must be called exactly once when subscription does not exists in cloud
	fakeClient.AssertNumberOfCalls(t, "Topic", 1)
	fakeClient.AssertNumberOfCalls(t, "TopicExists", 1)
	fakeClient.AssertNumberOfCalls(t, "CreateTopic", 1)
	fakeClient.AssertNumberOfCalls(t, "Subscription", 1)
	fakeClient.AssertNumberOfCalls(t, "SubscriptionExists", 1)
	fakeClient.AssertNumberOfCalls(t, "CreateSubscription", 1)
	fakeClient.AssertExpectations(t)
}

func TestCloudAdder_CreateSubscriptionErrorOnTopicExists(t *testing.T) {
	// Create fake mocked client In this case, we get an error while checking a topic's existence
	fakeClient := new(fakeClient)
	fakeClient.On("Topic", "test").Return(&pubsub.Topic{})
	// Got an error while checking a topic's existence
	fakeClient.On("TopicExists", context.Background(), &pubsub.Topic{}).
		Return(false, errors.New("error"))

	// Call real method
	cloudAdder := new(cloudAdder)
	cloudAdder.client = fakeClient
	cloudAdder.host = "http://localhost"
	err := cloudAdder.CreateSubscription(pusu.NewSubscription("test", "testing", new(dummySubscriber)))

	// Got an error
	assert.Error(t, err)

	// Following methods must be called exactly once
	fakeClient.AssertNumberOfCalls(t, "Topic", 1)
	fakeClient.AssertNumberOfCalls(t, "TopicExists", 1)

	// If got an error while checking topic exists, cloudAdder must return error and must not call following methods
	fakeClient.AssertNotCalled(t, "CreateTopic")
	fakeClient.AssertNotCalled(t, "Subscription")
	fakeClient.AssertNotCalled(t, "SubscriptionExists")
	fakeClient.AssertNotCalled(t, "CreateSubscription")
	fakeClient.AssertExpectations(t)
}

func TestCloudAdder_CreateSubscriptionErrorOnCreateTopic(t *testing.T) {
	// Create fake mocked client In this case, we get an error while creating topic
	fakeClient := new(fakeClient)
	fakeClient.On("Topic", "test").Return(&pubsub.Topic{})
	fakeClient.On("TopicExists", context.Background(), &pubsub.Topic{}).Return(false, nil)
	// Got an error while creating topic
	fakeClient.On("CreateTopic", context.Background(), "test").
		Return(&pubsub.Topic{}, errors.New("error"))

	// Call real method
	cloudAdder := new(cloudAdder)
	cloudAdder.client = fakeClient
	cloudAdder.host = "http://localhost"
	err := cloudAdder.CreateSubscription(pusu.NewSubscription("test", "testing", new(dummySubscriber)))

	// Got an error
	assert.Error(t, err)

	// Following methods must be called exactly once
	fakeClient.AssertNumberOfCalls(t, "Topic", 1)
	fakeClient.AssertNumberOfCalls(t, "TopicExists", 1)
	fakeClient.AssertNumberOfCalls(t, "CreateTopic", 1)

	// If got an error while creating topic, cloudAdder must return error and must not call following methods
	fakeClient.AssertNotCalled(t, "Subscription")
	fakeClient.AssertNotCalled(t, "SubscriptionExists")
	fakeClient.AssertNotCalled(t, "CreateSubscription")
	fakeClient.AssertExpectations(t)
}

func TestCloudAdder_CreateSubscriptionErrorOnSubscriptionExists(t *testing.T) {
	// Create fake mocked client In this case, we get an error while checking a subscription's existence
	fakeClient := new(fakeClient)
	fakeClient.On("Topic", "test").Return(&pubsub.Topic{})
	fakeClient.On("TopicExists", context.Background(), &pubsub.Topic{}).Return(false, nil)
	fakeClient.On("CreateTopic", context.Background(), "test").Return(&pubsub.Topic{}, nil)
	fakeClient.On("Subscription", "testing").Return(&pubsub.Subscription{})
	// We got an error while checking subscription's existence
	fakeClient.On("SubscriptionExists", context.Background(), &pubsub.Subscription{}).
		Return(false, errors.New("error"))

	// Call real method
	cloudAdder := new(cloudAdder)
	cloudAdder.client = fakeClient
	cloudAdder.host = "http://localhost"
	err := cloudAdder.CreateSubscription(pusu.NewSubscription("test", "testing", new(dummySubscriber)))

	// Got an error
	assert.Error(t, err)

	// Following methods must be called exactly once
	fakeClient.AssertNumberOfCalls(t, "Topic", 1)
	fakeClient.AssertNumberOfCalls(t, "TopicExists", 1)
	fakeClient.AssertNumberOfCalls(t, "CreateTopic", 1)
	fakeClient.AssertNumberOfCalls(t, "Subscription", 1)
	fakeClient.AssertNumberOfCalls(t, "SubscriptionExists", 1)

	// If got an error while checking subscription, cloudAdder must return error and must not call following methods
	fakeClient.AssertNotCalled(t, "CreateSubscription")
	fakeClient.AssertExpectations(t)
}

func TestCloudAdder_CreateSubscriptionErrorOnCreateSubscription(t *testing.T) {
	// Create fake mocked client In this case, we got an error while creating subscription on last step
	fakeClient := new(fakeClient)
	fakeClient.On("Topic", "test").Return(&pubsub.Topic{})
	fakeClient.On("TopicExists", context.Background(), &pubsub.Topic{}).Return(false, nil)
	fakeClient.On("CreateTopic", context.Background(), "test").Return(&pubsub.Topic{}, nil)
	fakeClient.On("Subscription", "testing").Return(&pubsub.Subscription{})
	fakeClient.On("SubscriptionExists", context.Background(), &pubsub.Subscription{}).
		Return(false, nil)
	// We got an error while creating subscription
	fakeClient.On("CreateSubscription", context.Background(), "testing", mock.Anything).
		Return(&pubsub.Subscription{}, errors.New("error"))

	// Call real method
	cloudAdder := new(cloudAdder)
	cloudAdder.client = fakeClient
	cloudAdder.host = "http://localhost"
	err := cloudAdder.CreateSubscription(pusu.NewSubscription("test", "testing", new(dummySubscriber)))

	// Everything is fine, so real method must return nil as error
	assert.Error(t, err)

	// All methods must be called exactly once
	fakeClient.AssertNumberOfCalls(t, "Topic", 1)
	fakeClient.AssertNumberOfCalls(t, "TopicExists", 1)
	fakeClient.AssertNumberOfCalls(t, "CreateTopic", 1)
	fakeClient.AssertNumberOfCalls(t, "Subscription", 1)
	fakeClient.AssertNumberOfCalls(t, "SubscriptionExists", 1)
	fakeClient.AssertNumberOfCalls(t, "CreateSubscription", 1)
	fakeClient.AssertExpectations(t)
}

type fakeClient struct {
	mock.Mock
}

func (f *fakeClient) Topic(name string) *pubsub.Topic {
	args := f.Called(name)
	return args.Get(0).(*pubsub.Topic)
}

func (f *fakeClient) TopicExists(ctx context.Context, topic *pubsub.Topic) (bool, error) {
	args := f.Called(ctx, topic)
	return args.Bool(0), args.Error(1)
}

func (f *fakeClient) CreateTopic(ctx context.Context, name string) (*pubsub.Topic, error) {
	args := f.Called(ctx, name)
	return args.Get(0).(*pubsub.Topic), args.Error(1)
}

func (f *fakeClient) Subscription(name string) *pubsub.Subscription {
	args := f.Called(name)
	return args.Get(0).(*pubsub.Subscription)
}

func (f *fakeClient) SubscriptionExists(ctx context.Context, subscription *pubsub.Subscription) (bool, error) {
	args := f.Called(ctx, subscription)
	return args.Bool(0), args.Error(1)
}

func (f *fakeClient) CreateSubscription(ctx context.Context, name string, config pubsub.SubscriptionConfig) (*pubsub.Subscription, error) {
	args := f.Called(ctx, name, config)
	return args.Get(0).(*pubsub.Subscription), args.Error(1)
}
