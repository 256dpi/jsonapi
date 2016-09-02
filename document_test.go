package jsonapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBodyInvalidDocument(t *testing.T) {
	readers := []io.Reader{
		stringReader(``),
		stringReader(`1`),
		stringReader(`"foo"`),
		stringReader(`true`),
		stringReader(`[]`),
	}

	for _, r := range readers {
		doc, err := ParseBody(r)
		assert.Error(t, err)
		assert.Nil(t, doc)
	}
}

func TestParseBodyEmptyDocument(t *testing.T) {
	doc, err := ParseBody(stringReader(`{}`))
	assert.Error(t, err)
	assert.Nil(t, doc)
}

func TestParseBodyMinimumSingleDocument(t *testing.T) {
	doc, err := ParseBody(stringReader(`{
  		"data": {
    		"type": "foo"
		}
	}`))
	assert.NoError(t, err)
	assert.Equal(t, &Document{
		Data: &HybridResource{
			One: &Resource{
				Type: "foo",
			},
		},
	}, doc)
}

func TestParseBodySingleDocument(t *testing.T) {
	doc, err := ParseBody(stringReader(`{
  		"data": {
    		"type": "foo",
    		"id": "1",
    		"attributes": {},
    		"relationships": {}
		}
	}`))
	assert.NoError(t, err)
	assert.Equal(t, &Document{
		Data: &HybridResource{
			One: &Resource{
				Type:          "foo",
				ID:            "1",
				Attributes:    make(Map),
				Relationships: make(map[string]HybridDocument),
			},
		},
	}, doc)
}

func TestWriteResponseMinimumSingleDocument(t *testing.T) {
	writer := httptest.NewRecorder()
	err := WriteResponse(writer, http.StatusOK, &Document{
		Data: &HybridResource{
			One: &Resource{
				Type: "foo",
			},
		},
	})
	assert.NoError(t, err)
	assert.JSONEq(t, `{
  		"data": {
    		"type": "foo"
		}
	}`, writer.Body.String())
}

func TestWriteResponseSingleDocument(t *testing.T) {
	writer := httptest.NewRecorder()
	err := WriteResponse(writer, http.StatusOK, &Document{
		Data: &HybridResource{
			One: &Resource{
				Type: "foo",
				ID:   "1",
			},
		},
	})
	assert.NoError(t, err)
	assert.JSONEq(t, `{
  		"data": {
    		"type": "foo",
    		"id": "1"
		}
	}`, writer.Body.String())
}

func BenchmarkParseBody(b *testing.B) {
	reader := stringReader(`{
		"links": {
			"self": "http://0.0.0.0:1234/api/foo/1"
		},
		"data": {
			"type": "foo",
			"id": "1",
			"attributes": {
				"foo": "bar",
				"bar": "baz"
			}
		}
	}`)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ParseBody(reader)
		if err != nil {
			panic(err)
		}

		reader.Seek(0, io.SeekStart)
	}
}

func BenchmarkWriteResponse(b *testing.B) {
	doc := &Document{
		Links: &DocumentLinks{
			Self: "http://0.0.0.0:1234/api/foo/1",
		},
		Data: &HybridResource{
			One: &Resource{
				Type: "foo",
				ID:   "1",
				Attributes: Map{
					"foo": "bar",
					"bar": "baz",
				},
			},
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := WriteResponse(httptest.NewRecorder(), http.StatusOK, doc)
		if err != nil {
			panic(err)
		}
	}
}

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
		err := WriteResource(httptest.NewRecorder(), http.StatusOK, res)
		if err != nil {
			panic(err)
		}
	}
}
