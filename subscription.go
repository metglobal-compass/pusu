package pusu

type Subscription interface {
	// Returns the name of topic
	Topic() string

	// Returns the name of subscription
	Name() string

	// Handles the pub/sub message
	// Function gets generic pusu.Message type and returns nil as error after successful handling of message.
	// If function return a non-nil error type, adapter try again for later attempt until gets a successful response.
	Handle(m *Message) error
}
