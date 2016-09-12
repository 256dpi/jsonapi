package jsonapi

import (
	"bytes"
	"io"
	"io/ioutil"
)

type testRequester struct {
	method      string
	header      map[string]string
	path        string
	queryParams map[string][]string
}

func (b *testRequester) Method() string {
	return b.method
}

func (b *testRequester) Get(key string) string {
	return b.header[key]
}

func (b *testRequester) Path() string {
	return b.path
}

func (b *testRequester) QueryParams() map[string][]string {
	return b.queryParams
}

type testResponder struct {
	header map[string]string
	status int
	buffer bytes.Buffer
}

func (b *testResponder) Set(key, value string) {
	b.header[key] = value
}

func (b *testResponder) WriteHeader(status int) {
	b.status = status
}

func (b *testResponder) Write(p []byte) (int, error) {
	return b.buffer.Write(p)
}

func newTestRequester(method, path string, qp map[string][]string) *testRequester {
	if qp == nil {
		qp = make(map[string][]string)
	}

	return &testRequester{
		method:      method,
		path:        path,
		header:      make(map[string]string),
		queryParams: qp,
	}
}

func newTestResponder() *testResponder {
	return &testResponder{
		header: make(map[string]string),
	}
}

func stringReader(str string) *bytes.Reader {
	return bytes.NewReader([]byte(str))
}

func readFullString(r io.Reader) string {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	return string(buf)
}
