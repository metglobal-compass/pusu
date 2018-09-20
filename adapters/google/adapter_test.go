package google

import (
	"errors"
	"github.com/metglobal-compass/pusu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAdapter_CreateSubscriptionErrorOnEmptyTopic(t *testing.T) {
	// Create adapter
	adapter := new(Adapter)

	// Create subscription which has empty topic name. Must return error when topic name is empty
	err := adapter.CreateSubscription(new(fakeSubscription).WithName("testing").WithTopic(""))
	assert.Error(t, err)
}
func TestAdapter_CreateSubscriptionErrorOnEmptyName(t *testing.T) {
	// Create adapter
	adapter := new(Adapter)

	// Create subscription which has empty subscription name. Must return error when name is empty
	err := adapter.CreateSubscription(new(fakeSubscription).WithName(""))
	assert.Error(t, err)
}

func TestAdapter_CreateSubscription(t *testing.T) {
	// Create mocked object
	successCreator := new(fakeCreator)
	successCreator.On("CreateSubscription", mock.Anything).Return(nil)

	// Create real instance and call real method
	adapter := new(Adapter)
	adapter.cloudAdder = successCreator
	adapter.httpHandlerAdder = successCreator
	err := adapter.CreateSubscription(new(fakeSubscription).WillHaveProperFields())

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
	err := adapter.CreateSubscription(new(fakeSubscription).WillHaveProperFields())
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
	err := adapter.CreateSubscription(new(fakeSubscription).WillHaveProperFields())
	assert.Error(t, err)
}

func TestAdapter_Run(t *testing.T) {
	// Create successor mocked object
	successRunner := new(fakeRunner)
	successRunner.On("Run", mock.Anything).Return(nil)

	// Create real object and call real method
	adapter := new(Adapter)
	adapter.runner = successRunner
	err := adapter.Run(new(fakeSubscription).WillHaveProperFields())

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
	err := adapter.Run(new(fakeSubscription).WillHaveProperFields())

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

// A fake creator definition
type fakeCreator struct {
	mock.Mock
}

func (f *fakeCreator) CreateSubscription(subscription pusu.Subscription) error {
	args := f.Called(subscription)
	return args.Error(0)
}

// A fake runner definition
type fakeRunner struct {
	mock.Mock
}

func (f *fakeRunner) Run(subscription pusu.Subscription) error {
	args := f.Called(subscription)
	return args.Error(0)
}

// A fake subscription definition
type fakeSubscription struct {
	mock.Mock
}

func (f *fakeSubscription) Handle(m *pusu.Message) error {
	args := f.Called(m)
	return args.Error(0)
}
func (f *fakeSubscription) Topic() string {
	args := f.Called()
	return args.String(0)
}

func (f *fakeSubscription) Name() string {
	args := f.Called()
	return args.String(0)
}

func (f *fakeSubscription) WithTopic(topic string) *fakeSubscription {
	f.On("Topic").Return(topic)
	return f
}

func (f *fakeSubscription) WithName(name string) *fakeSubscription {
	f.On("Name").Return(name)
	return f
}

func (f *fakeSubscription) WithReturning(err error) *fakeSubscription {
	f.On("Handle", mock.Anything).Return(err)
	return f
}

func (f *fakeSubscription) WillHaveProperFields() pusu.Subscription {
	return f.WithTopic("test").WithName("testing").WithReturning(nil)
}

func (f *fakeSubscription) WillReturnError() *fakeSubscription {
	return f.WithTopic("test").WithName("testing").WithReturning(errors.New("error"))
}
