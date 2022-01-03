package qst

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Do(t *testing.T) {
	_, err := WithClient(http.DefaultClient).Do("lol what", "")
	assert.EqualError(t, err, `net/http: invalid method "lol what"`)
}
