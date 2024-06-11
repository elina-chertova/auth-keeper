package sender

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient() *Client {
	return NewClient("")
}

func TestUnsupportedMethod(t *testing.T) {
	client := newTestClient()
	_, err := client.SendRequest("PUT", "/test", nil, "")
	assert.NotNil(t, err)
	assert.Equal(t, "error sending request: unsupported HTTP method: PUT", err.Error())
}

func TestSendGetRequest(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Authorization", "Bearer test-token")
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, `{"message":"success"}`)
			},
		),
	)
	defer server.Close()

	client := NewClient(server.URL)
	resp, err := client.SendRequest("GET", "/", nil, "")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "test-token", client.AuthToken)
}

func TestSendPostRequest(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Authorization", "Bearer test-token")
				w.WriteHeader(http.StatusCreated)
				fmt.Fprintln(w, `{"message":"created"}`)
			},
		),
	)
	defer server.Close()

	client := NewClient(server.URL)
	resp, err := client.SendRequest("POST", "/", map[string]string{"key": "value"}, "")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "test-token", client.AuthToken)
}
