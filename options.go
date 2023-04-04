package qst

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"path"
)

// Option is an option for building *http.Requests.
type Option interface {
	Apply(request *http.Request) (*http.Request, error)
}

// Pipeline is a collection of options, which can be applied as a whole.
type Pipeline []Option

// Apply applies the Pipeline to the *http.Request.
func (p Pipeline) Apply(request *http.Request) (*http.Request, error) {
	for _, option := range p {
		var err error
		request, err = option.Apply(request.Clone(request.Context()))
		if err != nil {
			return nil, err
		}
	}

	return request, nil
}

// OptionFunc is a function form of Option.
type OptionFunc func(request *http.Request) (*http.Request, error)

// Apply applies the OptionFunc to the *http.Request.
func (f OptionFunc) Apply(request *http.Request) (*http.Request, error) { return f(request) }

// Apply applies the Options to the *http.Request
func Apply(request *http.Request, options ...Option) (*http.Request, error) {
	return Pipeline(options).Apply(request)
}

func RawURL(url *neturl.URL) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		request.URL = url
		return request, nil
	})
}

func URL(url string) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		u, err := neturl.Parse(url)
		if err != nil {
			return nil, err
		}

		return RawURL(u).Apply(request)
	})
}

func Scheme(scheme string) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		request.URL.Scheme = scheme
		return request, nil
	})
}

func User(user *neturl.Userinfo) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		request.URL.User = user
		return request, nil
	})
}

func UserPassword(username, password string) Option {
	return User(neturl.UserPassword(username, password))
}

func Host(host string) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		request.Host = host
		request.URL.Host = host
		return request, nil
	})
}

func Path(segments ...string) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		elem := []string{request.URL.Path}
		elem = append(elem, segments...)
		request.URL.Path = path.Join(elem...)

		return request, nil
	})
}

// Query applies a key/value pair to the query parameters of the *http.Request.
func Query(key, value string) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		query := request.URL.Query()
		query.Add(key, value)
		request.URL.RawQuery = query.Encode()

		return request, nil
	})
}

// Queries applies multiple key/value pairs to the query parameters of the *http.Request. It wraps url.Values.
type Queries neturl.Values

// Apply applies the Queries to the *http.Request.
func (q Queries) Apply(request *http.Request) (*http.Request, error) {
	options := make(Pipeline, 0, len(q))
	for key, values := range q {
		for _, value := range values {
			options = append(options, Query(key, value))
		}
	}

	return options.Apply(request)
}

// Header applies a key/value pair to the headers of the *http.Request.
func Header(key, value string) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		request.Header.Add(key, value)
		return request, nil
	})
}

// Headers applies multiple key/value pairs to the headers of the *http.Request. It wraps http.Header.
type Headers http.Header

// Apply applies the Headers to the *http.Request.
func (h Headers) Apply(request *http.Request) (*http.Request, error) {
	options := make(Pipeline, 0, len(h))
	for key, values := range h {
		for _, value := range values {
			options = append(options, Header(key, value))
		}
	}

	return options.Apply(request)
}

// Cookie applies a cookie to the *http.Request
func Cookie(cookie *http.Cookie) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		request.AddCookie(cookie)
		return request, nil
	})
}

// Authorization applies an "Authorization" header to the *http.Request.
func Authorization(authorization string) Option {
	return Header("Authorization", authorization)
}

// BasicAuth applies a username and password basic auth header.
func BasicAuth(username, password string) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		request.SetBasicAuth(username, password)
		return request, nil
	})
}

// BearerAuth applies an "Authorization: Bearer <token>" header to the *http.Request.
func BearerAuth(token string) Option {
	return Authorization(fmt.Sprintf("Bearer %s", token))
}

// Context applies a context.Context to the *http.Request.
func Context(ctx context.Context) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		return request.WithContext(ctx), nil
	})
}

// ContextValue applies a context key/value pair to the *http.Request.
func ContextValue(key, value interface{}) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		return request.WithContext(context.WithValue(request.Context(), key, value)), nil
	})
}

// Body applies a io.ReadCloser to the *http.Request body.
func Body(body io.ReadCloser) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		request.Body = body
		return request, nil
	})
}

func BodyReader(body io.Reader) Option {
	return Body(io.NopCloser(body))
}

// BodyBytes applies a slice of bytes to the *http.Request body.
func BodyBytes(body []byte) Option {
	return BodyReader(bytes.NewBuffer(body))
}

// BodyString applies a string to the *http.Request body.
func BodyString(body string) Option {
	return BodyBytes([]byte(body))
}

// BodyForm URL-encodes multiple key/value pairs and applies the result to the *http.Request body.
type BodyForm neturl.Values

// Apply URL-encodes the BodyForm and applies the result to the *http.Request body.
func (f BodyForm) Apply(request *http.Request) (*http.Request, error) {
	return Apply(request,
		Header("Content-Type", "application/x-www-form-urlencoded"),
		BodyString(neturl.Values(f).Encode()),
	)
}

// BodyJSON encodes an object as JSON and applies it to the *http.Request body.
func BodyJSON(v interface{}) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		body := new(bytes.Buffer)
		if err := json.NewEncoder(body).Encode(v); err != nil {
			return nil, err
		}

		return Apply(request,
			Header("Content-Type", "application/json"),
			BodyReader(body),
		)
	})
}

// BodyXML encodes an object as XML and applies it to the *http.Request body.
func BodyXML(v interface{}) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		body := new(bytes.Buffer)
		if err := xml.NewEncoder(body).Encode(v); err != nil {
			return nil, err
		}

		return Apply(request,
			Header("Content-Type", "application/xml"),
			BodyReader(body),
		)
	})
}
