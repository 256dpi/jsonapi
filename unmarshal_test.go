package jsonapi

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidDocument(t *testing.T) {
	readers := []io.Reader{
		stringReader(``),
		stringReader(`1`),
		stringReader(`"foo"`),
		stringReader(`true`),
		stringReader(`[]`),
	}

	for _, r := range readers {
		doc, err := UnmarshalDocument(r)
		assert.Error(t, err)
		assert.Nil(t, doc)
	}
}

func TestEmptyDocument(t *testing.T) {
	// TODO: Should return error?

	doc, err := UnmarshalDocument(stringReader(`{}`))
	assert.NoError(t, err)
	assert.Equal(t, &Document{}, doc)
}

func TestMinimumSingleDocument(t *testing.T) {
	doc, err := UnmarshalDocument(stringReader(`{
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

func TestSingleDocument(t *testing.T) {
	doc, err := UnmarshalDocument(stringReader(`{
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

func BenchmarkUnmarshal(b *testing.B) {
	doc := []byte(`{
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

	reader := bytes.NewReader(doc)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := UnmarshalDocument(reader)
		if err != nil {
			panic(err)
		}

		reader.Seek(0, io.SeekStart)
	}
}
