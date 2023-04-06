package qst

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		_, err := New("lol what", "https://example.com")
		assert.EqualError(t, err, `net/http: invalid method "lol what"`)
	})
}

func ExampleNew() {
	// Output:
	// POST
	// https://httpbin.org/post?limit=10
	// Bearer some-token
	// {"key":"value"}

	req, _ := New(http.MethodPost, "https://example.com",
		Scheme("https"),
		Host("httpbin.org"),
		Path("/post"),
		BearerAuth("some-token"),
		Query("limit", "10"),
		BodyJSON(map[string]string{"key": "value"}),
	)

	body, _ := ioutil.ReadAll(req.Body)

	fmt.Println(req.Method)
	fmt.Println(req.URL)
	fmt.Println(req.Header.Get("Authorization"))
	fmt.Println(string(body))
}

func ExampleDo() {
	// Output:
	// POST
	// /?limit=10
	// Bearer some-token
	// {"key":"value"}

	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)

		fmt.Println(req.Method)
		fmt.Println(req.URL)
		fmt.Println(req.Header.Get("Authorization"))
		fmt.Println(string(body))
	}))
	defer server.Close()

	Do(http.MethodPost, "https://example.com",
		URL(server.URL),
		BearerAuth("some-token"),
		Query("limit", "10"),
		BodyJSON(map[string]string{"key": "value"}),
	)
}
