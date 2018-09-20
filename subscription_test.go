package pusu

import (
	"errors"
	"testing"
)

func TestNewSubscription(t *testing.T) {
	// Create a sample subscription with a returning a failed message processing.
	subscription := NewSubscription("test", "testing", new(failureSubscriber))

	// Check topic is created as sent to function
	if subscription.Topic() != "test" {
		t.Errorf("Subscription Topic Error:\n Expected: %s Actual: %s ", "test", subscription.Topic())
	}

	// Check name is created as sent to function
	if subscription.Name() != "testing" {
		t.Errorf("Subscription Name Error:\n Expected: %s Actual: %s ", "testing", subscription.Name())
	}

	// Check subscriber is created as sent to function. It must return error with pre-defined message.
	err := subscription.Subscriber().Handle(new(Message))

	if err.Error() != "test error" {
		t.Errorf("Subscriber error:\nExpected error message:\n%s \nActual:\n%s", "test error", err.Error())
	}
}

type failureSubscriber struct {
}

func (f *failureSubscriber) Handle(m *Message) error {
	return errors.New("test error")
}
