package jsonapi

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalMinimumSinglePayload(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	err := MarshalPayload(writer, &Payload{
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

func TestMarshalSinglePayload(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	err := MarshalPayload(writer, &Payload{
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
	payload := &Payload{
		Links: &PayloadLinks{
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
		err := MarshalPayload(writer, payload)
		if err != nil {
			panic(err)
		}

		writer.Reset()
	}
}
