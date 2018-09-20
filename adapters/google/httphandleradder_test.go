package google

import (
	"bytes"
	"fmt"
	"github.com/metglobal-compass/pusu"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHttpHandlerAdder_CreateSubscription(t *testing.T) {
	// Create subscription via httpdhandleradder
	handler := new(httpHandlerAdder)
	handler.CreateSubscription(pusu.NewSubscription("test", "testing", new(dummySubscriber)))

	// Create test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create request body and make request
	requestBody := []byte(`{"message": {"data": "W3sib2JqZWN0X2lkIjoxLCJvYmplY3RfbmFtZSI6IkFsbG90bWVudFBsYW4iLCJjaGlsZF9vYmplY3RfbmFtZSI6bnVsbCwib2JqZWN0X2RlZmluaXRpb24iOnsiaWQiOjEsImNsYXNzIjoiQWxsb3RtZW50UGxhbiJ9LCJhY3Rpb25fbmFtZSI6InVwZGF0ZSIsImxvZ190aW1lIjoiMjAxOC0wMi0xNiAxNTozNTowOSIsImNoYW5nZV9zZXQiOnsibmFtZSI6eyJvbGQiOiJCQVIiLCJuZXciOiJ0ZXN0cyJ9fSwiY29uc3VtZXJfbmFtZSI6IkNvbXBhc3MiLCJjb25zdW1lcl9pZCI6MSwiaXBfYWRkcmVzcyI6IjEwLjQuNC4xIiwidXNlcl9pZCI6MywidXNlcm5hbWUiOiJzZXlmaSIsImNsaWVudF9uYW1lIjoiSG90ZWxzcHJvIERNQ0MifV0="}}`)
	resp, _ := http.Post(
		fmt.Sprintf("%s%s", server.URL, handler.UrlPath(handler.subscription)),
		"application/json",
		bytes.NewReader(requestBody),
	)

	// Must return 200
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code is not valid. \nExcepted: 200\n Actual:%d", resp.StatusCode)
	}
}

func TestHttpHandlerAdder_UrlPath(t *testing.T) {
	// Create subscription via httpdhandleradder
	handler := new(httpHandlerAdder)
	subscription := pusu.NewSubscription("test", "testing", new(dummySubscriber))

	// Check generated path
	expectedPath := "/_handlers/topics/test/subscribers/testing"
	actualPath := handler.UrlPath(subscription)
	if actualPath != expectedPath {
		t.Errorf("Url path generation was broken. \nExpected:\n%s\nActual:\n%s", expectedPath, actualPath)
	}
}

func TestHttpHandlerAdder_ServeHTTP(t *testing.T) {
	// Create test request data
	body := []byte(`{"message": {"data": "W3sib2JqZWN0X2lkIjoxLCJvYmplY3RfbmFtZSI6IkFsbG90bWVudFBsYW4iLCJjaGlsZF9vYmplY3RfbmFtZSI6bnVsbCwib2JqZWN0X2RlZmluaXRpb24iOnsiaWQiOjEsImNsYXNzIjoiQWxsb3RtZW50UGxhbiJ9LCJhY3Rpb25fbmFtZSI6InVwZGF0ZSIsImxvZ190aW1lIjoiMjAxOC0wMi0xNiAxNTozNTowOSIsImNoYW5nZV9zZXQiOnsibmFtZSI6eyJvbGQiOiJCQVIiLCJuZXciOiJ0ZXN0cyJ9fSwiY29uc3VtZXJfbmFtZSI6IkNvbXBhc3MiLCJjb25zdW1lcl9pZCI6MSwiaXBfYWRkcmVzcyI6IjEwLjQuNC4xIiwidXNlcl9pZCI6MywidXNlcm5hbWUiOiJzZXlmaSIsImNsaWVudF9uYW1lIjoiSG90ZWxzcHJvIERNQ0MifV0="}}`)
	path := fmt.Sprintf("/_handlers/topics/%s/subscribers/%s", "test", "testing")
	req, _ := http.NewRequest("POST", path, bytes.NewReader(body))

	// Create fakeResponseWriter to hold response data
	w := httptest.NewRecorder()

	// Create http handler and call real method
	handler := new(httpHandlerAdder)
	handler.subscription = pusu.NewSubscription("test", "testing", new(dummySubscriber))
	handler.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Status code is not valid. \nExcepted: 200\n Actual:%d", w.Code)
	}
}

