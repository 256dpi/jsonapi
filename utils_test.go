package jsonapi

import (
	"bytes"
	"net/http"
	"net/url"
)

func constructRequest(method, path string) *http.Request {
	u, err := url.Parse(path)
	if err != nil {
		panic(err)
	}

	r := &http.Request{
		Method: method,
		URL:    u,
		Header: make(http.Header),
	}

	r.Header.Set("Accept", MediaType)

	return r
}

func stringReader(str string) *bytes.Reader {
	return bytes.NewReader([]byte(str))
}
