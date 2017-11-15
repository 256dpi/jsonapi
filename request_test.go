package jsonapi

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRequestError(t *testing.T) {
	invalidAccept := newTestRequest("GET", "")
	invalidAccept.Header.Set("Accept", "foo")

	invalidContentType := newTestRequest("GET", "")
	invalidContentType.Header.Set("Content-Type", "foo")

	missingContentType := newTestRequest("POST", "foo")

	list := []*http.Request{
		invalidAccept,
		invalidContentType,
		missingContentType,
		newTestRequest("PUT", ""),
		newTestRequest("GET", ""),
		newTestRequest("POST", ""),
		newTestRequest("GET", "/"),
		newTestRequest("GET", "foo/bar/baz/qux"),
		newTestRequest("GET", "foo/bar/baz/qux/quux"),
		newTestRequest("GET", "foo?page[number]=bar"),
		newTestRequest("GET", "foo?page[size]=bar"),
		newTestRequest("GET", "foo?page[number]=1"),
		newTestRequest("GET", "foo?page[size]=1"),
		newTestRequest("GET", "foo?page[number]=bar&page[number]=baz"),
		newTestRequest("GET", "foo?page[size]=bar&page[size]=baz"),
		newTestRequest("PATCH", "foo"),
	}

	for _, r := range list {
		req, err := ParseRequest(r, "")
		assert.Error(t, err)
		assert.Nil(t, req)
	}
}

func TestParseRequestPrefix(t *testing.T) {
	list := map[string]string{
		"bar":           "",
		"/bar":          "",
		"foo/bar":       "foo",
		"/foo/bar":      "foo",
		"foo/bar/":      "/foo",
		"/foo/bar/":     "foo/",
		"baz/foo/bar/":  "/baz/foo",
		"/baz/foo/bar/": "baz/foo/",
	}

	for url, prefix := range list {
		r := newTestRequest("GET", url)

		req, err := ParseRequest(r, prefix)
		assert.NoError(t, err)
		assert.Equal(t, "bar", req.ResourceType)
	}
}

func TestParseRequestResource(t *testing.T) {
	r := newTestRequest("GET", "foo")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
	}, req)
}

func TestParseRequestResourceID(t *testing.T) {
	r := newTestRequest("GET", "foo/1")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       FindResource,
		ResourceType: "foo",
		ResourceID:   "1",
	}, req)
}

func TestParseRequestRelatedResource(t *testing.T) {
	r := newTestRequest("GET", "foo/1/bar")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:          GetRelatedResources,
		ResourceType:    "foo",
		ResourceID:      "1",
		RelatedResource: "bar",
	}, req)
}

func TestParseRequestRelationship(t *testing.T) {
	r := newTestRequest("GET", "foo/1/relationships/bar")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       GetRelationship,
		ResourceType: "foo",
		ResourceID:   "1",
		Relationship: "bar",
	}, req)
}

func TestParseRequest(t *testing.T) {
	list := []struct {
		method string
		url    string
		intent Intent
		doc    bool
		base   string
	}{
		{"GET", "/posts", ListResources, false, "/posts"},
		{"GET", "/posts/1", FindResource, false, "/posts/1"},
		{"POST", "/posts", CreateResource, true, "/posts"},
		{"PATCH", "/posts/1", UpdateResource, true, "/posts/1"},
		{"DELETE", "/posts/1", DeleteResource, false, "/posts/1"},
		{"GET", "/posts/1/author", GetRelatedResources, false, "/posts/1"},
		{"GET", "/posts/1/relationships/author", GetRelationship, false, "/posts/1"},
		{"PATCH", "/posts/1/relationships/author", SetRelationship, true, "/posts/1"},
		{"POST", "/posts/1/relationships/comments", AppendToRelationship, true, "/posts/1"},
		{"DELETE", "/posts/1/relationships/comments", RemoveFromRelationship, true, "/posts/1"},
	}

	for _, entry := range list {
		r := newTestRequest(entry.method, entry.url)
		r.Header.Set("Content-Type", MediaType)

		req, err := ParseRequest(r, "")
		assert.NoError(t, err)
		assert.Equal(t, entry.method, req.Intent.RequestMethod())
		assert.Equal(t, entry.url, req.Self())
		assert.Equal(t, entry.intent, req.Intent)
		assert.Equal(t, entry.doc, req.Intent.DocumentExpected())
		assert.Equal(t, entry.base, req.Base())
	}
}

