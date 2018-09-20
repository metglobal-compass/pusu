package pusu

// Any adapter implementation of cloud pub/sub workflow must implement Creator and Runner interface
type Adapter interface {
	Creator
	Runner
}

// Creator implementation must check or create any relevant data about given subscription.
// Those may include creating topics, subscriber information or default configuration settings.
type Creator interface {
	CreateSubscription(subscription Subscription) error
}

// Runner implementation must trigger proper workflow for cloud vendor's must-have application up and running logic.
type Runner interface {
	Run(subscription Subscription) error
}
