package jsonapi

import (
	"bytes"
	"io"
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

	r.Header.Set("Accept", ContentType)

	return r
}

func stringReader(str string) io.Reader {
	return bytes.NewBufferString(str)
}
