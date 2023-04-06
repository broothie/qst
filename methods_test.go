package qst

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMethods(t *testing.T) {
	type testCase struct {
		new      func(url string, options ...Option) (*http.Request, error)
		do       func(url string, options ...Option) (*http.Response, error)
		clientDo func(options ...Option) (*http.Response, error)
	}

	client := http.DefaultClient
	var methods = map[string]testCase{
		http.MethodGet: {
			new:      NewGet,
			do:       Get,
			clientDo: NewClient(client).Get,
		},
		http.MethodHead: {
			new:      NewHead,
			do:       Head,
			clientDo: NewClient(client).Head,
		},
		http.MethodPost: {
			new:      NewPost,
			do:       Post,
			clientDo: NewClient(client).Post,
		},
		http.MethodPut: {
			new:      NewPut,
			do:       Put,
			clientDo: NewClient(client).Put,
		},
		http.MethodPatch: {
			new:      NewPatch,
			do:       Patch,
			clientDo: NewClient(client).Patch,
		},
		http.MethodDelete: {
			new:      NewDelete,
			do:       Delete,
			clientDo: NewClient(client).Delete,
		},
		http.MethodConnect: {
			new:      NewConnect,
			do:       Connect,
			clientDo: NewClient(client).Connect,
		},
		http.MethodOptions: {
			new:      NewOptions,
			do:       Options,
			clientDo: NewClient(client).Options,
		},
		http.MethodTrace: {
			new:      NewTrace,
			do:       Trace,
			clientDo: NewClient(client).Trace,
		},
	}

	for method, tc := range methods {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			assert.Equal(t, method, r.Method)
		}))

		t.Run("new", func(t *testing.T) {
			req, err := tc.new(server.URL)
			assert.NoError(t, err)

			_, err = http.DefaultClient.Do(req)
			require.NoError(t, err)
		})

		t.Run("do", func(t *testing.T) {
			_, err := tc.do(server.URL)
			require.NoError(t, err)
		})

		t.Run("clientDo", func(t *testing.T) {
			_, err := tc.clientDo(URL(server.URL))
			require.NoError(t, err)
		})

		server.Close()
	}
}
