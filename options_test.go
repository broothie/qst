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
		_, err := qst.NewPost("https://breakfast.com/api/cereals", qst.URL("%"))
		assert.EqualError(t, err, `parse %: invalid URL escape "%"`)
	})

	t.Run("BodyJSON error", func(t *testing.T) {
		_, err := qst.NewPost("https://breakfast.com/api/cereals", qst.BodyJSON(make(chan struct{})))
		assert.EqualError(t, err, "json: unsupported type: chan struct {}")
	})

	t.Run("BodyXML error", func(t *testing.T) {
		_, err := qst.NewPost("https://breakfast.com/api/cereals", qst.BodyXML(make(chan struct{})))
		assert.EqualError(t, err, "xml: unsupported type: chan struct {}")
	})

	t.Run("Dump read error", func(t *testing.T) {
		_, err := qst.NewGet("https://breakfast.com/api/cereals",
			qst.BodyReader(broken{}),
			qst.Dump(os.Stdout),
		)

		assert.EqualError(t, err, "broken")
	})

	t.Run("Dump write error", func(t *testing.T) {
		_, err := qst.NewGet("https://breakfast.com/api/cereals",
			qst.Dump(broken{}),
		)

		assert.EqualError(t, err, "broken")
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

func ExampleAccept() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.Accept("application/json"),
	)

	fmt.Println(request.Header.Get("Accept"))
	// Output: application/json
}

func ExampleContentType() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.ContentType("application/json"),
	)

	fmt.Println(request.Header.Get("Content-Type"))
	// Output: application/json
}

func ExampleReferer() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.Referer("https://breakfast.com"),
	)

	fmt.Println(request.Header.Get("Referer"))
	// Output: https://breakfast.com
}

func ExampleUserAgent() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.UserAgent("qst"),
	)

	fmt.Println(request.Header.Get("User-Agent"))
	// Output: qst
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

func ExampleTokenAuth() {
	request, _ := qst.NewGet("https://breakfast.com/api/cereals",
		qst.TokenAuth("c0rnfl@k3s"),
	)

	fmt.Println(request.Header.Get("Authorization"))
	// Output: Token c0rnfl@k3s
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
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.BodyJSON(map[string]string{"name": "Rice Krispies"}),
	)

	body, _ := ioutil.ReadAll(request.Body)
	fmt.Println(string(body))
	// Output: {"name":"Rice Krispies"}
}

func ExampleBodyXML() {
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.BodyXML("Part of a complete breakfast."),
	)
	body, _ := ioutil.ReadAll(request.Body)
	fmt.Println(string(body))
	// Output: <string>Part of a complete breakfast.</string>
}

func TestDump(t *testing.T) {
	var buffer bytes.Buffer
	qst.NewGet("https://breakfast.com/api/cereals",
		qst.BodyString("Part of a complete breakfast."),
		qst.Dump(&buffer),
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