func TestParseRequestWithPrefix(t *testing.T) {
	list := []struct {
		method string
		url    string
		intent Intent
		doc    bool
		base   string
	}{
		{"GET", "/api/posts", ListResources, false, "/api/posts"},
		{"GET", "/api/posts/1", FindResource, false, "/api/posts/1"},
		{"POST", "/api/posts", CreateResource, true, "/api/posts"},
		{"PATCH", "/api/posts/1", UpdateResource, true, "/api/posts/1"},
		{"DELETE", "/api/posts/1", DeleteResource, false, "/api/posts/1"},
		{"GET", "/api/posts/1/author", GetRelatedResources, false, "/api/posts/1"},
		{"GET", "/api/posts/1/relationships/author", GetRelationship, false, "/api/posts/1"},
		{"PATCH", "/api/posts/1/relationships/author", SetRelationship, true, "/api/posts/1"},
		{"POST", "/api/posts/1/relationships/comments", AppendToRelationship, true, "/api/posts/1"},
		{"DELETE", "/api/posts/1/relationships/comments", RemoveFromRelationship, true, "/api/posts/1"},
	}

	for _, entry := range list {
		r := newTestRequest(entry.method, entry.url)
		r.Header.Set("Content-Type", MediaType)

		req, err := ParseRequest(r, "api")
		assert.NoError(t, err)
		assert.Equal(t, entry.method, req.Intent.RequestMethod())
		assert.Equal(t, entry.url, req.Self())
		assert.Equal(t, entry.intent, req.Intent)
		assert.Equal(t, entry.doc, req.Intent.DocumentExpected())
		assert.Equal(t, entry.base, req.Base())
	}
}

func TestParseRequestInclude(t *testing.T) {
	r1 := newTestRequest("GET", "foo?include=bar,baz")

	req, err := ParseRequest(r1, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Include:      []string{"bar", "baz"},
	}, req)

	r2 := newTestRequest("GET", "foo?include=bar&include=baz,qux")

	req, err = ParseRequest(r2, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Include:      []string{"bar", "baz", "qux"},
	}, req)
}

func TestParseRequestSorting(t *testing.T) {
	r1 := newTestRequest("GET", "foo?sort=bar,baz")

	req, err := ParseRequest(r1, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Sorting:      []string{"bar", "baz"},
	}, req)

	r2 := newTestRequest("GET", "foo?sort=bar&sort=baz,qux")

	req, err = ParseRequest(r2, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Sorting:      []string{"bar", "baz", "qux"},
	}, req)
}

func TestParseRequestPage(t *testing.T) {
	r := newTestRequest("GET", "foo?page[number]=1&page[size]=2")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		PageNumber:   1,
		PageSize:     2,
	}, req)
}

func TestParseRequestFields(t *testing.T) {
	r := newTestRequest("GET", "foo?fields[foo]=bar,baz")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Fields: map[string][]string{
			"foo": {"bar", "baz"},
		},
	}, req)
}

func TestParseRequestFilters(t *testing.T) {
	r := newTestRequest("GET", "foo?filter[foo]=bar,baz")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Filters: map[string][]string{
			"foo": {"bar", "baz"},
		},
	}, req)
}

func TestZeroIntentRequestMethod(t *testing.T) {
	assert.Empty(t, Intent(0).RequestMethod())
}

func BenchmarkParseRequest(b *testing.B) {
	r := newTestRequest("GET", "foo/1")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseRequest(r, "")
	}
}

func BenchmarkParseRequestFilterAndSort(b *testing.B) {
	r := newTestRequest("GET", "foo/1?sort=bar&filter[baz]=qux")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseRequest(r, "")
	}
}
