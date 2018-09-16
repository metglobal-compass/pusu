package pusu

// Holds immutable data of subscription.
type Subscription struct {
	name       string
	topic      string
	subscriber Subscriber
}

// Returns name of subscription
func (s *Subscription) Name() string {
	return s.name
}

// Returns topic of subscription
func (s *Subscription) Topic() string {
	return s.topic
}

// Returns subscription's base logic
func (s *Subscription) Subscriber() Subscriber {
	return s.subscriber
}

// Creates new subscription
func NewSubscription(topic string, name string, subscriber Subscriber) *Subscription {
	s := new(Subscription)
	s.topic = topic
	s.name = name
	s.subscriber = subscriber

	return s
}
