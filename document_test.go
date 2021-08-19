package jsonapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDocumentInvalidInput(t *testing.T) {
	readers := []io.Reader{
		strings.NewReader(``),
		strings.NewReader(`1`),
		strings.NewReader(`"foo"`),
		strings.NewReader(`true`),
		strings.NewReader(`[]`),
		strings.NewReader(`{
			"data": "foo"
		}`),
		strings.NewReader(`{
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
		doc, err := ParseDocument(r)
		assert.Error(t, err)
		assert.Nil(t, doc)
	}
}

func TestParseDocumentDocumentWithErrors(t *testing.T) {
	doc, err := ParseDocument(strings.NewReader(`{
		"errors": [{
			"status": "404"
		}]
	}`))
	assert.Error(t, err)
	assert.Nil(t, doc)
	assert.Equal(t, &Error{
		Status: http.StatusNotFound,
	}, err)
}

func TestParseDocument(t *testing.T) {
	doc, err := ParseDocument(strings.NewReader(`{
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

func TestParseDocumentWithManyResources(t *testing.T) {
	doc, err := ParseDocument(strings.NewReader(`{
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

func TestParseDocumentDocumentWithRelationships(t *testing.T) {
	doc, err := ParseDocument(strings.NewReader(`{
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

func TestParseDocumentWithBigNumbers(t *testing.T) {
	type test struct {
		Num int64 `json:"num"`
	}

	doc, err := ParseDocument(strings.NewReader(`{
  		"data": {
    		"type": "foo",
    		"id": "1",
    		"attributes": {
    			"num": 4699539
    		},
    		"relationships": {}
		}
	}`))
	assert.NoError(t, err)
	assert.Equal(t, &Document{
		Data: &HybridResource{
			One: &Resource{
				Type: "foo",
				ID:   "1",
				Attributes: map[string]interface{}{
					"num": json.Number("4699539"),
				},
				Relationships: map[string]*Document{},
			},
		},
	}, doc)

	m := &test{}
	err = doc.Data.One.Attributes.Assign(m)
	assert.NoError(t, err)
}

func TestParseDocumentNullLink(t *testing.T) {
	doc, err := ParseDocument(strings.NewReader(`{
  		"links": {
			"self": null
		}
	}`))
	assert.NoError(t, err)
	assert.Equal(t, &Document{
		Links: &DocumentLinks{
			Self: NullLink,
		},
	}, doc)
}

func TestWriteResponseOneResource(t *testing.T) {
	res := httptest.NewRecorder()

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
	}`, res.Body.String())
}

func TestWriteResponseManyResources(t *testing.T) {
	res := httptest.NewRecorder()

	err := WriteResponse(res, http.StatusOK, &Document{
		Data: &HybridResource{
			Many: []*Resource{
				{
					Type: "foo",
					ID:   "1",
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.JSONEq(t, `{
  		"data": [
			{
    			"type": "foo",
    			"id": "1"
			}
		]
	}`, res.Body.String())
}

func TestWriteResponseNoResources(t *testing.T) {
	res := httptest.NewRecorder()

	err := WriteResponse(res, http.StatusOK, &Document{
		Data: &HybridResource{},
	})
	assert.NoError(t, err)
	assert.JSONEq(t, `{
  		"data": null
	}`, res.Body.String())
}

func TestWriteResponseNullLink(t *testing.T) {
	res := httptest.NewRecorder()

	err := WriteResponse(res, http.StatusOK, &Document{
		Links: &DocumentLinks{
			Self: NullLink,
		},
	})
	assert.NoError(t, err)
	assert.JSONEq(t, `{
  		"links": {
			"self": null
		}
	}`, res.Body.String())
}

func BenchmarkParseDocument(b *testing.B) {
	reader := strings.NewReader(`{
		"links": {
			"self": "/api/foo/1"
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
		_, err := ParseDocument(reader)
		if err != nil {
			panic(err)
		}

		_, err = reader.Seek(0, io.SeekStart)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkWriteResponse(b *testing.B) {
	doc := &Document{
		Links: &DocumentLinks{
			Self: "/api/foo/1",
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
		res := httptest.NewRecorder()

		err := WriteResponse(res, http.StatusOK, doc)
		if err != nil {
			panic(err)
		}
	}
}
