package qst

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	t.Run("BodyJSON error", func(t *testing.T) {
		_, err := NewPost("https://example.com", BodyJSON(make(chan struct{})))
		assert.EqualError(t, err, "json: unsupported type: chan struct {}")
	})

	t.Run("BodyXML error", func(t *testing.T) {
		_, err := NewPost("https://example.com", BodyXML(make(chan struct{})))
		assert.EqualError(t, err, "xml: unsupported type: chan struct {}")
	})

	t.Run("Context", func(t *testing.T) {
		type keyType struct{}
		var key keyType

		request, err := NewPost("https://example.com",
			Context(context.WithValue(context.TODO(), key, "here")),
		)

		assert.NoError(t, err)

		assert.Equal(t, "here", request.Context().Value(key))
	})

	t.Run("ContextValue", func(t *testing.T) {
		type keyType struct{}
		var key keyType

		request, err := NewPost("https://example.com", ContextValue(key, "here"))

		assert.NoError(t, err)

		assert.Equal(t, "here", request.Context().Value(key))
	})
}

func ExampleQuery() {
	req, _ := NewGet("https://example.com", Query("page", "10"))
	fmt.Println(req.URL.Query().Encode())
	// Output: page=10
}

func ExamplePath() {
	req, _ := NewGet("https://example.com/api", Path("to", "some", "resource"))
	fmt.Println(req.URL.Path)
	// Output: /api/to/some/resource
}

func ExampleQueries() {
	req, _ := NewGet("https://example.com", Queries{"page": {"10"}, "count": {"50"}})
	fmt.Println(req.URL.Query().Encode())
	// Output: count=50&page=10
}

func ExampleHeader() {
	req, _ := NewGet("https://example.com", Header("X-Trace-Id", "asdf"))
	fmt.Println(req.Header.Get("x-trace-id"))
	// Output: asdf
}

func ExampleHeaders() {
	req, _ := NewGet("https://example.com", Headers{"X-Trace-Id": {"asdf"}})
	fmt.Println(req.Header.Get("x-trace-id"))
	// Output: asdf
}

func ExampleCookie() {
	req, _ := NewGet("https://example.com",
		Cookie(&http.Cookie{
			Name:    "some-cookie",
			Value:   "some-value",
			Path:    "/",
			Expires: time.Now().Add(time.Hour),
		}),
	)

	fmt.Println(req.Cookie("some-cookie"))
	// Output: some-cookie=some-value <nil>
}

func ExampleAuthorization() {
	req, _ := NewGet("https://example.com", Authorization("some-token"))
	fmt.Println(req.Header.Get("Authorization"))
	// Output: some-token
}

func ExampleBasicAuth() {
	req, _ := NewGet("https://example.com", BasicAuth("someone", "hunter12"))
	fmt.Println(req.Header.Get("Authorization"))
	// Output: Basic c29tZW9uZTpodW50ZXIxMg==
}

func ExampleBearerAuth() {
	req, _ := NewGet("https://example.com", BearerAuth("some-token"))
	fmt.Println(req.Header.Get("Authorization"))
	// Output: Bearer some-token
}

func ExampleContextValue() {
	req, _ := NewGet("https://example.com", ContextValue("key", "value"))
	fmt.Println(req.Context().Value("key"))
	// Output: value
}

func ExampleBody() {
	req, _ := NewPost("https://example.com", Body(ioutil.NopCloser(bytes.NewBufferString("something"))))
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	// Output: something
}

func ExampleBodyBytes() {
	req, _ := NewPost("https://example.com", BodyBytes([]byte("something")))
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	// Output: something
}

func ExampleBodyString() {
	req, _ := NewPost("https://example.com", BodyString("something"))
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	// Output: something
}

func ExampleBodyForm() {
	req, _ := NewPost("https://example.com", BodyForm{"something": {"here"}})
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	// Output: something=here
}

func ExampleBodyJSON() {
	req, _ := NewPost("https://example.com", BodyJSON(map[string]interface{}{"something": "here"}))
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	// Output: {"something":"here"}
}

func ExampleBodyXML() {
	req, _ := NewPost("https://example.com", BodyXML("something"))
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	// Output: <string>something</string>
}
