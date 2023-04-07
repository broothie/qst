package qst_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/broothie/qst"
	"github.com/stretchr/testify/assert"
)

func TestClient_Do(t *testing.T) {
	_, err := qst.NewClient(http.DefaultClient).Do("lol what")
	assert.EqualError(t, err, `net/http: invalid method "lol what"`)
}

func ExampleNewClient() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header.Get("Authorization"))
		fmt.Println(r.URL.RawQuery)
	}))
	defer server.Close()

	client := qst.NewClient(http.DefaultClient,
		qst.URL(server.URL),
		qst.BearerAuth("c0rnfl@k3s"),
	)

	client.Get(qst.Query("page", "10"))
	// Output: Bearer c0rnfl@k3s
	// page=10
}
