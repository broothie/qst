package qst

import (
	"net/http"
)

// Option is an option for building an *http.Request.
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

func (p Pipeline) With(options ...Option) Pipeline {
	return append(p, options...)
}

// OptionFunc is a function form of Option.
type OptionFunc func(request *http.Request) (*http.Request, error)

// Apply applies the OptionFunc to the *http.Request.
func (f OptionFunc) Apply(request *http.Request) (*http.Request, error) { return f(request) }

// Apply applies the Options to the *http.Request.
func Apply(request *http.Request, options ...Option) (*http.Request, error) {
	return Pipeline(options).Apply(request)
}
