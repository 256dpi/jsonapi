package jsonapi

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRequestError(t *testing.T) {
	invalidAccept := constructRequest("GET", "")
	invalidAccept.Header.Set("Accept", "foo")

	invalidContentType := constructRequest("GET", "")
	invalidContentType.Header.Set("Content-Type", "foo")

	list := []*http.Request{
		invalidAccept,
		invalidContentType,
		constructRequest("PUT", ""),
		constructRequest("GET", ""),
		constructRequest("POST", ""),
		constructRequest("GET", "/"),
		constructRequest("GET", "foo/bar/baz/qux"),
		constructRequest("GET", "foo/bar/baz/qux/quux"),
		constructRequest("GET", "foo?page[number]=bar"),
		constructRequest("GET", "foo?page[size]=bar"),
		constructRequest("GET", "foo?page[number]=1"),
		constructRequest("GET", "foo?page[size]=1"),
		constructRequest("GET", "foo?page[number]=bar&page[number]=baz"),
		constructRequest("GET", "foo?page[size]=bar&page[size]=baz"),
	}

	for _, r := range list {
		req, err := ParseRequest(r, "")
		assert.Error(t, err)
		assert.Nil(t, req)
	}
}

func TestParseRequestMethodOverride(t *testing.T) {
	r := constructRequest("GET", "foo")
	r.Header.Set("Content-Type", ContentType)
	r.Header.Set("X-HTTP-Method-Override", "PATCH")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Action:       Update,
		Target:       ResourceCollection,
		ResourceType: "foo",
	}, req)
}

func TestParseRequestPrefix(t *testing.T) {
	r := constructRequest("GET", "foo/bar")

	req, err := ParseRequest(r, "foo/")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Action:       Fetch,
		Target:       ResourceCollection,
		ResourceType: "bar",
	}, req)
}

func TestParseRequestResource(t *testing.T) {
	r := constructRequest("GET", "foo")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Action:       Fetch,
		Target:       ResourceCollection,
		ResourceType: "foo",
	}, req)
}

func TestParseRequestResourceID(t *testing.T) {
	r := constructRequest("GET", "foo/1")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Action:       Fetch,
		Target:       SingleResource,
		ResourceType: "foo",
		ResourceID:   "1",
	}, req)
}

func TestParseRequestRelatedResource(t *testing.T) {
	r := constructRequest("GET", "foo/1/bar")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Action:          Fetch,
		Target:          RelatedResource,
		ResourceType:    "foo",
		ResourceID:      "1",
		RelatedResource: "bar",
	}, req)
}

func TestParseRequestRelationship(t *testing.T) {
	r := constructRequest("GET", "foo/1/relationships/bar")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Action:       Fetch,
		Target:       Relationship,
		ResourceType: "foo",
		ResourceID:   "1",
		Relationship: "bar",
	}, req)
}

func TestParseRequestInclude(t *testing.T) {
	r1 := constructRequest("GET", "foo?include=bar,baz")

	req, err := ParseRequest(r1, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Action:       Fetch,
		Target:       ResourceCollection,
		ResourceType: "foo",
		Include:      []string{"bar", "baz"},
	}, req)

	r2 := constructRequest("GET", "foo?include=bar&include=baz,qux")

	req, err = ParseRequest(r2, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Action:       Fetch,
		Target:       ResourceCollection,
		ResourceType: "foo",
		Include:      []string{"bar", "baz", "qux"},
	}, req)
}

func TestParseRequestSorting(t *testing.T) {
	r1 := constructRequest("GET", "foo?sort=bar,baz")

	req, err := ParseRequest(r1, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Action:       Fetch,
		Target:       ResourceCollection,
		ResourceType: "foo",
		Sorting:      []string{"bar", "baz"},
	}, req)

	r2 := constructRequest("GET", "foo?sort=bar&sort=baz,qux")

	req, err = ParseRequest(r2, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Action:       Fetch,
		Target:       ResourceCollection,
		ResourceType: "foo",
		Sorting:      []string{"bar", "baz", "qux"},
	}, req)
}

func TestParseRequestPage(t *testing.T) {
	r1 := constructRequest("GET", "foo?page[number]=1&page[size]=2")

	req, err := ParseRequest(r1, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Action:       Fetch,
		Target:       ResourceCollection,
		ResourceType: "foo",
		PageNumber:   1,
		PageSize:     2,
	}, req)
}

func TestParseRequestFields(t *testing.T) {
	r1 := constructRequest("GET", "foo?fields[foo]=bar,baz")

	req, err := ParseRequest(r1, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Action:       Fetch,
		Target:       ResourceCollection,
		ResourceType: "foo",
		Fields: map[string][]string{
			"foo": {"bar", "baz"},
		},
	}, req)
}

func TestParseRequestFilters(t *testing.T) {
	r1 := constructRequest("GET", "foo?filter[foo]=bar,baz")

	req, err := ParseRequest(r1, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Action:       Fetch,
		Target:       ResourceCollection,
		ResourceType: "foo",
		Filters: map[string][]string{
			"foo": {"bar", "baz"},
		},
	}, req)
}

func TestParseRequestDocumentInvalidDocument(t *testing.T) {
	readers := []io.Reader{
		stringReader(``),
		stringReader(`1`),
		stringReader(`"foo"`),
		stringReader(`true`),
		stringReader(`[]`),
	}

	for _, r := range readers {
		doc, err := ParseRequestDocument(r)
		assert.Error(t, err)
		assert.Nil(t, doc)
	}
}

func TestParseRequestDocumentEmptyDocument(t *testing.T) {
	doc, err := ParseRequestDocument(stringReader(`{}`))
	assert.Error(t, err)
	assert.Nil(t, doc)
}

func TestParseRequestDocumentMinimumSingleDocument(t *testing.T) {
	doc, err := ParseRequestDocument(stringReader(`{
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

func TestParseRequestDocumentSingleDocument(t *testing.T) {
	doc, err := ParseRequestDocument(stringReader(`{
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

func BenchmarkParseRequest(b *testing.B) {
	r := constructRequest("GET", "foo/1")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseRequest(r, "")
	}
}

func BenchmarkParseRequestFilterAndSort(b *testing.B) {
	r := constructRequest("GET", "foo/1?sort=bar&filter[baz]=qux")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseRequest(r, "")
	}
}

func BenchmarkParseRequestDocument(b *testing.B) {
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
		_, err := ParseRequestDocument(reader)
		if err != nil {
			panic(err)
		}

		reader.Seek(0, io.SeekStart)
	}
}
