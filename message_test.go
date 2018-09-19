package pusu

import "testing"

func TestMessage_NewMessage(t *testing.T) {
	expectedPayload := "testmessagedata"
	message := NewMessage(expectedPayload)

	actualPayload := message.Message()

	if actualPayload != expectedPayload {
		t.Errorf("Error: Expected: %s, Actual: %s", expectedPayload, actualPayload)
	}
}
