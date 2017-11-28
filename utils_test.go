package jsonapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
)

func newTestRequest(method, path string) *http.Request {
	r, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}

	return r
}

func newTestResponseRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

func stringReader(str string) *bytes.Reader {
	return bytes.NewReader([]byte(str))
}
