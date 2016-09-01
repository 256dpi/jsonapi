package jsonapi

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalMinimumSingleDocument(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	err := MarshalDocument(writer, &Document{
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
	}`, writer.String())
}

func TestMarshalSingleDocument(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	err := MarshalDocument(writer, &Document{
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
	}`, writer.String())
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

	writer := bytes.NewBuffer(nil)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := MarshalDocument(writer, doc)
		if err != nil {
			panic(err)
		}

		writer.Reset()
	}
}
