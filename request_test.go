package jsonapi

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRequestInvalidContentType(t *testing.T) {
	r := constructRequest("")
	r.Header.Set("Content-Type", "foo")

	req, err := ParseRequest(r, "")
	assert.Equal(t, ErrInvalidContentType, err)
	assert.Nil(t, req)
}

func TestParseRequestInvalidURL(t *testing.T) {
	list := []*http.Request{
		constructRequest(""),
		constructRequest("/"),
		constructRequest("foo/bar/baz/qux"),
		constructRequest("foo/bar/baz/qux/quux"),
		constructRequest("foo?page[number]=bar"),
		constructRequest("foo?page[size]=bar"),
		constructRequest("foo?page[number]=bar&page[number]=baz"),
		constructRequest("foo?page[size]=bar&page[size]=baz"),
	}

	for _, r := range list {
		req, err := ParseRequest(r, "")
		assert.Equal(t, ErrInvalidURL, err)
		assert.Nil(t, req)
	}
}

func TestParseRequestPrefix(t *testing.T) {
	r := constructRequest("foo/bar")

	req, err := ParseRequest(r, "foo/")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Resource: "bar",
	}, req)
}

func TestParseRequestResource(t *testing.T) {
	r := constructRequest("foo")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Resource: "foo",
	}, req)
}

func TestParseRequestResourceID(t *testing.T) {
	r := constructRequest("foo/1")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Resource:   "foo",
		ResourceID: "1",
	}, req)
}

func TestParseRequestRelatedResource(t *testing.T) {
	r := constructRequest("foo/1/bar")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Resource:        "foo",
		ResourceID:      "1",
		RelatedResource: "bar",
	}, req)
}

func TestParseRequestRelationship(t *testing.T) {
	r := constructRequest("foo/1/relationships/bar")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Resource:     "foo",
		ResourceID:   "1",
		Relationship: "bar",
	}, req)
}

func TestParseRequestInclude(t *testing.T) {
	r1 := constructRequest("foo?include=bar,baz")

	req, err := ParseRequest(r1, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Resource: "foo",
		Include:  []string{"bar", "baz"},
	}, req)

	r2 := constructRequest("foo?include=bar&include=baz,qux")

	req, err = ParseRequest(r2, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Resource: "foo",
		Include:  []string{"bar", "baz", "qux"},
	}, req)
}

func TestParseRequestSorting(t *testing.T) {
	r1 := constructRequest("foo?sort=bar,baz")

	req, err := ParseRequest(r1, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Resource: "foo",
		Sorting:  []string{"bar", "baz"},
	}, req)

	r2 := constructRequest("foo?sort=bar&sort=baz,qux")

	req, err = ParseRequest(r2, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Resource: "foo",
		Sorting:  []string{"bar", "baz", "qux"},
	}, req)
}

func TestParseRequestPage(t *testing.T) {
	r1 := constructRequest("foo?page[number]=1&page[size]=2")

	req, err := ParseRequest(r1, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Resource:   "foo",
		PageNumber: 1,
		PageSize:   2,
	}, req)
}

func TestParseRequestFields(t *testing.T) {
	r1 := constructRequest("foo?fields[foo]=bar,baz")

	req, err := ParseRequest(r1, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Resource: "foo",
		Fields: map[string][]string{
			"foo": {"bar", "baz"},
		},
	}, req)
}

func TestParseRequestFilters(t *testing.T) {
	r1 := constructRequest("foo?filter[foo]=bar,baz")

	req, err := ParseRequest(r1, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Resource: "foo",
		Filters: map[string][]string{
			"foo": {"bar", "baz"},
		},
	}, req)
}

func BenchmarkParseRequest(b *testing.B) {
	r := constructRequest("foo/1")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseRequest(r, "")
	}
}

func BenchmarkParseRequestFilterAndSort(b *testing.B) {
	r := constructRequest("foo/1?sort=bar&filter[baz]=qux")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseRequest(r, "")
	}
}
