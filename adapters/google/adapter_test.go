package google

import (
	"errors"
	"github.com/metglobal-compass/pusu"
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
	adapter := new(Adapter)
	adapter.topicAdder = new(fakeCreatorSuccessor)
	adapter.subscriberAdder = new(fakeCreatorSuccessor)
	adapter.httpHandlerAdder = new(fakeCreatorSuccessor)

	subscription := pusu.NewSubscription("test", "testing", func(m *pusu.Message) error {
		return nil
	})

	// There must be no error
	err := adapter.CreateSubscription(subscription)
	if err != nil {
		t.Errorf("Error while creating subscription with following message: \n %s", err.Error())
	}

	// Topic adder must be called once
	topicAdder := adapter.topicAdder.(*fakeCreatorSuccessor)
	if topicAdder.called != 1 {
		t.Error("Error while creating subscription. Topic adder must be called exactly once")
	}

	// Subscriber adder must be called once
	subscriberAdder := adapter.subscriberAdder.(*fakeCreatorSuccessor)
	if subscriberAdder.called != 1 {
		t.Error("Error while creating subscription. Subscriber adder must be called exactly once")
	}

	// Http handler must be called once
	httpHandlerAdder := adapter.httpHandlerAdder.(*fakeCreatorSuccessor)
	if httpHandlerAdder.called != 1 {
		t.Error("Error while creating subscription. Http handler adder must be called exactly once")
	}
}

func TestAdapter_CreateSubscriptionOnTopicFailure(t *testing.T) {
	adapter := new(Adapter)
	adapter.topicAdder = new(fakeCreatorFailure)
	adapter.subscriberAdder = new(fakeCreatorSuccessor)
	adapter.httpHandlerAdder = new(fakeCreatorSuccessor)

	subscription := pusu.NewSubscription("test", "testing", func(m *pusu.Message) error {
		return nil
	})

	// Subscription creation must return error when topic adder return error
	err := adapter.CreateSubscription(subscription)
	topicAdder := adapter.topicAdder.(*fakeCreatorFailure)
	if err == nil || topicAdder.called != 1 {
		t.Error("subscription creation must return error when topic adder returns error")
	}

	// Subscriber adder must not be called after failure on topic creation
	subscriberAdder := adapter.subscriberAdder.(*fakeCreatorSuccessor)
	if subscriberAdder.called > 0 {
		t.Error("Subscriber adder must not be called after failure on topic creation")
	}
}

func TestAdapter_CreateSubscriptionOnSubscriberFailure(t *testing.T) {
	adapter := new(Adapter)
	adapter.topicAdder = new(fakeCreatorSuccessor)
	adapter.subscriberAdder = new(fakeCreatorFailure)
	adapter.httpHandlerAdder = new(fakeCreatorSuccessor)

	subscription := pusu.NewSubscription("test", "testing", func(m *pusu.Message) error {
		return nil
	})

	// Subscription creation must return error when subscriber adder return error
	err := adapter.CreateSubscription(subscription)
	subscriberAdder := adapter.subscriberAdder.(*fakeCreatorFailure)
	if err == nil || subscriberAdder.called != 1 {
		t.Error("subscription creation must return error when subscriber adder returns error")
	}
}

func TestAdapter_CreateSubscriptionOnHttpHandlerFailure(t *testing.T) {
	adapter := new(Adapter)
	adapter.topicAdder = new(fakeCreatorSuccessor)
	adapter.subscriberAdder = new(fakeCreatorSuccessor)
	adapter.httpHandlerAdder = new(fakeCreatorFailure)

	subscription := pusu.NewSubscription("test", "testing", func(m *pusu.Message) error {
		return nil
	})

	// Subscription creation must return error when subscriber adder return error
	err := adapter.CreateSubscription(subscription)
	httpHandlerAdder := adapter.httpHandlerAdder.(*fakeCreatorFailure)
	if err == nil || httpHandlerAdder.called != 1 {
		t.Error("subscription creation must return error when http handler adder returns error")
	}
}

func TestAdapter_Run(t *testing.T) {
	// Adapter must return nil error when runner successful
	adapter := new(Adapter)
	adapter.runner = new(fakeRunnerSuccessor)

	subscription := pusu.NewSubscription("test", "testing", func(m *pusu.Message) error {
		return nil
	})

	err := adapter.Run(subscription)
	if err != nil {
		t.Errorf("Error while running subscription with following message: \n %s", err.Error())
	}

	// Adapter must return error when runner unsuccessful
	adapter.runner = new(fakeRunnerFailure)
	err = adapter.Run(subscription)
	if err == nil {
		t.Error("Adapter must return error when runner unsuccessful")
	}
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

type fakeCreatorSuccessor struct {
	called int
}

func (f *fakeCreatorSuccessor) CreateSubscription(subscription *pusu.Subscription) error {
	f.called = f.called + 1
	return nil
}

type fakeCreatorFailure struct {
	called int
}

func (f *fakeCreatorFailure) CreateSubscription(subscription *pusu.Subscription) error {
	f.called = f.called + 1
	return errors.New("create subscription error")
}

type fakeRunnerSuccessor struct {
	called int
}

func (f *fakeRunnerSuccessor) Run(subscription *pusu.Subscription) error {
	f.called = f.called + 1
	return nil
}

type fakeRunnerFailure struct {
	called int
}

func (f *fakeRunnerFailure) Run(subscription *pusu.Subscription) error {
	f.called = f.called + 1
	return errors.New("running subscription error")
}
