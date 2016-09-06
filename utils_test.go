package jsonapi

import (
	"bytes"
	"net/http"
	"net/url"

	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/standard"
)

func constructRequest(method, path string) engine.Request {
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

	return standard.NewRequest(r, nil)
}

func stringReader(str string) *bytes.Reader {
	return bytes.NewReader([]byte(str))
}
