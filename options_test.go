package qst_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/broothie/qst"
	"github.com/stretchr/testify/assert"
)

func TestOptions_errors(t *testing.T) {
	t.Run("URL error", func(t *testing.T) {
		_, err := qst.NewPost("https://breakfast.com/api/cereals", qst.WithURL("%"))
		assert.EqualError(t, err, `failed to apply option 0: parse "%": invalid URL escape "%"`)
	})

	t.Run("BodyJSON error", func(t *testing.T) {
		_, err := qst.NewPost("https://breakfast.com/api/cereals", qst.WithBodyJSON(make(chan struct{})))
		assert.EqualError(t, err, "failed to apply option 0: json: unsupported type: chan struct {}")
	})

	t.Run("BodyXML error", func(t *testing.T) {
		_, err := qst.NewPost("https://breakfast.com/api/cereals", qst.WithBodyXML(make(chan struct{})))
		assert.EqualError(t, err, "failed to apply option 0: xml: unsupported type: chan struct {}")
	})

	t.Run("Dump read error", func(t *testing.T) {
		_, err := qst.NewGet("https://breakfast.com/api/cereals",
			qst.WithBodyReader(broken{}),
			qst.WithDump(os.Stdout),
		)

		assert.EqualError(t, err, "failed to apply option 1: broken")
	})

	t.Run("Dump write error", func(t *testing.T) {
		_, err := qst.NewGet("https://breakfast.com/api/cereals",
			qst.WithDump(broken{}),
		)

		assert.EqualError(t, err, "failed to apply option 0: broken")
	})
}

func ExampleWithPath() {
	request, _ := qst.NewGet("https://breakfast.com/api/",
		qst.WithPath("/cereals", "1234/variants", "frosted"),
	)

	fmt.Println(request.URL.Path)
	// Output: /api/cereals/1234/variants/frosted
}

func ExampleWithUsername() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.WithUsername("TonyTheTiger"),
	)

	fmt.Println(request.URL)
	// Output: https://TonyTheTiger@breakfast.com/api/cereals
}

func ExampleWithUserPassword() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.WithUserPassword("TonyTheTiger", "grrreat"),
	)

	fmt.Println(request.URL)
	// Output: https://TonyTheTiger:grrreat@breakfast.com/api/cereals
}

func ExampleWithQuery() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.WithQuery("page", "10"),
	)

	fmt.Println(request.URL.Query().Encode())
	// Output: page=10
}

func ExampleWithQueries() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.WithQueries{
			"page":  {"10"},
			"limit": {"50"},
		},
	)

	fmt.Println(request.URL.Query().Encode())
	// Output: limit=50&page=10
}

func ExampleWithHeader() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.WithHeader("grain", "oats"),
	)

	fmt.Println(request.Header.Get("grain"))
	// Output: oats
}

func ExampleWithHeaders() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.WithHeaders{
			"grain": {"oats"},
			"style": {"toasted"},
		},
	)

	fmt.Println(request.Header)
	// Output: map[Grain:[oats] Style:[toasted]]
}

func ExampleWithAcceptHeader() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.WithAcceptHeader("application/json"),
	)

	fmt.Println(request.Header.Get("Accept"))
	// Output: application/json
}

func ExampleWithContentTypeHeader() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.WithContentTypeHeader("application/json"),
	)

	fmt.Println(request.Header.Get("Content-Type"))
	// Output: application/json
}

func ExampleWithRefererHeader() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.WithRefererHeader("https://breakfast.com"),
	)

	fmt.Println(request.Header.Get("Referer"))
	// Output: https://breakfast.com
}

func ExampleWithUserAgentHeader() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.WithUserAgentHeader("qst"),
	)

	fmt.Println(request.Header.Get("User-Agent"))
	// Output: qst
}

func ExampleWithCookie() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.WithCookie(&http.Cookie{
			Name:    "cookie-crisp",
			Value:   "COOOOKIE CRISP!",
			Path:    "/",
			Expires: time.Now().Add(time.Hour),
		}),
	)

	fmt.Println(request.Cookie("cookie-crisp"))
	// Output: cookie-crisp="COOOOKIE CRISP!" <nil>
}

func ExampleWithAuthorizationHeader() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.WithAuthorizationHeader("c0rnfl@k3s"),
	)

	fmt.Println(request.Header.Get("Authorization"))
	// Output: c0rnfl@k3s
}

func ExampleWithBasicAuth() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.WithBasicAuth("TonyTheTiger", "grrreat"),
	)

	fmt.Println(request.Header.Get("Authorization"))
	// Output: Basic VG9ueVRoZVRpZ2VyOmdycnJlYXQ=
}

func ExampleWithTokenAuth() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.WithTokenAuth("c0rnfl@k3s"),
	)

	fmt.Println(request.Header.Get("Authorization"))
	// Output: Token c0rnfl@k3s
}

func ExampleWithBearerAuth() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.WithBearerAuth("c0rnfl@k3s"),
	)

	fmt.Println(request.Header.Get("Authorization"))
	// Output: Bearer c0rnfl@k3s
}

func ExampleWithContextValue() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.WithContextValue("frosted", true),
	)

	fmt.Println(request.Context().Value("frosted"))
	// Output: true
}

func ExampleWithBody() {
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.WithBody(ioutil.NopCloser(bytes.NewBufferString("Part of a complete breakfast."))),
	)

	body, _ := ioutil.ReadAll(request.Body)
	fmt.Println(string(body))
	// Output: Part of a complete breakfast.
}

func ExampleWithBodyBytes() {
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.WithBodyBytes([]byte("Part of a complete breakfast.")),
	)

	body, _ := ioutil.ReadAll(request.Body)
	fmt.Println(string(body))
	// Output: Part of a complete breakfast.
}

func ExampleWithBodyString() {
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.WithBodyString("Part of a complete breakfast."),
	)

	body, _ := ioutil.ReadAll(request.Body)
	fmt.Println(string(body))
	// Output: Part of a complete breakfast.
}

func ExampleWithBodyForm() {
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.WithBodyForm{"name": {"Grape Nuts"}},
	)

	body, _ := ioutil.ReadAll(request.Body)
	fmt.Println(string(body))
	// Output: name=Grape+Nuts
}

func ExampleWithBodyJSON() {
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.WithBodyJSON(map[string]string{"name": "Rice Krispies"}),
	)

	body, _ := ioutil.ReadAll(request.Body)
	fmt.Println(string(body))
	// Output: {"name":"Rice Krispies"}
}

func ExampleWithBodyXML() {
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.WithBodyXML("Part of a complete breakfast."),
	)
	body, _ := ioutil.ReadAll(request.Body)
	fmt.Println(string(body))
	// Output: <string>Part of a complete breakfast.</string>
}

func TestDump(t *testing.T) {
	var buffer bytes.Buffer
	qst.NewGet("https://breakfast.com/api/cereals",
		qst.WithBodyString("Part of a complete breakfast."),
		qst.WithDump(&buffer),
	)

	expected := "" +
		"GET /api/cereals HTTP/1.1\r\n" +
		"Host: breakfast.com\r\n" +
		"\r\n" +
		"Part of a complete breakfast."

	assert.Equal(t, expected, buffer.String())
}

type broken struct{}

func (broken) Write([]byte) (int, error) {
	return 0, errors.New("broken")
}

func (broken) Read([]byte) (int, error) {
	return 0, errors.New("broken")
}
