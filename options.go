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
	"net/url"
	"strings"
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

// QueryValue applies a key/value pair to the query parameters of the *http.Request.
func QueryValue(key, value string) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		query := request.URL.Query()
		query.Add(key, value)
		request.URL.RawQuery = query.Encode()

		return request, nil
	})
}

// Query applies multiple key/value pairs to the query parameters of the *http.Request. It wraps url.Values.
type Query url.Values

// Apply applies the Query to the *http.Request.
func (q Query) Apply(request *http.Request) (*http.Request, error) {
	options := make(Pipeline, 0, len(q))
	for key, values := range q {
		for _, value := range values {
			options = append(options, QueryValue(key, value))
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

// Body applies a io.ReadCloser to the *http.Request body.
func Body(body io.ReadCloser) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		request.Body = body
		return request, nil
	})
}

// BodyBytes applies a slice of bytes to the *http.Request body.
func BodyBytes(body []byte) Option {
	return Body(ioutil.NopCloser(bytes.NewBuffer(body)))
}

// BodyString applies a string to the *http.Request body.
func BodyString(body string) Option {
	return Body(ioutil.NopCloser(strings.NewReader(body)))
}

// BodyForm URL-encodes multiple key/value pairs and applies the result to the *http.Request body.
type BodyForm url.Values

// Apply URL-encodes the BodyForm and applies the result to the *http.Request body.
func (f BodyForm) Apply(request *http.Request) (*http.Request, error) {
	return Pipeline{
		Header("Content-Type", "application/x-www-form-urlencoded"),
		BodyString(url.Values(f).Encode()),
	}.Apply(request)
}

// BodyJSON encodes an object as JSON and applies it to the *http.Request body.
func BodyJSON(v interface{}) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		body := new(bytes.Buffer)
		if err := json.NewEncoder(body).Encode(v); err != nil {
			return nil, err
		}

		return Pipeline{
			Header("Content-Type", "application/json"),
			Body(ioutil.NopCloser(body)),
		}.Apply(request)
	})
}

// BodyXML encodes an object as XML and applies it to the *http.Request body.
func BodyXML(v interface{}) Option {
	return OptionFunc(func(request *http.Request) (*http.Request, error) {
		body := new(bytes.Buffer)
		if err := xml.NewEncoder(body).Encode(v); err != nil {
			return nil, err
		}

		return Pipeline{
			Header("Content-Type", "application/xml"),
			Body(ioutil.NopCloser(body)),
		}.Apply(request)
	})
}
