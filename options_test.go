package qst

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	t.Run("BodyJSON error", func(t *testing.T) {
		_, err := NewPost("", BodyJSON(make(chan struct{})))
		assert.EqualError(t, err, "json: unsupported type: chan struct {}")
	})

	t.Run("BodyXML error", func(t *testing.T) {
		_, err := NewPost("", BodyXML(make(chan struct{})))
		assert.EqualError(t, err, "xml: unsupported type: chan struct {}")
	})

	t.Run("Context", func(t *testing.T) {
		type keyType struct{}
		var key keyType

		request, err := NewPost("", Context(context.WithValue(context.TODO(), key, "here")))
		assert.NoError(t, err)

		assert.Equal(t, "here", request.Context().Value(key))
	})
}

func ExampleQueryValue() {
	req, _ := NewGet("http://httpbin.org/get", QueryValue("page", "10"))
	fmt.Println(req.URL.Query().Encode())
	// Output: page=10
}

func ExampleQuery() {
	req, _ := NewGet("http://httpbin.org/get", Query{"page": {"10"}, "count": {"50"}})
	fmt.Println(req.URL.Query().Encode())
	// Output: count=50&page=10
}

func ExampleHeader() {
	req, _ := NewGet("http://httpbin.org/get", Header("X-Trace-Id", "asdf"))
	fmt.Println(req.Header.Get("x-trace-id"))
	// Output: asdf
}

func ExampleHeaders() {
	req, _ := NewGet("http://httpbin.org/get", Headers{"X-Trace-Id": {"asdf"}})
	fmt.Println(req.Header.Get("x-trace-id"))
	// Output: asdf
}

func ExampleAuthorization() {
	req, _ := NewGet("http://httpbin.org/get", Authorization("some-token"))
	fmt.Println(req.Header.Get("Authorization"))
	// Output: some-token
}

func ExampleBasicAuth() {
	req, _ := NewGet("http://httpbin.org/get", BasicAuth("someone", "hunter12"))
	fmt.Println(req.Header.Get("Authorization"))
	// Output: Basic c29tZW9uZTpodW50ZXIxMg==
}

func ExampleBearer() {
	req, _ := NewGet("http://httpbin.org/get", BearerAuth("some-token"))
	fmt.Println(req.Header.Get("Authorization"))
	// Output: Bearer some-token
}

func ExampleBody() {
	req, _ := NewPost("http://httpbin.org/post", Body(ioutil.NopCloser(bytes.NewBufferString("something"))))
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	// Output: something
}

func ExampleBodyBytes() {
	req, _ := NewPost("http://httpbin.org/post", BodyBytes([]byte("something")))
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	// Output: something
}

func ExampleBodyString() {
	req, _ := NewPost("http://httpbin.org/post", BodyString("something"))
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	// Output: something
}

func ExampleBodyForm() {
	req, _ := NewPost("http://httpbin.org/post", BodyForm{"something": {"here"}})
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	// Output: something=here
}

func ExampleBodyJSON() {
	req, _ := NewPost("http://httpbin.org/post", BodyJSON(map[string]interface{}{"something": "here"}))
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	// Output: {"something":"here"}
}

func ExampleBodyXML() {
	req, _ := NewPost("http://httpbin.org/post", BodyXML("something"))
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	// Output: <string>something</string>
}
