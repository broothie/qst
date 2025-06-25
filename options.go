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

// WithRawURL applies the URL to the *http.Request.
func WithRawURL(url *pkgurl.URL) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		request.URL = url
		return request, nil
	})
}

// WithURL applies a url string to the *http.Request.
func WithURL(url string) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		u, err := pkgurl.Parse(url)
		if err != nil {
			return nil, err
		}

		return WithRawURL(u).Apply(request)
	})
}

// WithScheme applies the scheme to the *http.Request URL.
func WithScheme(scheme string) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		request.URL.Scheme = scheme
		return request, nil
	})
}

// WithUser applies the Userinfo to the *http.Request URL User.
func WithUser(user *pkgurl.Userinfo) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		request.URL.User = user
		return request, nil
	})
}

// WithUsername applies the username to *http.Request URL User.
func WithUsername(username string) option.Option[*http.Request] {
	return WithUser(pkgurl.User(username))
}

// WithUserPassword applies the username and password to *http.Request URL User.
func WithUserPassword(username, password string) option.Option[*http.Request] {
	return WithUser(pkgurl.UserPassword(username, password))
}

// WithHost applies the host to the *http.Request and *http.Request URL.
func WithHost(host string) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		request.Host = host
		request.URL.Host = host
		return request, nil
	})
}

// WithPath joins the segments with path.Join, and appends the result to the *http.Request URL.
func WithPath(segments ...string) option.Option[*http.Request] {
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

// WithQuery applies a key/value pair to the query parameters of the *http.Request.
func WithQuery(key, value string) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		query := request.URL.Query()
		query.Add(key, value)
		request.URL.RawQuery = query.Encode()

		return request, nil
	})
}

// WithQueries applies multiple key/value pairs to the query parameters of the *http.Request. It wraps url.Values.
type WithQueries pkgurl.Values

// Apply applies the WithQueries to the *http.Request.
func (q WithQueries) Apply(request *http.Request) (*http.Request, error) {
	var options []option.Option[*http.Request]
	for key, values := range q {
		for _, value := range values {
			options = append(options, WithQuery(key, value))
		}
	}

	return option.Apply(request, options...)
}

// WithHeader applies a key/value pair to the headers of the *http.Request, retaining the existing headers for the key.
func WithHeader(key, value string) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		request.Header.Add(key, value)
		return request, nil
	})
}

// WithHeaders applies multiple key/value pairs to the headers of the *http.Request. It wraps http.Header.
type WithHeaders http.Header

// Apply applies the WithHeaders to the *http.Request.
func (h WithHeaders) Apply(request *http.Request) (*http.Request, error) {
	var options []option.Option[*http.Request]
	for key, values := range h {
		for _, value := range values {
			options = append(options, WithHeader(key, value))
		}
	}

	return option.Apply(request, options...)
}

// WithAcceptHeader applies an "Accept" header to the *http.Request.
func WithAcceptHeader(accept string) option.Option[*http.Request] {
	return WithHeader("Accept", accept)
}

// WithContentTypeHeader applies a "Content-Type" to the *http.Request.
func WithContentTypeHeader(contentType string) option.Option[*http.Request] {
	return WithHeader("Content-Type", contentType)
}

// WithRefererHeader applies a "Referer" header to the *http.Request.
func WithRefererHeader(referer string) option.Option[*http.Request] {
	return WithHeader("Referer", referer)
}

// WithUserAgentHeader applies a "User-Agent" header to the *http.Request.
func WithUserAgentHeader(userAgent string) option.Option[*http.Request] {
	return WithHeader("User-Agent", userAgent)
}

// WithCookie applies a cookie to the *http.Request.
func WithCookie(cookie *http.Cookie) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		request.AddCookie(cookie)
		return request, nil
	})
}

// WithAuthorizationHeader applies an "Authorization" header to the *http.Request.
func WithAuthorizationHeader(authorization string) option.Option[*http.Request] {
	return WithHeader("Authorization", authorization)
}

// WithBasicAuth applies a username and password basic auth header.
func WithBasicAuth(username, password string) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		request.SetBasicAuth(username, password)
		return request, nil
	})
}

// WithTokenAuth applies an "Authorization: Token <token>" header to the *http.Request.
func WithTokenAuth(token string) option.Option[*http.Request] {
	return WithAuthorizationHeader(fmt.Sprintf("Token %s", token))
}

// WithBearerAuth applies an "Authorization: Bearer <token>" header to the *http.Request.
func WithBearerAuth(token string) option.Option[*http.Request] {
	return WithAuthorizationHeader(fmt.Sprintf("Bearer %s", token))
}

// WithContext applies a context.Context to the *http.Request.
func WithContext(ctx context.Context) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		return request.WithContext(ctx), nil
	})
}

// WithContextValue applies a context key/value pair to the *http.Request.
func WithContextValue(key, value interface{}) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		return WithContext(context.WithValue(request.Context(), key, value)).Apply(request)
	})
}

// WithBody applies an io.ReadCloser to the *http.Request body.
func WithBody(body io.ReadCloser) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		request.Body = body
		return request, nil
	})
}

// WithBodyReader applies an io.Reader to the *http.Request body.
func WithBodyReader(body io.Reader) option.Option[*http.Request] {
	return WithBody(ioutil.NopCloser(body))
}

// WithBodyBytes applies a slice of bytes to the *http.Request body.
func WithBodyBytes(body []byte) option.Option[*http.Request] {
	return WithBodyReader(bytes.NewBuffer(body))
}

// WithBodyString applies a string to the *http.Request body.
func WithBodyString(body string) option.Option[*http.Request] {
	return WithBodyBytes([]byte(body))
}

// WithBodyForm URL-encodes multiple key/value pairs and applies the result to the *http.Request body.
type WithBodyForm pkgurl.Values

// Apply URL-encodes the WithBodyForm and applies the result to the *http.Request body.
func (f WithBodyForm) Apply(request *http.Request) (*http.Request, error) {
	return option.Apply(request,
		WithContentTypeHeader("application/x-www-form-urlencoded"),
		WithBodyString(pkgurl.Values(f).Encode()),
	)
}

// WithBodyJSON encodes an object as JSON and applies it to the *http.Request body.
func WithBodyJSON(v interface{}) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		body := new(bytes.Buffer)
		if err := json.NewEncoder(body).Encode(v); err != nil {
			return nil, err
		}

		return option.Apply(request,
			WithContentTypeHeader("application/json"),
			WithBodyReader(body),
		)
	})
}

// WithBodyXML encodes an object as XML and applies it to the *http.Request body.
func WithBodyXML(v interface{}) option.Option[*http.Request] {
	return option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
		body := new(bytes.Buffer)
		if err := xml.NewEncoder(body).Encode(v); err != nil {
			return nil, err
		}

		return option.Apply(request,
			WithContentTypeHeader("application/xml"),
			WithBodyReader(body),
		)
	})
}

// WithDump writes the request to w.
func WithDump(w io.Writer) option.Option[*http.Request] {
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
