package pusu

// Message struct holds immutable message payload.
type Message struct {
	message interface{}
}

// Creates new message with payload message data
func NewMessage(message interface{}) *Message {
	m := new(Message)
	m.message = message

	return m
}

// Get payload message
func (m *Message) GetMessage() interface{} {
	return m.message
}
