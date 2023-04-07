package qst_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/broothie/qst"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMethods(t *testing.T) {
	type testCase struct {
		new      func(url string, options ...qst.Option) (*http.Request, error)
		do       func(url string, options ...qst.Option) (*http.Response, error)
		clientDo func(options ...qst.Option) (*http.Response, error)
	}

	client := http.DefaultClient
	var methods = map[string]testCase{
		http.MethodGet: {
			new:      qst.NewGet,
			do:       qst.Get,
			clientDo: qst.NewClient(client).Get,
		},
		http.MethodHead: {
			new:      qst.NewHead,
			do:       qst.Head,
			clientDo: qst.NewClient(client).Head,
		},
		http.MethodPost: {
			new:      qst.NewPost,
			do:       qst.Post,
			clientDo: qst.NewClient(client).Post,
		},
		http.MethodPut: {
			new:      qst.NewPut,
			do:       qst.Put,
			clientDo: qst.NewClient(client).Put,
		},
		http.MethodPatch: {
			new:      qst.NewPatch,
			do:       qst.Patch,
			clientDo: qst.NewClient(client).Patch,
		},
		http.MethodDelete: {
			new:      qst.NewDelete,
			do:       qst.Delete,
			clientDo: qst.NewClient(client).Delete,
		},
		http.MethodConnect: {
			new:      qst.NewConnect,
			do:       qst.Connect,
			clientDo: qst.NewClient(client).Connect,
		},
		http.MethodOptions: {
			new:      qst.NewOptions,
			do:       qst.Options,
			clientDo: qst.NewClient(client).Options,
		},
		http.MethodTrace: {
			new:      qst.NewTrace,
			do:       qst.Trace,
			clientDo: qst.NewClient(client).Trace,
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
			_, err := tc.clientDo(qst.URL(server.URL))
			require.NoError(t, err)
		})

		server.Close()
	}
}
