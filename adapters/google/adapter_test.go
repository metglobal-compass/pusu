package google

import (
	"errors"
	"github.com/metglobal-compass/pusu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAdapter_CreateSubscriptionValidations(t *testing.T) {
	adapter := new(Adapter)

	// Must return error when topic name is empty
	subscriptionWithEmptyTopic := pusu.NewSubscription("", "testing", new(dummySubscriber))
	err := adapter.CreateSubscription(subscriptionWithEmptyTopic)
	if err == nil {
		t.Error("Error while validation of subscription. Topic must be validated.")
	}

	// Must return error when subscription name is empty
	subscriptionWithEmptySubscriberName := pusu.NewSubscription("test", "", new(dummySubscriber))
	err = adapter.CreateSubscription(subscriptionWithEmptySubscriberName)
	if err == nil {
		t.Error("Error while validation of subscription. Subscription name must be validated.")
	}

	// Must return error when subscription handler is empty
	subscriptionWithEmptySubscriber := pusu.NewSubscription("test", "testing", nil)
	err = adapter.CreateSubscription(subscriptionWithEmptySubscriber)
	if err == nil {
		t.Error("Error while validation of subscription. Subscription handler must be validated")
	}
}

func TestAdapter_CreateSubscription(t *testing.T) {
	// Create mocked object
	successCreator := new(fakeCreator)
	successCreator.On("CreateSubscription", mock.Anything).Return(nil)

	// Create real instance and call real method
	adapter := new(Adapter)
	adapter.cloudAdder = successCreator
	adapter.httpHandlerAdder = successCreator
	err := adapter.CreateSubscription(pusu.NewSubscription("test", "testing", new(dummySubscriber)))

	// There must be no error
	assert.Nil(t, err)

	// We used same interface so our successCreator must ve called twice for both cloud adder and http handler adder
	successCreator.AssertNumberOfCalls(t, "CreateSubscription", 2)
}

func TestAdapter_CreateSubscriptionOnCloudFailure(t *testing.T) {
	// Create successor mocked object
	successCreator := new(fakeCreator)
	successCreator.On("CreateSubscription", mock.Anything).Return(nil)

	// Create failure mocked object
	failCreator := new(fakeCreator)
	failCreator.On("CreateSubscription", mock.Anything).Return(errors.New("error... "))

	// Create real instance and call real method
	adapter := new(Adapter)
	adapter.cloudAdder = successCreator
	adapter.httpHandlerAdder = failCreator

	// Subscription creation must return error when cloud adder return error
	err := adapter.CreateSubscription(pusu.NewSubscription("test", "testing", new(dummySubscriber)))
	assert.Error(t, err)
}

func TestAdapter_CreateSubscriptionOnHttpHandlerFailure(t *testing.T) {
	// Create successor mocked object
	successCreator := new(fakeCreator)
	successCreator.On("CreateSubscription", mock.Anything).Return(nil)

	// Create failure mocked object
	failCreator := new(fakeCreator)
	failCreator.On("CreateSubscription", mock.Anything).Return(errors.New("error... "))

	// Create real instance and call real method
	adapter := new(Adapter)
	adapter.cloudAdder = failCreator
	adapter.httpHandlerAdder = successCreator

	// Subscription creation must return error when cloud adder return error
	err := adapter.CreateSubscription(pusu.NewSubscription("test", "testing", new(dummySubscriber)))
	assert.Error(t, err)
}

func TestAdapter_Run(t *testing.T) {
	// Create successor mocked object
	successRunner := new(fakeRunner)
	successRunner.On("Run", mock.Anything).Return(nil)

	// Create real object and call real method
	adapter := new(Adapter)
	adapter.runner = successRunner
	err := adapter.Run(pusu.NewSubscription("test", "testing", new(dummySubscriber)))

	// Error must be nil
	assert.Nil(t, err)
}

func TestAdapter_RunError(t *testing.T) {
	// Create successor mocked object
	failRunner := new(fakeRunner)
	failRunner.On("Run", mock.Anything).Return(errors.New("error"))

	// Create real object and call real method
	adapter := new(Adapter)
	adapter.runner = failRunner
	err := adapter.Run(pusu.NewSubscription("test", "testing", new(dummySubscriber)))

	// Error must not be nil
	assert.Error(t, err)
}

func TestAdapter_CreateAdapter(t *testing.T) {
	// Call real method and check each interfaces have proper types
	adapter, err := CreateAdapter("my-project", "http://localhost")

	assert.Nil(t, err)
	assert.NotNil(t, adapter.httpHandlerAdder)
	assert.NotNil(t, adapter.cloudAdder)
	assert.NotNil(t, adapter.runner)

	_, httpHandlerProper := adapter.httpHandlerAdder.(*httpHandlerAdder)
	assert.True(t, httpHandlerProper, "Http handler is not proper")

	cloudAdder, cloudAdderProper := adapter.cloudAdder.(*cloudAdder)
	assert.True(t, cloudAdderProper, "Cloud adder is not proper")

	// Check cloud adder created with proper host configuration
	assert.Equal(t, "http://localhost", cloudAdder.host)

	_, pubSubClientProper := cloudAdder.client.(*pubSubClientWrapper)
	assert.True(t, pubSubClientProper, "Pubsub client is not proper")

	_, appEngineRunnerProper := adapter.runner.(*appEngineRunner)
	assert.True(t, appEngineRunnerProper, "Runner  is not proper")
}

func TestAdapter_CreateAdapterErrorWithEmptyProject(t *testing.T) {
	// Call real method and check each interfaces have proper type
	_, err := CreateAdapter("", "http://localhost")
	assert.Error(t, err)
}

func TestAdapter_CreateAdapterErrorWithEmptyHost(t *testing.T) {
	// Call real method and check each interfaces have proper type
	_, err := CreateAdapter("my-project", "")
	assert.Error(t, err)
}

type fakeCreator struct {
	mock.Mock
}

func (f *fakeCreator) CreateSubscription(subscription *pusu.Subscription) error {
	args := f.Called(subscription)
	return args.Error(0)
}

type fakeRunner struct {
	mock.Mock
}

func (f *fakeRunner) Run(subscription *pusu.Subscription) error {
	args := f.Called(subscription)
	return args.Error(0)
}

type dummySubscriber struct {
}

func (d *dummySubscriber) Handle(m *pusu.Message) error {
	return nil
}

type failureSubscriber struct {
}

func (d *failureSubscriber) Handle(m *pusu.Message) error {
	return errors.New("Message error. ")
}
