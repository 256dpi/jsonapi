package jsonapi

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"net/http/httptest"
)

func TestMarshalMinimumSingleDocument(t *testing.T) {
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

func TestMarshalSingleDocument(t *testing.T) {
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

func BenchmarkMarshal(b *testing.B) {
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
