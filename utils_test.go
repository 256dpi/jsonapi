package jsonapi

import (
	"net/http"
	"net/url"
)

func constructRequest(path string) *http.Request {
	u, err := url.Parse(path)
	if err != nil {
		panic(err)
	}

	r := &http.Request{
		URL:    u,
		Header: make(http.Header),
	}

	r.Header.Set("Content-Type", ContentType)

	return r
}
