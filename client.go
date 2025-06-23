package qst

import (
	"net/http"

	"github.com/broothie/option"
)

// DefaultClient captures the current Doer.
var DefaultClient = NewClient(http.DefaultClient)

// Doer is typically an *http.Client.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// NewClient creates a new Client.
func NewClient(doer Doer, options ...option.Option[*http.Request]) *Client {
	return &Client{doer: doer, options: options}
}

// Client captures a Doer and Options to apply to every request.
type Client struct {
	doer    Doer
	options []option.Option[*http.Request]
}

// Do makes an *http.Request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Do(method string, options ...option.Option[*http.Request]) (*http.Response, error) {
	opts := append(c.options, options...)
	request, err := New(method, "", opts...)
	if err != nil {
		return nil, err
	}

	return c.doer.Do(request)
}
