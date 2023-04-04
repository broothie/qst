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
		_, err := NewPost(BodyJSON(make(chan struct{})))
		assert.EqualError(t, err, "json: unsupported type: chan struct {}")
	})

	t.Run("BodyXML error", func(t *testing.T) {
		_, err := NewPost(BodyXML(make(chan struct{})))
		assert.EqualError(t, err, "xml: unsupported type: chan struct {}")
	})

	t.Run("Context", func(t *testing.T) {
		type keyType struct{}
		var key keyType

		request, err := NewPost(Context(context.WithValue(context.TODO(), key, "here")))
		assert.NoError(t, err)

		assert.Equal(t, "here", request.Context().Value(key))
	})
}

func ExampleQuery() {
	req, _ := NewGet(Query("page", "10"))
	fmt.Println(req.URL.Query().Encode())
	// Output: page=10
}

func ExamplePath() {
	req, _ := NewGet(URL("https://example.com/api"), Path("to", "some", "resource"))
	fmt.Println(req.URL.Path)
	// Output: /api/to/some/resource
}

func ExampleQueries() {
	req, _ := NewGet(Queries{"page": {"10"}, "count": {"50"}})
	fmt.Println(req.URL.Query().Encode())
	// Output: count=50&page=10
}

func ExampleHeader() {
	req, _ := NewGet(Header("X-Trace-Id", "asdf"))
	fmt.Println(req.Header.Get("x-trace-id"))
	// Output: asdf
}

func ExampleHeaders() {
	req, _ := NewGet(Headers{"X-Trace-Id": {"asdf"}})
	fmt.Println(req.Header.Get("x-trace-id"))
	// Output: asdf
}

func ExampleCookie() {
	req, _ := NewGet(Cookie(&http.Cookie{
		Name:    "some-cookie",
		Value:   "some-value",
		Path:    "/",
		Expires: time.Now().Add(time.Hour),
	}))

	fmt.Println(req.Cookie("some-cookie"))
	// Output: some-cookie=some-value <nil>
}

func ExampleAuthorization() {
	req, _ := NewGet(Authorization("some-token"))
	fmt.Println(req.Header.Get("Authorization"))
	// Output: some-token
}

func ExampleBasicAuth() {
	req, _ := NewGet(BasicAuth("someone", "hunter12"))
	fmt.Println(req.Header.Get("Authorization"))
	// Output: Basic c29tZW9uZTpodW50ZXIxMg==
}

func ExampleBearerAuth() {
	req, _ := NewGet(BearerAuth("some-token"))
	fmt.Println(req.Header.Get("Authorization"))
	// Output: Bearer some-token
}

func ExampleContextValue() {
	req, _ := NewGet(ContextValue("key", "value"))
	fmt.Println(req.Context().Value("key"))
	// Output: value
}

func ExampleBody() {
	req, _ := NewPost(Body(ioutil.NopCloser(bytes.NewBufferString("something"))))
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	// Output: something
}

func ExampleBodyBytes() {
	req, _ := NewPost(BodyBytes([]byte("something")))
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	// Output: something
}

func ExampleBodyString() {
	req, _ := NewPost(BodyString("something"))
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	// Output: something
}

func ExampleBodyForm() {
	req, _ := NewPost(BodyForm{"something": {"here"}})
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	// Output: something=here
}

func ExampleBodyJSON() {
	req, _ := NewPost(BodyJSON(map[string]interface{}{"something": "here"}))
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	// Output: {"something":"here"}
}

func ExampleBodyXML() {
	req, _ := NewPost(BodyXML("something"))
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	// Output: <string>something</string>
}
