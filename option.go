package qst

import (
	"net/http"

	"github.com/broothie/option"
)

// Option is an alias for the generic option interface.
type Option = option.Option[*http.Request]

// Pipeline is a collection of options, which can be applied as a whole.
type Pipeline = option.Pipeline[*http.Request]

// OptionFunc is a function form of Option.
type OptionFunc = option.OptionFunc[*http.Request]

// Apply applies the Options to the *http.Request.
func Apply(request *http.Request, options ...Option) (*http.Request, error) {
	return option.Apply(request, options...)
}
