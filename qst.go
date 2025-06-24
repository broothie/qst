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

// Do makes an *http.Request using the current DefaultClient and returns the *http.Response.
func Do(method, url string, options ...option.Option[*http.Request]) (*http.Response, error) {
	return DefaultClient.Do(method, append(option.Options[*http.Request]{URL(url)}, options...)...)
}
