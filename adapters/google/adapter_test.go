package google

import (
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

type fakeCreatorSuccessor struct {
	called int
}

func (f *fakeCreatorSuccessor) CreateSubscription(subscription *pusu.Subscription) error {
	f.called = f.called + 1
	return nil
}
