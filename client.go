package qst

import "net/http"

// DefaultClient captures the current Doer.
var DefaultClient = Client{doer: http.DefaultClient}

// WithClient captures a *http.Client (or anything that implements Do(*http.Request) (*http.Response, error)).
func WithClient(doer Doer) *Client {
	return &Client{doer: doer}
}

// Doer is typically a *http.Client.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Client captures a Doer.
type Client struct {
	doer Doer
}

// Do makes a *http.Request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Do(method, url string, options ...Option) (*http.Response, error) {
	request, err := New(method, url, options...)
	if err != nil {
		return nil, err
	}

	return c.doer.Do(request)
}
