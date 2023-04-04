package qst

import "net/http"

// NewGet builds a new *http.Request with method GET.
func NewGet(options ...Option) (*http.Request, error) {
    return New(http.MethodGet, options...)
}

// Get makes a GET request using the current DefaultClient and returns the *http.Response.
func Get(options ...Option) (*http.Response, error) {
    return Do(http.MethodGet, options...)
}

// Get makes a GET request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Get(options ...Option) (*http.Response, error) {
    return c.Do(http.MethodGet, options...)
}

// NewHead builds a new *http.Request with method HEAD.
func NewHead(options ...Option) (*http.Request, error) {
    return New(http.MethodHead, options...)
}

// Head makes a HEAD request using the current DefaultClient and returns the *http.Response.
func Head(options ...Option) (*http.Response, error) {
    return Do(http.MethodHead, options...)
}

// Head makes a HEAD request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Head(options ...Option) (*http.Response, error) {
    return c.Do(http.MethodHead, options...)
}

// NewPost builds a new *http.Request with method POST.
func NewPost(options ...Option) (*http.Request, error) {
    return New(http.MethodPost, options...)
}

// Post makes a POST request using the current DefaultClient and returns the *http.Response.
func Post(options ...Option) (*http.Response, error) {
    return Do(http.MethodPost, options...)
}

// Post makes a POST request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Post(options ...Option) (*http.Response, error) {
    return c.Do(http.MethodPost, options...)
}

// NewPut builds a new *http.Request with method PUT.
func NewPut(options ...Option) (*http.Request, error) {
    return New(http.MethodPut, options...)
}

// Put makes a PUT request using the current DefaultClient and returns the *http.Response.
func Put(options ...Option) (*http.Response, error) {
    return Do(http.MethodPut, options...)
}

// Put makes a PUT request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Put(options ...Option) (*http.Response, error) {
    return c.Do(http.MethodPut, options...)
}

// NewPatch builds a new *http.Request with method PATCH.
func NewPatch(options ...Option) (*http.Request, error) {
    return New(http.MethodPatch, options...)
}

// Patch makes a PATCH request using the current DefaultClient and returns the *http.Response.
func Patch(options ...Option) (*http.Response, error) {
    return Do(http.MethodPatch, options...)
}

// Patch makes a PATCH request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Patch(options ...Option) (*http.Response, error) {
    return c.Do(http.MethodPatch, options...)
}

// NewDelete builds a new *http.Request with method DELETE.
func NewDelete(options ...Option) (*http.Request, error) {
    return New(http.MethodDelete, options...)
}

// Delete makes a DELETE request using the current DefaultClient and returns the *http.Response.
func Delete(options ...Option) (*http.Response, error) {
    return Do(http.MethodDelete, options...)
}

// Delete makes a DELETE request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Delete(options ...Option) (*http.Response, error) {
    return c.Do(http.MethodDelete, options...)
}

// NewConnect builds a new *http.Request with method CONNECT.
func NewConnect(options ...Option) (*http.Request, error) {
    return New(http.MethodConnect, options...)
}

// Connect makes a CONNECT request using the current DefaultClient and returns the *http.Response.
func Connect(options ...Option) (*http.Response, error) {
    return Do(http.MethodConnect, options...)
}

// Connect makes a CONNECT request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Connect(options ...Option) (*http.Response, error) {
    return c.Do(http.MethodConnect, options...)
}

// NewOptions builds a new *http.Request with method OPTIONS.
func NewOptions(options ...Option) (*http.Request, error) {
    return New(http.MethodOptions, options...)
}

// Options makes a OPTIONS request using the current DefaultClient and returns the *http.Response.
func Options(options ...Option) (*http.Response, error) {
    return Do(http.MethodOptions, options...)
}

// Options makes a OPTIONS request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Options(options ...Option) (*http.Response, error) {
    return c.Do(http.MethodOptions, options...)
}

// NewTrace builds a new *http.Request with method TRACE.
func NewTrace(options ...Option) (*http.Request, error) {
    return New(http.MethodTrace, options...)
}

// Trace makes a TRACE request using the current DefaultClient and returns the *http.Response.
func Trace(options ...Option) (*http.Response, error) {
    return Do(http.MethodTrace, options...)
}

// Trace makes a TRACE request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Trace(options ...Option) (*http.Response, error) {
    return c.Do(http.MethodTrace, options...)
}
