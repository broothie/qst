package qst_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/broothie/qst"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		_, err := qst.New("lol what", "https://breakfast.com/api/cereals")
		assert.EqualError(t, err, `net/http: invalid method "lol what"`)
	})
}

func ExampleNew() {
	req, _ := qst.New(http.MethodPost, "http://bfast.com/api",
		qst.Scheme("https"),
		qst.Host("breakfast.com"),
		qst.Path("/cereals", "1234"),
		qst.BearerAuth("c0rnfl@k3s"),
		qst.BodyJSON(map[string]string{"name": "Honey Bunches of Oats"}),
	)

	fmt.Println(req.Method)
	fmt.Println(req.URL)
	fmt.Println(req.Header.Get("Authorization"))
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))

	// Output:
	// POST
	// https://breakfast.com/api/cereals/1234
	// Bearer c0rnfl@k3s
	// {"name":"Honey Bunches of Oats"}
}

func ExampleDo() {
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method)
		fmt.Println(r.URL)
		fmt.Println(r.Header.Get("Authorization"))
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Println(string(body))
	}))
	defer server.Close()

	qst.Do(http.MethodPost, server.URL,
		qst.Path("api", "/cereals", "1234"),
		qst.BearerAuth("c0rnfl@k3s"),
		qst.BodyJSON(map[string]string{"name": "Honey Bunches of Oats"}),
	)

	// Output:
	// POST
	// /api/cereals/1234
	// Bearer c0rnfl@k3s
	// {"name":"Honey Bunches of Oats"}
}
