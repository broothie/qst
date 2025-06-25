package qst_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	pkgurl "net/url"
	"os"
	"testing"
	"time"

	"github.com/broothie/qst"
	"github.com/stretchr/testify/assert"
)

func TestOptions_errors(t *testing.T) {
	t.Run("URL error", func(t *testing.T) {
		_, err := qst.NewPost("https://breakfast.com/api/cereals", qst.URL("%"))
		assert.EqualError(t, err, `failed to apply option 0: parse "%": invalid URL escape "%"`)
	})

	t.Run("BodyJSON error", func(t *testing.T) {
		_, err := qst.NewPost("https://breakfast.com/api/cereals", qst.BodyJSON(make(chan struct{})))
		assert.EqualError(t, err, "failed to apply option 0: json: unsupported type: chan struct {}")
	})

	t.Run("BodyXML error", func(t *testing.T) {
		_, err := qst.NewPost("https://breakfast.com/api/cereals", qst.BodyXML(make(chan struct{})))
		assert.EqualError(t, err, "failed to apply option 0: xml: unsupported type: chan struct {}")
	})

	t.Run("Dump read error", func(t *testing.T) {
		_, err := qst.NewGet("https://breakfast.com/api/cereals",
			qst.BodyReader(broken{}),
			qst.Dump(os.Stdout),
		)

		assert.EqualError(t, err, "failed to apply option 1: broken")
	})

	t.Run("Dump write error", func(t *testing.T) {
		_, err := qst.NewGet("https://breakfast.com/api/cereals",
			qst.Dump(broken{}),
		)

		assert.EqualError(t, err, "failed to apply option 0: broken")
	})
}

func TestComplexFieldOptions(t *testing.T) {
	t.Run("MultipartForm", func(t *testing.T) {
		form := &multipart.Form{
			Value: map[string][]string{"key": {"value"}},
			File:  map[string][]*multipart.FileHeader{},
		}
		
		request, err := qst.NewPost("https://breakfast.com/api/cereals",
			qst.MultipartForm(form),
		)
		
		assert.NoError(t, err)
		assert.Equal(t, form, request.MultipartForm)
		assert.Equal(t, "value", request.MultipartForm.Value["key"][0])
	})
	
	t.Run("MultipartFormValue", func(t *testing.T) {
		request, err := qst.NewPost("https://breakfast.com/api/cereals",
			qst.MultipartFormValue("name", "Frosted Flakes"),
			qst.MultipartFormValue("brand", "Kellogg's"),
			qst.MultipartFormValue("name", "Corn Flakes"), // Multiple values for same key
		)
		
		assert.NoError(t, err)
		assert.NotNil(t, request.MultipartForm)
		assert.Equal(t, "Frosted Flakes", request.MultipartForm.Value["name"][0])
		assert.Equal(t, "Corn Flakes", request.MultipartForm.Value["name"][1])
		assert.Equal(t, "Kellogg's", request.MultipartForm.Value["brand"][0])
	})
	
	t.Run("Form and PostForm", func(t *testing.T) {
		form := pkgurl.Values{"name": {"Lucky Charms"}}
		postForm := pkgurl.Values{"type": {"cereal"}}
		
		request, err := qst.NewPost("https://breakfast.com/api/cereals",
			qst.Form(form),
			qst.PostForm(postForm),
		)
		
		assert.NoError(t, err)
		assert.Equal(t, "Lucky Charms", request.Form.Get("name"))
		assert.Equal(t, "cereal", request.PostForm.Get("type"))
	})
	
	t.Run("FormValue and PostFormValue", func(t *testing.T) {
		request, err := qst.NewPost("https://breakfast.com/api/cereals",
			qst.FormValue("name", "Captain Crunch"),
			qst.FormValue("brand", "Quaker"),
			qst.PostFormValue("type", "cereal"),
			qst.PostFormValue("sweetness", "high"),
		)
		
		assert.NoError(t, err)
		assert.Equal(t, "Captain Crunch", request.Form.Get("name"))
		assert.Equal(t, "Quaker", request.Form.Get("brand"))
		assert.Equal(t, "cereal", request.PostForm.Get("type"))
		assert.Equal(t, "high", request.PostForm.Get("sweetness"))
	})
	
	t.Run("Trailer", func(t *testing.T) {
		trailer := http.Header{"X-Checksum": {"abc123"}}
		
		request, err := qst.NewPost("https://breakfast.com/api/cereals",
			qst.Trailer(trailer),
		)
		
		assert.NoError(t, err)
		assert.Equal(t, trailer, request.Trailer)
		assert.Equal(t, "abc123", request.Trailer.Get("X-Checksum"))
	})
	
	t.Run("TrailerHeader", func(t *testing.T) {
		request, err := qst.NewPost("https://breakfast.com/api/cereals",
			qst.TrailerHeader("X-Checksum", "def456"),
			qst.TrailerHeader("X-Signature", "ghi789"),
		)
		
		assert.NoError(t, err)
		assert.Equal(t, "def456", request.Trailer.Get("X-Checksum"))
		assert.Equal(t, "ghi789", request.Trailer.Get("X-Signature"))
	})
	
	t.Run("TransferEncoding", func(t *testing.T) {
		encodings := []string{"chunked", "gzip"}
		
		request, err := qst.NewPost("https://breakfast.com/api/cereals",
			qst.TransferEncoding(encodings),
		)
		
		assert.NoError(t, err)
		assert.Equal(t, encodings, request.TransferEncoding)
	})
	
	t.Run("TransferEncodingAppend", func(t *testing.T) {
		request, err := qst.NewPost("https://breakfast.com/api/cereals",
			qst.TransferEncodingAppend("chunked"),
			qst.TransferEncodingAppend("gzip"),
		)
		
		assert.NoError(t, err)
		assert.Equal(t, []string{"chunked", "gzip"}, request.TransferEncoding)
	})
	
	t.Run("GetBody", func(t *testing.T) {
		getBodyFunc := func() (io.ReadCloser, error) {
			return ioutil.NopCloser(bytes.NewBufferString("test body")), nil
		}
		
		request, err := qst.NewPost("https://breakfast.com/api/cereals",
			qst.GetBody(getBodyFunc),
		)
		
		assert.NoError(t, err)
		assert.NotNil(t, request.GetBody)
		
		body, err := request.GetBody()
		assert.NoError(t, err)
		content, err := ioutil.ReadAll(body)
		assert.NoError(t, err)
		assert.Equal(t, "test body", string(content))
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

func ExampleMultipartForm() {
	form := &multipart.Form{
		Value: map[string][]string{"name": {"Cheerios"}},
		File:  map[string][]*multipart.FileHeader{},
	}
	
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.MultipartForm(form),
	)

	fmt.Println(request.MultipartForm.Value["name"][0])
	// Output: Cheerios
}

func ExampleMultipartFormValue() {
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.MultipartFormValue("name", "Frosted Flakes"),
		qst.MultipartFormValue("brand", "Kellogg's"),
	)

	fmt.Println(request.MultipartForm.Value["name"][0])
	fmt.Println(request.MultipartForm.Value["brand"][0])
	// Output: Frosted Flakes
	// Kellogg's
}

func ExampleForm() {
	form := pkgurl.Values{"name": {"Lucky Charms"}, "marshmallows": {"yes"}}
	
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.Form(form),
	)

	fmt.Println(request.Form.Get("name"))
	fmt.Println(request.Form.Get("marshmallows"))
	// Output: Lucky Charms
	// yes
}

