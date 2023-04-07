package qst_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/broothie/qst"
	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	t.Run("BodyJSON error", func(t *testing.T) {
		_, err := qst.NewPost("https://breakfast.com/api", qst.BodyJSON(make(chan struct{})))
		assert.EqualError(t, err, "json: unsupported type: chan struct {}")
	})

	t.Run("BodyXML error", func(t *testing.T) {
		_, err := qst.NewPost("https://breakfast.com/api", qst.BodyXML(make(chan struct{})))
		assert.EqualError(t, err, "xml: unsupported type: chan struct {}")
	})

	t.Run("Context", func(t *testing.T) {
		type keyType struct{}
		var key keyType

		request, err := qst.NewPost("https://breakfast.com/api",
			qst.Context(context.WithValue(context.TODO(), key, "milk")),
		)

		assert.NoError(t, err)
		assert.Equal(t, "milk", request.Context().Value(key))
	})

	t.Run("ContextValue", func(t *testing.T) {
		type keyType struct{}
		var key keyType

		request, err := qst.NewPost("https://breakfast.com/api", qst.ContextValue(key, "milk"))

		assert.NoError(t, err)

		assert.Equal(t, "milk", request.Context().Value(key))
	})
}

func ExamplePath() {
	request, _ := qst.NewGet("https://breakfast.com/api/",
		qst.Path("/cereals", "1234/variants", "frosted"),
	)

	fmt.Println(request.URL.Path)
	// Output: /api/cereals/1234/variants/frosted
}

func ExampleUsername() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.Username("TonyTheTiger"),
	)

	fmt.Println(request.URL)
	// Output: https://TonyTheTiger@breakfast.com/api/cereals
}

func ExampleUserPassword() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.UserPassword("TonyTheTiger", "grrreat"),
	)

	fmt.Println(request.URL)
	// Output: https://TonyTheTiger:grrreat@breakfast.com/api/cereals
}

func ExampleQuery() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.Query("page", "10"),
	)

	fmt.Println(request.URL.Query().Encode())
	// Output: page=10
}

func ExampleQueries() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.Queries{
			"page":  {"10"},
			"limit": {"50"},
		},
	)

	fmt.Println(request.URL.Query().Encode())
	// Output: limit=50&page=10
}

func ExampleHeader() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.Header("grain", "oats"),
	)

	fmt.Println(request.Header.Get("grain"))
	// Output: oats
}

func ExampleHeaders() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.Headers{
			"grain": {"oats"},
			"style": {"toasted"},
		},
	)

	fmt.Println(request.Header)
	// Output: map[Grain:[oats] Style:[toasted]]
}

func ExampleCookie() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.Cookie(&http.Cookie{
			Name:    "cookie-crisp",
			Value:   "COOOOKIE CRISP!",
			Path:    "/",
			Expires: time.Now().Add(time.Hour),
		}),
	)

	fmt.Println(request.Cookie("cookie-crisp"))
	// Output: cookie-crisp="COOOOKIE CRISP!" <nil>
}

func ExampleAuthorization() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.Authorization("c0rnfl@k3s"),
	)

	fmt.Println(request.Header.Get("Authorization"))
	// Output: c0rnfl@k3s
}

func ExampleBasicAuth() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.BasicAuth("TonyTheTiger", "grrreat"),
	)

	fmt.Println(request.Header.Get("Authorization"))
	// Output: Basic VG9ueVRoZVRpZ2VyOmdycnJlYXQ=
}

func ExampleBearerAuth() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.BearerAuth("c0rnfl@k3s"),
	)

	fmt.Println(request.Header.Get("Authorization"))
	// Output: Bearer c0rnfl@k3s
}

func ExampleContextValue() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.ContextValue("frosted", true),
	)

	fmt.Println(request.Context().Value("frosted"))
	// Output: true
}

func ExampleBody() {
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.Body(ioutil.NopCloser(bytes.NewBufferString("Part of a complete breakfast."))),
	)

	body, _ := ioutil.ReadAll(request.Body)
	fmt.Println(string(body))
	// Output: Part of a complete breakfast.
}

func ExampleBodyBytes() {
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.BodyBytes([]byte("Part of a complete breakfast.")),
	)

	body, _ := ioutil.ReadAll(request.Body)
	fmt.Println(string(body))
	// Output: Part of a complete breakfast.
}

func ExampleBodyString() {
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.BodyString("Part of a complete breakfast."),
	)

	body, _ := ioutil.ReadAll(request.Body)
	fmt.Println(string(body))
	// Output: Part of a complete breakfast.
}

func ExampleBodyForm() {
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.BodyForm{"name": {"Grape Nuts"}},
	)

	body, _ := ioutil.ReadAll(request.Body)
	fmt.Println(string(body))
	// Output: name=Grape+Nuts
}

func ExampleBodyJSON() {
	request, _ := qst.NewPost("https://breakfast.com/api",
		qst.BodyJSON(map[string]string{"name": "Rice Krispies"}),
	)

	body, _ := ioutil.ReadAll(request.Body)
	fmt.Println(string(body))
	// Output: {"name":"Rice Krispies"}
}

func ExampleBodyXML() {
	request, _ := qst.NewPost("https://breakfast.com/api",
		qst.BodyXML("Part of a complete breakfast."),
	)
	body, _ := ioutil.ReadAll(request.Body)
	fmt.Println(string(body))
	// Output: <string>Part of a complete breakfast.</string>
}