func TestHttpHandlerAdder_ServeHTTPErrorJson(t *testing.T) {
	// Create test request data
	body := []byte(`{JSONERROR}`)
	path := fmt.Sprintf("/_handlers/topics/%s/subscribers/%s", "test", "testing")
	req, _ := http.NewRequest("POST", path, bytes.NewReader(body))

	// Create fakeResponseWriter to hold response data
	w := httptest.NewRecorder()

	// Create http handler and call real method
	handler := new(httpHandlerAdder)
	handler.subscription = pusu.NewSubscription("test", "testing", new(dummySubscriber))
	handler.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status code is not valid. \nExcepted: 500\n Actual:%d", w.Code)
	}

	if strings.TrimSpace(w.Body.String()) != ErrorJsonSyntax {
		t.Errorf("\nExcepted Error Message: \n%s\nActual:\n%s", ErrorJsonSyntax, w.Body.String())
	}
}

func TestHttpHandlerAdder_ServeHTTPErrorBase64(t *testing.T) {
	// Create test request data
	body := []byte(`{"message": {"data": "WRONGMESSAGE="}}`)
	path := fmt.Sprintf("/_handlers/topics/%s/subscribers/%s", "test", "testing")
	req, _ := http.NewRequest("POST", path, bytes.NewReader(body))

	// Create fakeResponseWriter to hold response data
	w := httptest.NewRecorder()

	// Create http handler and call real method
	handler := new(httpHandlerAdder)
	handler.subscription = pusu.NewSubscription("test", "testing", new(dummySubscriber))
	handler.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status code is not valid. \nExcepted: 500\n Actual:%d", w.Code)
	}

	if strings.TrimSpace(w.Body.String()) != ErrorBase64MessageSyntax {
		t.Errorf("\nExcepted Error Message: \n%s\nActual:\n%s", ErrorBase64MessageSyntax, w.Body.String())
	}
}

func TestHttpHandlerAdder_ServeHTTPSubscriberError(t *testing.T) {
	// Create test request data
	body := []byte(`{"message": {"data": "W3sib2JqZWN0X2lkIjoxLCJvYmplY3RfbmFtZSI6IkFsbG90bWVudFBsYW4iLCJjaGlsZF9vYmplY3RfbmFtZSI6bnVsbCwib2JqZWN0X2RlZmluaXRpb24iOnsiaWQiOjEsImNsYXNzIjoiQWxsb3RtZW50UGxhbiJ9LCJhY3Rpb25fbmFtZSI6InVwZGF0ZSIsImxvZ190aW1lIjoiMjAxOC0wMi0xNiAxNTozNTowOSIsImNoYW5nZV9zZXQiOnsibmFtZSI6eyJvbGQiOiJCQVIiLCJuZXciOiJ0ZXN0cyJ9fSwiY29uc3VtZXJfbmFtZSI6IkNvbXBhc3MiLCJjb25zdW1lcl9pZCI6MSwiaXBfYWRkcmVzcyI6IjEwLjQuNC4xIiwidXNlcl9pZCI6MywidXNlcm5hbWUiOiJzZXlmaSIsImNsaWVudF9uYW1lIjoiSG90ZWxzcHJvIERNQ0MifV0="}}`)
	path := fmt.Sprintf("/_handlers/topics/%s/subscribers/%s", "test", "testing")
	req, _ := http.NewRequest("POST", path, bytes.NewReader(body))

	// Create fakeResponseWriter to hold response data
	w := httptest.NewRecorder()

	// Create http handler and call real method. Subscriber must return error
	handler := new(httpHandlerAdder)
	handler.subscription = pusu.NewSubscription("test", "testing", new(failureSubscriber))
	handler.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status code is not valid. \nExcepted: 500\n Actual:%d", w.Code)
	}

	if strings.TrimSpace(w.Body.String()) != ErrorMessageExecution {
		t.Errorf("\nExcepted Error Message: \n%s\nActual:\n%s", ErrorMessageExecution, w.Body.String())
	}
}