func ExamplePostForm() {
	form := pkgurl.Values{"name": {"Honey Nut Cheerios"}}
	
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.PostForm(form),
	)

	fmt.Println(request.PostForm.Get("name"))
	// Output: Honey Nut Cheerios
}

func ExampleFormValue() {
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.FormValue("name", "Captain Crunch"),
		qst.FormValue("crunch", "high"),
	)

	fmt.Println(request.Form.Get("name"))
	fmt.Println(request.Form.Get("crunch"))
	// Output: Captain Crunch
	// high
}

func ExamplePostFormValue() {
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.PostFormValue("name", "Fruit Loops"),
		qst.PostFormValue("colors", "many"),
	)

	fmt.Println(request.PostForm.Get("name"))
	fmt.Println(request.PostForm.Get("colors"))
	// Output: Fruit Loops
	// many
}

func ExampleTrailer() {
	trailer := http.Header{"X-Checksum": {"abc123"}}
	
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.Trailer(trailer),
	)

	fmt.Println(request.Trailer.Get("X-Checksum"))
	// Output: abc123
}

func ExampleTrailerHeader() {
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.TrailerHeader("X-Checksum", "def456"),
		qst.TrailerHeader("X-Signature", "ghi789"),
	)

	fmt.Println(request.Trailer.Get("X-Checksum"))
	fmt.Println(request.Trailer.Get("X-Signature"))
	// Output: def456
	// ghi789
}

func ExampleTransferEncoding() {
	encodings := []string{"chunked", "gzip"}
	
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.TransferEncoding(encodings),
	)

	fmt.Println(request.TransferEncoding[0])
	fmt.Println(request.TransferEncoding[1])
	// Output: chunked
	// gzip
}

func ExampleTransferEncodingAppend() {
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.TransferEncodingAppend("chunked"),
		qst.TransferEncodingAppend("gzip"),
	)

	fmt.Println(request.TransferEncoding[0])
	fmt.Println(request.TransferEncoding[1])
	// Output: chunked
	// gzip
}

func ExampleGetBody() {
	getBodyFunc := func() (io.ReadCloser, error) {
		return ioutil.NopCloser(bytes.NewBufferString("test body")), nil
	}
	
	request, _ := qst.NewPost("https://breakfast.com/api/cereals",
		qst.GetBody(getBodyFunc),
	)

	body, _ := request.GetBody()
	content, _ := ioutil.ReadAll(body)
	fmt.Println(string(content))
	// Output: test body
}
