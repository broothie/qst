package qst

import "net/http"

// client is the global HTTP client used by the Do function.
var client = http.DefaultClient

// SetClient sets the global HTTP client used by the Do function.
func SetClient(c *http.Client) {
	client = c
}
