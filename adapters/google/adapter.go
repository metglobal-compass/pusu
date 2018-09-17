// Package google provides recommended implementation of pub/sub workflow for Golang.
// Google Cloud Platform Subscribers are based on Google App Engine Flexible Environment for Go.
// Each subscriber are separate App Engine service and scalable as needed.
// Google Cloud Pub/Sub pushes triggered messages to those App Engine services.
// If message processing is successful, pusu returns a 200 OK response code and Pub/Sub acknowledges the message
// If it is unsuccessful, pusu returns 400 or 500 response code and Pub/Sub tries later until gets a success message.
package google

import (
	"errors"
	"github.com/metglobal-compass/pusu"
)

type Adapter struct {
	// Adds topic information to Google Cloud
	topicAdder pusu.Creator

	// Adds subscriber information to Google Cloud
	subscriberAdder pusu.Creator

	// Adds relevant routing information to http package and handles pre-processing http message to pusu.Message
	httpHandlerAdder pusu.Creator
}

// Implementation of pusu.Creator interface as part of pusu.Adapter interface
func (g *Adapter) CreateSubscription(subscription *pusu.Subscription) error {
	// Validate subscription
	if subscription.Name() == "" {
		return errors.New("Subscription name must not be empty. ")
	}
	if subscription.Topic() == "" {
		return errors.New("Subscription topic must not be empty. ")
	}
	if subscription.Subscriber() == nil {
		return errors.New("Subscription handler must not be empty. ")
	}

	err := g.httpHandlerAdder.CreateSubscription(subscription)
	if err != nil {
		return err
	}

	err = g.topicAdder.CreateSubscription(subscription)
	if err != nil {
		return err
	}

	err = g.subscriberAdder.CreateSubscription(subscription)
	if err != nil {
		return err
	}

	return nil
}
