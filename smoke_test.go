package qst_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/broothie/qst"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSmoke(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		// Set up
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			// Verify
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "Bearer c0rnfl@k3s", r.Header.Get("Authorization"))
			assert.Equal(t, "oats", r.Header.Get("grain"))
		}))
		defer server.Close()

		// Exercise
		_, err := qst.Post(server.URL,
			qst.WithBearerAuth("c0rnfl@k3s"),
			qst.WithHeader("grain", "oats"),
		)

		require.NoError(t, err)
	})

	t.Run("post", func(t *testing.T) {
		// Set up
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			// Verify
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "Bearer c0rnfl@k3s", r.Header.Get("Authorization"))

			var body map[string]interface{}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))

			assert.Equal(t, "Raisin Bran Crunch", body["name"])
			assert.True(t, body["raisins"].(bool))
		}))
		defer server.Close()

		// Exercise
		_, err := qst.Post(server.URL,
			qst.WithBearerAuth("c0rnfl@k3s"),
			qst.WithBodyJSON(map[string]interface{}{"name": "Raisin Bran Crunch", "raisins": true}),
		)
		require.NoError(t, err)
	})
}