func TestHttpHandlerAdder_ServeHTTPSubscriberPathError(t *testing.T) {
	// Create test request data with wrong url path
	body := []byte(`{"message": {"data": "W3sib2JqZWN0X2lkIjoxLCJvYmplY3RfbmFtZSI6IkFsbG90bWVudFBsYW4iLCJjaGlsZF9vYmplY3RfbmFtZSI6bnVsbCwib2JqZWN0X2RlZmluaXRpb24iOnsiaWQiOjEsImNsYXNzIjoiQWxsb3RtZW50UGxhbiJ9LCJhY3Rpb25fbmFtZSI6InVwZGF0ZSIsImxvZ190aW1lIjoiMjAxOC0wMi0xNiAxNTozNTowOSIsImNoYW5nZV9zZXQiOnsibmFtZSI6eyJvbGQiOiJCQVIiLCJuZXciOiJ0ZXN0cyJ9fSwiY29uc3VtZXJfbmFtZSI6IkNvbXBhc3MiLCJjb25zdW1lcl9pZCI6MSwiaXBfYWRkcmVzcyI6IjEwLjQuNC4xIiwidXNlcl9pZCI6MywidXNlcm5hbWUiOiJzZXlmaSIsImNsaWVudF9uYW1lIjoiSG90ZWxzcHJvIERNQ0MifV0="}}`)
	path := fmt.Sprintf("/_handlers/topics/%s/subscribers/%s", "wrong_topic", "wrong_subscription_name")
	req, _ := http.NewRequest("POST", path, bytes.NewReader(body))

	// Create fakeResponseWriter to hold response data
	w := httptest.NewRecorder()

	// Create http handler and call real method
	handler := new(httpHandlerAdder)
	handler.subscription = pusu.NewSubscription("test", "testing", new(dummySubscriber))
	handler.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusNotFound {
		t.Errorf("Status code is not valid. \nExcepted: 200\n Actual:%d", w.Code)
	}
}

func TestHttpHandlerAdder_ServeHTTPSubscriberMethodTypeError(t *testing.T) {
	// Create test request data with wrong method type
	body := []byte(`{"message": {"data": "W3sib2JqZWN0X2lkIjoxLCJvYmplY3RfbmFtZSI6IkFsbG90bWVudFBsYW4iLCJjaGlsZF9vYmplY3RfbmFtZSI6bnVsbCwib2JqZWN0X2RlZmluaXRpb24iOnsiaWQiOjEsImNsYXNzIjoiQWxsb3RtZW50UGxhbiJ9LCJhY3Rpb25fbmFtZSI6InVwZGF0ZSIsImxvZ190aW1lIjoiMjAxOC0wMi0xNiAxNTozNTowOSIsImNoYW5nZV9zZXQiOnsibmFtZSI6eyJvbGQiOiJCQVIiLCJuZXciOiJ0ZXN0cyJ9fSwiY29uc3VtZXJfbmFtZSI6IkNvbXBhc3MiLCJjb25zdW1lcl9pZCI6MSwiaXBfYWRkcmVzcyI6IjEwLjQuNC4xIiwidXNlcl9pZCI6MywidXNlcm5hbWUiOiJzZXlmaSIsImNsaWVudF9uYW1lIjoiSG90ZWxzcHJvIERNQ0MifV0="}}`)
	path := fmt.Sprintf("/_handlers/topics/%s/subscribers/%s", "test", "testing")
	req, _ := http.NewRequest("GET", path, bytes.NewReader(body))

	// Create fakeResponseWriter to hold response data
	w := httptest.NewRecorder()

	// Create http handler and call real method
	handler := new(httpHandlerAdder)
	handler.subscription = pusu.NewSubscription("test", "testing", new(dummySubscriber))
	handler.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusNotFound {
		t.Errorf("Status code is not valid. \nExcepted: 200\n Actual:%d", w.Code)
	}
}
