package jsonapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteResourceEmpty(t *testing.T) {
	// TODO: Should this raise an error?

	rec := httptest.NewRecorder()
	err := WriteResource(rec, http.StatusOK, &Resource{}, nil)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"data": {
			"type": ""
		}
	}`, rec.Body.String())
}

func TestWriteResourcesEmpty(t *testing.T) {
	// TODO: Should this raise an error?

	rec := httptest.NewRecorder()
	err := WriteResources(rec, http.StatusOK, []*Resource{{}}, nil)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"data": [
			{
				"type": ""
			}
		]
	}`, rec.Body.String())
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
		err := WriteResource(httptest.NewRecorder(), http.StatusOK, res, nil)
		if err != nil {
			panic(err)
		}
	}
}
