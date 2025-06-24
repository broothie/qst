package qst

import (
	"net/http"

	"github.com/broothie/option"
)

// New builds a new *http.Request.
func New(method, url string, options ...option.Option[*http.Request]) (*http.Request, error) {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	return option.Apply(request, options...)
}

// Do makes an *http.Request using the global client and returns the *http.Response.
func Do(method, url string, options ...option.Option[*http.Request]) (*http.Response, error) {
	request, err := New(method, url, options...)
	if err != nil {
		return nil, err
	}

	return client.Do(request)
}
