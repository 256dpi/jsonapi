package jsonapi

import (
	"bytes"
	"net/http"
	"net/url"

	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/standard"
)

func constructHTTPRequest(method, path string) *http.Request {
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

func constructRequest(method, path string) engine.Request {
	return standard.NewRequest(constructHTTPRequest(method, path), nil)
}

func stringReader(str string) *bytes.Reader {
	return bytes.NewReader([]byte(str))
}
