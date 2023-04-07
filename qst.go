package qst

import "net/http"

// New builds a new *http.Request.
func New(method, url string, options ...Option) (*http.Request, error) {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	return Apply(request, options...)
}

// Do makes an *http.Request using the current DefaultClient and returns the *http.Response.
func Do(method, url string, options ...Option) (*http.Response, error) {
	return DefaultClient.Do(method, Pipeline{URL(url)}.With(options...)...)
}
