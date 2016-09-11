package jsonapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBridgeRequest(t *testing.T) {
	r, err := http.NewRequest("GET", "http://localhost/foo?foo=bar", nil)
	assert.NoError(t, err)

	r.Header.Set("foo", "bar")

	b := BridgeRequest(r)
	assert.Equal(t, "GET", b.Method())
	assert.Equal(t, "/foo", b.Path())
	assert.Equal(t, "bar", b.Get("foo"))
	assert.Equal(t, []string{"bar"}, b.QueryParams()["foo"])
}

func TestBridgeResponseWriter(t *testing.T) {
	rec := httptest.NewRecorder()

	b := BridgeResponseWriter(rec)
	b.Set("foo", "bar")
	b.WriteHeader(100)
	b.Write([]byte("foo"))

	assert.Equal(t, "bar", rec.Header().Get("foo"))
	assert.Equal(t, 100, rec.Result().StatusCode)
	assert.Equal(t, "foo", rec.Body.String())
}
