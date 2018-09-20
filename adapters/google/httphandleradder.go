package google

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/metglobal-compass/pusu"
	"net/http"
)

type httpHandlerAdder struct {
	subscription pusu.Subscription
}

// Implementation of internal Creator interface for Google Adapter
func (h *httpHandlerAdder) CreateSubscription(subscription pusu.Subscription) error {
	h.subscription = subscription
	http.Handle(h.UrlPath(subscription), h)
	return nil
}

// Implementation of handler interface of net/http Handler interface
func (h *httpHandlerAdder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check path in case of anything
	if r.URL.Path != h.UrlPath(h.subscription) || r.Method != http.MethodPost {
		http.Error(w, ErrorMessageExecution, http.StatusNotFound)
		return
	}

	// Convert pubsubmessage structure to pusu.Message
	var m message
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil || m.Message.Data == "" {
		http.Error(w, ErrorJsonSyntax, http.StatusInternalServerError)
		return
	}

	// Decode base64 encoded pub/sub data to string
	s, err := base64.StdEncoding.DecodeString(m.Message.Data)
	if err != nil {
		http.Error(w, ErrorBase64MessageSyntax, http.StatusInternalServerError)
		return
	}

	// Create message
	pusuMessage := pusu.NewMessage(string(s))

	// Execute real method of subscription
	err = h.subscription.Handle(pusuMessage)

	// Return 500 status code in case of any error, otherwise do nothing
	if err != nil {
		http.Error(w, ErrorMessageExecution, http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// Get url path of subscriber
func (h *httpHandlerAdder) UrlPath(subscription pusu.Subscription) string {
	return fmt.Sprintf("/_handlers/topics/%s/subscribers/%s", subscription.Topic(), subscription.Name())
}
