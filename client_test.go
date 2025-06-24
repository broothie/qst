package qst_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/broothie/qst"
	"github.com/stretchr/testify/assert"
)

func TestSetClient(t *testing.T) {
	// Create a custom client with a specific timeout
	customClient := &http.Client{
		Timeout: 1 * time.Second,
	}

	// Set the custom client
	qst.SetClient(customClient)

	// Create a request that would use the custom client
	// We can't easily test the exact client being used without making an actual request,
	// but we can at least verify the function doesn't panic
	_, err := qst.New(http.MethodGet, "https://example.com")
	assert.NoError(t, err)

	// Reset to default client to avoid affecting other tests
	qst.SetClient(http.DefaultClient)
}
