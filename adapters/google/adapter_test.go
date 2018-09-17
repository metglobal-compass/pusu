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
