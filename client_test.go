package qst

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Do(t *testing.T) {
	_, err := NewClient(http.DefaultClient).Do("lol what", "")
	assert.EqualError(t, err, `net/http: invalid method "lol what"`)
}

func ExampleNewClient() {
	client := NewClient(http.DefaultClient, BearerAuth("asdf"))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header.Get("Authorization"))
		fmt.Println(r.URL.RawQuery)
	}))

	client.Get(server.URL, QueryValue("page", "10"))
	// Output: Bearer asdf
	// page=10
}
