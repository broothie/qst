package qst

import (
	"net/http"

	"github.com/broothie/option"
)

// Apply applies the Options to the *http.Request.
func Apply(request *http.Request, options ...option.Option[*http.Request]) (*http.Request, error) {
	return option.Apply(request, options...)
}
