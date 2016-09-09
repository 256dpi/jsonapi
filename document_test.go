package jsonapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/engine/standard"
	"github.com/stretchr/testify/assert"
)

func TestParseBodyInvalidDocument(t *testing.T) {
	readers := []io.Reader{
		stringReader(``),
		stringReader(`1`),
		stringReader(`"foo"`),
		stringReader(`true`),
		stringReader(`[]`),
		stringReader(`{
			"data": "foo"
		}`),
		stringReader(`{
			"data": {
				"type": "foo",
				"id": "1",
				"relationships": {
					"bar": "foo"
				}
			}
		}`),
	}

	for _, r := range readers {
		doc, err := ParseBody(r)
		assert.Error(t, err)
		assert.Nil(t, doc)
	}
}

func TestParseBodyDocumentWithErrors(t *testing.T) {
	doc, err := ParseBody(stringReader(`{
		"errors": [{
			"status": "404"
		}]
	}`))
	assert.Error(t, err)
	assert.Nil(t, doc)
}

func TestParseBodyMinimumDocument(t *testing.T) {
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

func TestParseBodyDocument(t *testing.T) {
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
				Attributes:    make(map[string]interface{}),
				Relationships: make(map[string]*Document),
			},
		},
	}, doc)
}

func TestParseBodyDocuments(t *testing.T) {
	doc, err := ParseBody(stringReader(`{
  		"data": [
  			{
				"type": "foo",
				"id": "1",
				"attributes": {},
				"relationships": {}
			}
		]
	}`))
	assert.NoError(t, err)
	assert.Equal(t, &Document{
		Data: &HybridResource{
			Many: []*Resource{
				{
					Type:          "foo",
					ID:            "1",
					Attributes:    make(map[string]interface{}),
					Relationships: make(map[string]*Document),
				},
			},
		},
	}, doc)
}

func TestParseBodyDocumentWithRelationship(t *testing.T) {
	doc, err := ParseBody(stringReader(`{
  		"data": {
    		"type": "foo",
    		"id": "1",
    		"relationships": {
    			"bar": {
    				"data": {
    					"type": "bar"
    				}
				}
    		}
		}
	}`))
	assert.NoError(t, err)
	assert.Equal(t, &Document{
		Data: &HybridResource{
			One: &Resource{
				Type: "foo",
				ID:   "1",
				Relationships: map[string]*Document{
					"bar": {
						Data: &HybridResource{
							One: &Resource{
								Type: "bar",
							},
						},
					},
				},
			},
		},
	}, doc)
}

func TestWriteResponseMinimumSingleDocument(t *testing.T) {
	res, rec := constructResponseAndRecorder()

	err := WriteResponse(res, http.StatusOK, &Document{
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
	}`, rec.Body.String())
}

func TestWriteResponseSingleDocument(t *testing.T) {
	res, rec := constructResponseAndRecorder()

	err := WriteResponse(res, http.StatusOK, &Document{
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
	}`, rec.Body.String())
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
		res := standard.NewResponse(httptest.NewRecorder(), nil)

		err := WriteResponse(res, http.StatusOK, doc)
		if err != nil {
			panic(err)
		}
	}
}
