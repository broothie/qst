package qst_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/broothie/option"
	"github.com/broothie/qst"
)

func ExampleOptionFunc() {
	token := "c0rnfl@k3s"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header.Get("Authorization"))
	}))
	defer server.Close()

	client := qst.NewClient(server.Client(),
		qst.URL(server.URL),
		option.Func[*http.Request](func(request *http.Request) (*http.Request, error) {
			return qst.BearerAuth(token).Apply(request)
		}),
	)

	client.Get()
	token = "c00ki3cr!5p"
	client.Get()

	// Output: Bearer c0rnfl@k3s
	// Bearer c00ki3cr!5p
}
