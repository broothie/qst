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
			assert.Equal(t, "Bearer asdf", r.Header.Get("Authorization"))
			assert.Equal(t, "here", r.Header.Get("something"))
		}))
		defer server.Close()

		// Exercise
		_, err := qst.Post(
			qst.URL(server.URL),
			qst.BearerAuth("asdf"),
			qst.Header("something", "here"),
		)

		require.NoError(t, err)
	})

	t.Run("post", func(t *testing.T) {
		// Set up
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			// Verify
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "Bearer asdf", r.Header.Get("Authorization"))

			var body map[string]interface{}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))

			assert.Equal(t, "here", body["something"])
			assert.EqualValues(t, 1, body["a"])
		}))
		defer server.Close()

		// Exercise
		_, err := qst.Post(
			qst.URL(server.URL),
			qst.BearerAuth("asdf"),
			qst.BodyJSON(map[string]interface{}{"something": "here", "a": 1}),
		)
		require.NoError(t, err)
	})
}
