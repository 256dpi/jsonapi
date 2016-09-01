package jsonapi

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidPayload(t *testing.T) {
	readers := []io.Reader{
		stringReader(``),
		stringReader(`1`),
		stringReader(`"foo"`),
		stringReader(`true`),
		stringReader(`[]`),
	}

	for _, r := range readers {
		payload, err := UnmarshalPayload(r)
		assert.Error(t, err)
		assert.Nil(t, payload)
	}
}

func TestEmptyPayload(t *testing.T) {
	// TODO: Should return error?

	payload, err := UnmarshalPayload(stringReader(`{}`))
	assert.NoError(t, err)
	assert.Equal(t, &Payload{}, payload)
}

func TestMinimumSinglePayload(t *testing.T) {
	payload, err := UnmarshalPayload(stringReader(`{
  		"data": {
    		"type": "foo"
		}
	}`))
	assert.NoError(t, err)
	assert.Equal(t, &Payload{
		Data: &HybridResource{
			One: &Resource{
				Type: "foo",
			},
		},
	}, payload)
}

func TestSinglePayload(t *testing.T) {
	payload, err := UnmarshalPayload(stringReader(`{
  		"data": {
    		"type": "foo",
    		"id": "1",
    		"attributes": {},
    		"relationships": {}
		}
	}`))
	assert.NoError(t, err)
	assert.Equal(t, &Payload{
		Data: &HybridResource{
			One: &Resource{
				Type:          "foo",
				ID:            "1",
				Attributes:    make(Map),
				Relationships: make(map[string]HybridPayload),
			},
		},
	}, payload)
}

func BenchmarkUnmarshal(b *testing.B) {
	payload := []byte(`{
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

	reader := bytes.NewReader(payload)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := UnmarshalPayload(reader)
		if err != nil {
			panic(err)
		}

		reader.Seek(0, io.SeekStart)
	}
}
