// Package google provides recommended implementation of pub/sub workflow for Golang.
// Google Cloud Platform Subscribers are based on Google App Engine Flexible Environment for Go.
// Each subscriber are separate App Engine service and scalable as needed.
// Google Cloud Pub/Sub pushes triggered messages to those App Engine services.
// If message processing is successful, pusu returns a 200 OK response code and Pub/Sub acknowledges the message
// If it is unsuccessful, pusu returns 400 or 500 response code and Pub/Sub tries later until gets a success message.
package google

import (
	"cloud.google.com/go/pubsub"
	"context"
	"errors"
	"github.com/metglobal-compass/pusu"
)

type Adapter struct {
	// Adds topic and subscription information to Google Cloud Pub/Sub
	cloudAdder pusu.Creator

	// Adds relevant routing information to http package and handles pre-processing http message to pusu.Message
	httpHandlerAdder pusu.Creator

	// Runs subscription as HTTP App Engine service
	runner pusu.Runner
}

// Implementation of pusu.Creator interface as part of pusu.Adapter interface
func (g *Adapter) CreateSubscription(subscription pusu.Subscription) error {
	// Validate subscription
	if subscription.Name() == "" {
		return errors.New("Subscription name must not be empty. ")
	}
	if subscription.Topic() == "" {
		return errors.New("Subscription topic must not be empty. ")
	}

	err := g.cloudAdder.CreateSubscription(subscription)
	if err != nil {
		return err
	}

	err = g.httpHandlerAdder.CreateSubscription(subscription)
	if err != nil {
		return err
	}

	return nil
}

// Implementation of pusu.Runner interface as part of pusu.Adapter interface
func (g *Adapter) Run(subscription pusu.Subscription) error {
	err := g.runner.Run(subscription)
	return err
}

// Creates Google Adapter
// projectId: Google Cloud Project Id
// host: Base host uri of app engine based subscriber http handlers. (Ex: https://servicename.appspot.com/)
func CreateAdapter(projectId string, host string) (*Adapter, error) {
	// Validate parameters
	if projectId == "" {
		return nil, errors.New("projectId must not be empty")
	}

	if host == "" {
		return nil, errors.New("host for subscriber http handlers must not be empty")
	}

	googleAdapter := new(Adapter)
	googleAdapter.httpHandlerAdder = new(httpHandlerAdder)

	// Add pub/sub client
	client, err := pubsub.NewClient(context.Background(), projectId)
	if err != nil {
		return nil, err
	}
	googleAdapter.cloudAdder = &cloudAdder{client: &pubSubClientWrapper{client: client}, host: host}

	// Add appengine runner
	googleAdapter.runner = new(appEngineRunner)

	return googleAdapter, nil
}
