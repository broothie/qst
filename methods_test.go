package qst_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/broothie/option"
	"github.com/broothie/qst"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMethods(t *testing.T) {
	type testCase struct {
		new func(url string, options ...option.Option[*http.Request]) (*http.Request, error)
		do  func(url string, options ...option.Option[*http.Request]) (*http.Response, error)
	}

	var methods = map[string]testCase{
		http.MethodGet: {
			new: qst.NewGet,
			do:  qst.Get,
		},
		http.MethodHead: {
			new: qst.NewHead,
			do:  qst.Head,
		},
		http.MethodPost: {
			new: qst.NewPost,
			do:  qst.Post,
		},
		http.MethodPut: {
			new: qst.NewPut,
			do:  qst.Put,
		},
		http.MethodPatch: {
			new: qst.NewPatch,
			do:  qst.Patch,
		},
		http.MethodDelete: {
			new: qst.NewDelete,
			do:  qst.Delete,
		},
		http.MethodConnect: {
			new: qst.NewConnect,
			do:  qst.Connect,
		},
		http.MethodOptions: {
			new: qst.NewOptions,
			do:  qst.Options,
		},
		http.MethodTrace: {
			new: qst.NewTrace,
			do:  qst.Trace,
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

		server.Close()
	}
}
