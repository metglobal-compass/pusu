package main

import (
	"fmt"
	"github.com/metglobal-compass/pusu"
	"github.com/metglobal-compass/pusu/adapters/google"
	"os"
)

var (
	subscription *Subscription
	adapter      *google.Adapter
)

func init() {
	var err error
	subscription := new(Subscription)

	// Create adapter with environment variables
	adapter, err = google.CreateAdapter(os.Getenv("PUB_SUB_PROJECT_ID"), os.Getenv("BASE_HOST"))
	if err != nil {
		panic(err)
	}

	// Create subscription in google cloud
	err = adapter.CreateSubscription(subscription)
	if err != nil {
		panic(err)
	}
}

func main() {
	// Run subscription as http handler
	fmt.Println("Subscription is listening")
	err := adapter.Run(subscription)
	if err != nil {
		panic(err)
	}
}

// Subscription definition and method implementation of pusu.Subscription interface
type Subscription struct {
}

func (l *Subscription) Handle(m *pusu.Message) error {
	var err error

	fmt.Println("printing message...", m.Message().(string))

	return err
}

func (l *Subscription) Topic() string {
	return "printer"
}

func (l *Subscription) Name() string {
	return "printing"
}
