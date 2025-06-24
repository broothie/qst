package qst

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	pkgurl "net/url"
	pkgpath "path"

	"github.com/broothie/option"
)

// RawURL applies the URL to the *http.Request.
func RawURL(url *pkgurl.URL) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		request.URL = url
		return request, nil
	})
}

// URL applies a url string to the *http.Request.
func URL(url string) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		u, err := pkgurl.Parse(url)
		if err != nil {
			return nil, err
		}

		return RawURL(u).Apply(request)
	})
}

// Scheme applies the scheme to the *http.Request URL.
func Scheme(scheme string) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		request.URL.Scheme = scheme
		return request, nil
	})
}

// User applies the Userinfo to the *http.Request URL User.
func User(user *pkgurl.Userinfo) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		request.URL.User = user
		return request, nil
	})
}

// Username applies the username to *http.Request URL User.
func Username(username string) option.Option[*http.Request] {
	return User(pkgurl.User(username))
}

// UserPassword applies the username and password to *http.Request URL User.
func UserPassword(username, password string) option.Option[*http.Request] {
	return User(pkgurl.UserPassword(username, password))
}

// Host applies the host to the *http.Request and *http.Request URL.
func Host(host string) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		request.Host = host
		request.URL.Host = host
		return request, nil
	})
}

// Path joins the segments with path.Join, and appends the result to the *http.Request URL.
func Path(segments ...string) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		elem := []string{request.URL.Path}
		elem = append(elem, segments...)
		path := pkgpath.Join(elem...)
		if !pkgpath.IsAbs(path) {
			path = fmt.Sprintf("/%s", path)
		}

		request.URL.Path = path
		return request, nil
	})
}

// Query applies a key/value pair to the query parameters of the *http.Request.
func Query(key, value string) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		query := request.URL.Query()
		query.Add(key, value)
		request.URL.RawQuery = query.Encode()

		return request, nil
	})
}

// Queries applies multiple key/value pairs to the query parameters of the *http.Request. It wraps url.Values.
type Queries pkgurl.Values

// Apply applies the Queries to the *http.Request.
func (q Queries) Apply(request *http.Request) (*http.Request, error) {
	var options []option.Option[*http.Request]
	for key, values := range q {
		for _, value := range values {
			options = append(options, Query(key, value))
		}
	}

	return option.Apply(request, options...)
}

// Header applies a key/value pair to the headers of the *http.Request, retaining the existing headers for the key.
func Header(key, value string) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		request.Header.Add(key, value)
		return request, nil
	})
}

// Headers applies multiple key/value pairs to the headers of the *http.Request. It wraps http.Header.
type Headers http.Header

// Apply applies the Headers to the *http.Request.
func (h Headers) Apply(request *http.Request) (*http.Request, error) {
	var options []option.Option[*http.Request]
	for key, values := range h {
		for _, value := range values {
			options = append(options, Header(key, value))
		}
	}

	return option.Apply(request, options...)
}

// Accept applies an "Accept" header to the *http.Request.
func Accept(accept string) option.Option[*http.Request] {
	return Header("Accept", accept)
}

// ContentType applies a "Content-Type" to the *http.Request.
func ContentType(contentType string) option.Option[*http.Request] {
	return Header("Content-Type", contentType)
}

// Referer applies a "Referer" header to the *http.Request.
func Referer(referer string) option.Option[*http.Request] {
	return Header("Referer", referer)
}

// UserAgent applies a "User-Agent" header to the *http.Request.
func UserAgent(userAgent string) option.Option[*http.Request] {
	return Header("User-Agent", userAgent)
}

// Cookie applies a cookie to the *http.Request.
func Cookie(cookie *http.Cookie) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		request.AddCookie(cookie)
		return request, nil
	})
}

// Authorization applies an "Authorization" header to the *http.Request.
func Authorization(authorization string) option.Option[*http.Request] {
	return Header("Authorization", authorization)
}

// BasicAuth applies a username and password basic auth header.
func BasicAuth(username, password string) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		request.SetBasicAuth(username, password)
		return request, nil
	})
}

// TokenAuth applies an "Authorization: Token <token>" header to the *http.Request.
func TokenAuth(token string) option.Option[*http.Request] {
	return Authorization(fmt.Sprintf("Token %s", token))
}

// BearerAuth applies an "Authorization: Bearer <token>" header to the *http.Request.
func BearerAuth(token string) option.Option[*http.Request] {
	return Authorization(fmt.Sprintf("Bearer %s", token))
}

// Context applies a context.Context to the *http.Request.
func Context(ctx context.Context) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		return request.WithContext(ctx), nil
	})
}

// ContextValue applies a context key/value pair to the *http.Request.
func ContextValue(key, value interface{}) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		return Context(context.WithValue(request.Context(), key, value)).Apply(request)
	})
}

// Body applies an io.ReadCloser to the *http.Request body.
func Body(body io.ReadCloser) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		request.Body = body
		return request, nil
	})
}

// BodyReader applies an io.Reader to the *http.Request body.
func BodyReader(body io.Reader) option.Option[*http.Request] {
	return Body(ioutil.NopCloser(body))
}

// BodyBytes applies a slice of bytes to the *http.Request body.
func BodyBytes(body []byte) option.Option[*http.Request] {
	return BodyReader(bytes.NewBuffer(body))
}

// BodyString applies a string to the *http.Request body.
func BodyString(body string) option.Option[*http.Request] {
	return BodyBytes([]byte(body))
}

// BodyForm URL-encodes multiple key/value pairs and applies the result to the *http.Request body.
type BodyForm pkgurl.Values

// Apply URL-encodes the BodyForm and applies the result to the *http.Request body.
func (f BodyForm) Apply(request *http.Request) (*http.Request, error) {
	return option.Apply(request,
		ContentType("application/x-www-form-urlencoded"),
		BodyString(pkgurl.Values(f).Encode()),
	)
}

// BodyJSON encodes an object as JSON and applies it to the *http.Request body.
func BodyJSON(v interface{}) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		body := new(bytes.Buffer)
		if err := json.NewEncoder(body).Encode(v); err != nil {
			return nil, err
		}

		return option.Apply(request,
			ContentType("application/json"),
			BodyReader(body),
		)
	})
}

// BodyXML encodes an object as XML and applies it to the *http.Request body.
func BodyXML(v interface{}) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		body := new(bytes.Buffer)
		if err := xml.NewEncoder(body).Encode(v); err != nil {
			return nil, err
		}

		return option.Apply(request,
			ContentType("application/xml"),
			BodyReader(body),
		)
	})
}

// Dump writes the request to w.
func Dump(w io.Writer) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		dump, err := httputil.DumpRequest(request, true)
		if err != nil {
			return nil, err
		}

		if _, err := w.Write(dump); err != nil {
			return nil, err
		}

		return request, nil
	})
}
