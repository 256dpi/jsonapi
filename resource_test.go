package jsonapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkWriteResource(b *testing.B) {
	res := &Resource{
		Type: "foo",
		ID:   "1",
		Attributes: Map{
			"foo": "bar",
			"bar": "baz",
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := WriteResource(httptest.NewRecorder(), http.StatusOK, res, nil)
		if err != nil {
			panic(err)
		}
	}
}
