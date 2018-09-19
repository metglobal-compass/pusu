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
	subscriptionWithEmptyTopic := pusu.NewSubscription("", "testing", func(m *pusu.Message) error {
		return nil
	})
	err := adapter.CreateSubscription(subscriptionWithEmptyTopic)
	if err == nil {
		t.Error("Error while validation of subscription. Topic must be validated.")
	}

	// Must return error when subscription name is empty
	subscriptionWithEmptySubscriberName := pusu.NewSubscription("test", "", func(m *pusu.Message) error {
		return nil
	})
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
	err := adapter.CreateSubscription(pusu.NewSubscription("test", "testing", func(m *pusu.Message) error {
		return nil
	}))

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
	err := adapter.CreateSubscription(pusu.NewSubscription("test", "testing", func(m *pusu.Message) error {
		return nil
	}))
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
	err := adapter.CreateSubscription(pusu.NewSubscription("test", "testing", func(m *pusu.Message) error {
		return nil
	}))
	assert.Error(t, err)
}

func TestAdapter_Run(t *testing.T) {
	// Create successor mocked object
	successRunner := new(fakeRunner)
	successRunner.On("Run", mock.Anything).Return(nil)

	// Create real object and call real method
	adapter := new(Adapter)
	adapter.runner = successRunner
	err := adapter.Run(pusu.NewSubscription("test", "testing", func(m *pusu.Message) error {
		return nil
	}))

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
	err := adapter.Run(pusu.NewSubscription("test", "testing", func(m *pusu.Message) error {
		return nil
	}))

	// Error must not be nil
	assert.Error(t, err)
}

func TestAdapter_CreateAdapter(t *testing.T) {
	// Call real method and check each interfaces have proper type
	adapter := CreateAdapter()

	_, httpHandlerProper := adapter.httpHandlerAdder.(*httpHandlerAdder)

	if adapter.httpHandlerAdder != nil && !httpHandlerProper {
		t.Errorf("Create Adapter Error: \nFollowing interfaces does not have proper type:"+
			"httpHandlerAdder: \n%t",
			httpHandlerProper)
	}
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
