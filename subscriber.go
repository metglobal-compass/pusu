package pusu

// Subscriber type is base definition of any subscription's business logic.
// Function gets generic pusu.Message type and returns nil as error after successful handling of message.
// If function return a non-nil error type, adapter try again for later attempt until gets a successful response.
type Subscriber func(m *Message) error
