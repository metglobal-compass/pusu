package google

// Google Cloud Pub/Sub message structure
type message struct {
	Message struct {
		Data string `json:"data"`
	} `json:"message"`
}
