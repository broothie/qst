package qst

import "net/http"

// NewGet builds a new *http.Request with method GET.
func NewGet(url string, options ...Option) (*http.Request, error) {
    return New(http.MethodGet, url, options...)
}

// Get makes a GET request using the current DefaultClient and returns the *http.Response.
func Get(url string, options ...Option) (*http.Response, error) {
    return Do(http.MethodGet, url, options...)
}

// Get makes a GET request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Get(url string, options ...Option) (*http.Response, error) {
    return c.Do(http.MethodGet, url, options...)
}

// NewHead builds a new *http.Request with method HEAD.
func NewHead(url string, options ...Option) (*http.Request, error) {
    return New(http.MethodHead, url, options...)
}

// Head makes a HEAD request using the current DefaultClient and returns the *http.Response.
func Head(url string, options ...Option) (*http.Response, error) {
    return Do(http.MethodHead, url, options...)
}

// Head makes a HEAD request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Head(url string, options ...Option) (*http.Response, error) {
    return c.Do(http.MethodHead, url, options...)
}

// NewPost builds a new *http.Request with method POST.
func NewPost(url string, options ...Option) (*http.Request, error) {
    return New(http.MethodPost, url, options...)
}

// Post makes a POST request using the current DefaultClient and returns the *http.Response.
func Post(url string, options ...Option) (*http.Response, error) {
    return Do(http.MethodPost, url, options...)
}

// Post makes a POST request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Post(url string, options ...Option) (*http.Response, error) {
    return c.Do(http.MethodPost, url, options...)
}

// NewPut builds a new *http.Request with method PUT.
func NewPut(url string, options ...Option) (*http.Request, error) {
    return New(http.MethodPut, url, options...)
}

// Put makes a PUT request using the current DefaultClient and returns the *http.Response.
func Put(url string, options ...Option) (*http.Response, error) {
    return Do(http.MethodPut, url, options...)
}

// Put makes a PUT request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Put(url string, options ...Option) (*http.Response, error) {
    return c.Do(http.MethodPut, url, options...)
}

// NewPatch builds a new *http.Request with method PATCH.
func NewPatch(url string, options ...Option) (*http.Request, error) {
    return New(http.MethodPatch, url, options...)
}

// Patch makes a PATCH request using the current DefaultClient and returns the *http.Response.
func Patch(url string, options ...Option) (*http.Response, error) {
    return Do(http.MethodPatch, url, options...)
}

// Patch makes a PATCH request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Patch(url string, options ...Option) (*http.Response, error) {
    return c.Do(http.MethodPatch, url, options...)
}

// NewDelete builds a new *http.Request with method DELETE.
func NewDelete(url string, options ...Option) (*http.Request, error) {
    return New(http.MethodDelete, url, options...)
}

// Delete makes a DELETE request using the current DefaultClient and returns the *http.Response.
func Delete(url string, options ...Option) (*http.Response, error) {
    return Do(http.MethodDelete, url, options...)
}

// Delete makes a DELETE request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Delete(url string, options ...Option) (*http.Response, error) {
    return c.Do(http.MethodDelete, url, options...)
}

// NewConnect builds a new *http.Request with method CONNECT.
func NewConnect(url string, options ...Option) (*http.Request, error) {
    return New(http.MethodConnect, url, options...)
}

// Connect makes a CONNECT request using the current DefaultClient and returns the *http.Response.
func Connect(url string, options ...Option) (*http.Response, error) {
    return Do(http.MethodConnect, url, options...)
}

// Connect makes a CONNECT request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Connect(url string, options ...Option) (*http.Response, error) {
    return c.Do(http.MethodConnect, url, options...)
}

// NewOptions builds a new *http.Request with method OPTIONS.
func NewOptions(url string, options ...Option) (*http.Request, error) {
    return New(http.MethodOptions, url, options...)
}

// Options makes a OPTIONS request using the current DefaultClient and returns the *http.Response.
func Options(url string, options ...Option) (*http.Response, error) {
    return Do(http.MethodOptions, url, options...)
}

// Options makes a OPTIONS request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Options(url string, options ...Option) (*http.Response, error) {
    return c.Do(http.MethodOptions, url, options...)
}

// NewTrace builds a new *http.Request with method TRACE.
func NewTrace(url string, options ...Option) (*http.Request, error) {
    return New(http.MethodTrace, url, options...)
}

// Trace makes a TRACE request using the current DefaultClient and returns the *http.Response.
func Trace(url string, options ...Option) (*http.Response, error) {
    return Do(http.MethodTrace, url, options...)
}

// Trace makes a TRACE request and returns the *http.Response using the Doer assigned to c.
func (c *Client) Trace(url string, options ...Option) (*http.Response, error) {
    return c.Do(http.MethodTrace, url, options...)
}
